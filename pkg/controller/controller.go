package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/prometheus/prometheus-dump-operator/pkg/apis/prometheus/v1alpha1"
	"github.com/prometheus/prometheus-dump-operator/pkg/loader"
	"github.com/prometheus/prometheus-dump-operator/pkg/sar"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// PrometheusDumpLoaderReconciler reconciles PrometheusDumpLoader resources
type PrometheusDumpLoaderReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Config *rest.Config
}

// Reconcile handles reconciliation of PrometheusDumpLoader resources
func (r *PrometheusDumpLoaderReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling PrometheusDumpLoader", "name", req.Name, "namespace", req.Namespace)

	// Fetch the PrometheusDumpLoader instance
	instance := &v1alpha1.PrometheusDumpLoader{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("PrometheusDumpLoader not found, may have been deleted")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Skip if already completed
	if instance.Status.Phase == v1alpha1.PhaseCompleted {
		logger.Info("PrometheusDumpLoader already completed", "name", instance.Name)
		return ctrl.Result{}, nil
	}

	// Update status to InProgress if pending
	if instance.Status.Phase == "" || instance.Status.Phase == v1alpha1.PhasePending {
		instance.Status.Phase = v1alpha1.PhaseInProgress
		instance.Status.Message = "Starting dump load operation"
		instance.Status.LastUpdateTime = metav1.Now()
		if err := r.Status().Update(ctx, instance); err != nil {
			return ctrl.Result{}, err
		}
	}

	// Create loader
	dumpLoader := loader.NewDumpLoader(
		instance.Spec.SourcePath,
		instance.Spec.TargetPath,
		instance.Spec.Compression,
	)

	// Perform the load operation
	result, err := dumpLoader.Load(ctx, &instance.Spec)
	if err != nil {
		logger.Error(err, "Failed to load dumps")
		instance.Status.Phase = v1alpha1.PhaseFailed
		instance.Status.Message = fmt.Sprintf("Failed to load dumps: %v", err)
		instance.Status.LastUpdateTime = metav1.Now()

		// Add error condition
		instance.Status.Conditions = append(instance.Status.Conditions, metav1.Condition{
			Type:               "Failed",
			Status:             metav1.ConditionTrue,
			LastTransitionTime: metav1.Now(),
			Reason:             "LoadFailed",
			Message:            err.Error(),
		})

		if updateErr := r.Status().Update(ctx, instance); updateErr != nil {
			logger.Error(updateErr, "Failed to update status")
		}
		return ctrl.Result{RequeueAfter: 30 * time.Second}, err
	}

	// Update status with results
	instance.Status.FilesCopied = result.FilesCopied
	instance.Status.BytesCopied = result.BytesCopied
	instance.Status.LastUpdateTime = metav1.Now()

	// Perform sar conversion if enabled
	if instance.Spec.SarConversion != nil && instance.Spec.SarConversion.Enabled {
		logger.Info("Starting sar conversion")

		sarConverter := sar.NewConverter(instance.Spec.SarConversion)

		startTime := time.Now().Add(-24 * time.Hour)
		endTime := time.Now()

		if instance.Spec.TimeRange != nil {
			if !instance.Spec.TimeRange.Start.IsZero() {
				startTime = instance.Spec.TimeRange.Start.Time
			}
			if !instance.Spec.TimeRange.End.IsZero() {
				endTime = instance.Spec.TimeRange.End.Time
			}
		}

		sarStatus, sarErr := sarConverter.Convert(ctx, instance.Spec.TargetPath, startTime, endTime)
		if sarErr != nil {
			logger.Error(sarErr, "Failed to convert to sar format")
			instance.Status.Message = fmt.Sprintf("Loaded %d files but sar conversion failed: %v", result.FilesCopied, sarErr)

			instance.Status.Conditions = append(instance.Status.Conditions, metav1.Condition{
				Type:               "SarConversionFailed",
				Status:             metav1.ConditionTrue,
				LastTransitionTime: metav1.Now(),
				Reason:             "ConversionError",
				Message:            sarErr.Error(),
			})
		} else {
			logger.Info("Sar conversion completed successfully",
				"metricsConverted", sarStatus.MetricsConverted,
				"filesGenerated", sarStatus.SarFilesGenerated)

			instance.Status.SarConversionStatus = sarStatus

			instance.Status.Conditions = append(instance.Status.Conditions, metav1.Condition{
				Type:               "SarConversionCompleted",
				Status:             metav1.ConditionTrue,
				LastTransitionTime: metav1.Now(),
				Reason:             "ConversionSucceeded",
				Message:            fmt.Sprintf("Generated %d sar files", sarStatus.SarFilesGenerated),
			})
		}
	}

	if len(result.Errors) > 0 {
		instance.Status.Phase = v1alpha1.PhaseFailed
		instance.Status.Message = fmt.Sprintf("Completed with %d errors", len(result.Errors))

		// Add partial failure condition
		instance.Status.Conditions = append(instance.Status.Conditions, metav1.Condition{
			Type:               "PartialFailure",
			Status:             metav1.ConditionTrue,
			LastTransitionTime: metav1.Now(),
			Reason:             "SomeFilesFailed",
			Message:            fmt.Sprintf("%d files failed to copy", len(result.Errors)),
		})
	} else {
		instance.Status.Phase = v1alpha1.PhaseCompleted
		instance.Status.Message = fmt.Sprintf("Successfully copied %d files (%d bytes)", result.FilesCopied, result.BytesCopied)

		// Add success condition
		instance.Status.Conditions = append(instance.Status.Conditions, metav1.Condition{
			Type:               "Completed",
			Status:             metav1.ConditionTrue,
			LastTransitionTime: metav1.Now(),
			Reason:             "LoadSucceeded",
			Message:            "All files copied successfully",
		})
	}

	if err := r.Status().Update(ctx, instance); err != nil {
		logger.Error(err, "Failed to update status")
		return ctrl.Result{}, err
	}

	logger.Info("Reconciliation complete",
		"filesCopied", result.FilesCopied,
		"bytesCopied", result.BytesCopied,
		"errors", len(result.Errors))

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager
func (r *PrometheusDumpLoaderReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.PrometheusDumpLoader{}).
		Complete(r)
}
