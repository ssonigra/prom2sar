package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PrometheusDumpLoader represents a request to load Prometheus dumps
type PrometheusDumpLoader struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PrometheusDumpLoaderSpec   `json:"spec"`
	Status PrometheusDumpLoaderStatus `json:"status,omitempty"`
}

// PrometheusDumpLoaderSpec defines the desired state
type PrometheusDumpLoaderSpec struct {
	// SourcePath is the path where Prometheus dumps are located
	SourcePath string `json:"sourcePath"`

	// TargetPath is where dumps should be copied
	TargetPath string `json:"targetPath"`

	// TimeRange optionally filters dumps by time
	TimeRange *TimeRange `json:"timeRange,omitempty"`

	// Filters for selecting specific dumps
	Filters *DumpFilters `json:"filters,omitempty"`

	// Compression indicates whether to compress copied dumps
	Compression bool `json:"compression,omitempty"`

	// SarConversion enables conversion to sar-compatible format
	SarConversion *SarConversionSpec `json:"sarConversion,omitempty"`
}

// TimeRange specifies a time window
type TimeRange struct {
	Start metav1.Time `json:"start,omitempty"`
	End   metav1.Time `json:"end,omitempty"`
}

// DumpFilters allows filtering dumps
type DumpFilters struct {
	// Metrics is a list of metric names to include
	Metrics []string `json:"metrics,omitempty"`

	// Labels are label matchers
	Labels map[string]string `json:"labels,omitempty"`
}

// SarConversionSpec defines sar conversion options
type SarConversionSpec struct {
	// Enabled enables sar conversion
	Enabled bool `json:"enabled"`

	// OutputPath is where sar files will be written
	OutputPath string `json:"outputPath"`

	// Format specifies the output format (text, binary)
	// text: human-readable sar-like output
	// binary: sar-compatible binary format (sadc format)
	Format SarOutputFormat `json:"format,omitempty"`

	// Interval is the sampling interval in seconds (default: 60)
	Interval int `json:"interval,omitempty"`

	// MetricsProfile defines which system metrics to extract
	// Options: all, cpu, memory, disk, network, custom
	MetricsProfile string `json:"metricsProfile,omitempty"`

	// CustomMetrics allows specifying custom metric mappings
	CustomMetrics []MetricMapping `json:"customMetrics,omitempty"`
}

// SarOutputFormat represents the output format
type SarOutputFormat string

const (
	SarFormatText   SarOutputFormat = "text"
	SarFormatBinary SarOutputFormat = "binary"
)

// MetricMapping maps Prometheus metrics to sar fields
type MetricMapping struct {
	// PrometheusMetric is the source metric name
	PrometheusMetric string `json:"prometheusMetric"`

	// SarField is the target sar field
	SarField string `json:"sarField"`

	// LabelSelector filters metrics by labels
	LabelSelector map[string]string `json:"labelSelector,omitempty"`

	// Aggregation method (avg, sum, max, min)
	Aggregation string `json:"aggregation,omitempty"`
}

// PrometheusDumpLoaderStatus represents the current state
type PrometheusDumpLoaderStatus struct {
	// Phase is the current phase
	Phase string `json:"phase,omitempty"`

	// Message is a human-readable status message
	Message string `json:"message,omitempty"`

	// BytesCopied is the total bytes copied
	BytesCopied int64 `json:"bytesCopied,omitempty"`

	// FilesCopied is the number of files copied
	FilesCopied int `json:"filesCopied,omitempty"`

	// LastUpdateTime is when status was last updated
	LastUpdateTime metav1.Time `json:"lastUpdateTime,omitempty"`

	// Conditions represent the latest available observations
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// SarConversionStatus tracks sar conversion progress
	SarConversionStatus *SarConversionStatus `json:"sarConversionStatus,omitempty"`
}

// SarConversionStatus tracks sar conversion progress
type SarConversionStatus struct {
	// Phase is the conversion phase
	Phase string `json:"phase,omitempty"`

	// MetricsConverted is the number of metrics converted
	MetricsConverted int `json:"metricsConverted,omitempty"`

	// SarFilesGenerated is the number of sar files created
	SarFilesGenerated int `json:"sarFilesGenerated,omitempty"`

	// OutputLocation is where sar files were written
	OutputLocation string `json:"outputLocation,omitempty"`

	// TimeRangeCovered shows the time range of converted data
	TimeRangeCovered *TimeRange `json:"timeRangeCovered,omitempty"`
}

// Phase constants
const (
	PhasePending    = "Pending"
	PhaseInProgress = "InProgress"
	PhaseCompleted  = "Completed"
	PhaseFailed     = "Failed"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PrometheusDumpLoaderList contains a list of PrometheusDumpLoader
type PrometheusDumpLoaderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PrometheusDumpLoader `json:"items"`
}
