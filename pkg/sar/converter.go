package sar

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	v1alpha1 "github.com/prometheus/prometheus-dump-operator/pkg/apis/prometheus/v1alpha1"
	"github.com/prometheus/prometheus-dump-operator/pkg/tsdb"
)

// Converter orchestrates the conversion from Prometheus TSDB to sar format
type Converter struct {
	spec *v1alpha1.SarConversionSpec
}

// NewConverter creates a new converter
func NewConverter(spec *v1alpha1.SarConversionSpec) *Converter {
	return &Converter{spec: spec}
}

// Convert performs the conversion from TSDB to sar format
func (c *Converter) Convert(ctx context.Context, tsdbPath string, startTime, endTime time.Time) (*v1alpha1.SarConversionStatus, error) {
	status := &v1alpha1.SarConversionStatus{
		Phase: "Starting",
	}

	reader, err := tsdb.NewReader(tsdbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create TSDB reader: %w", err)
	}
	defer reader.Close()

	status.Phase = "ExtractingMetrics"

	mapper := NewMetricsMapper(reader)

	interval := int64(60)
	if c.spec.Interval > 0 {
		interval = int64(c.spec.Interval)
	}

	profile := c.spec.MetricsProfile
	if profile == "" {
		profile = "all"
	}

	start := tsdb.TimeRangeToMillis(startTime)
	end := tsdb.TimeRangeToMillis(endTime)

	sarData, err := mapper.ExtractSarData(ctx, start, end, interval, profile)
	if err != nil {
		return nil, fmt.Errorf("failed to extract sar data: %w", err)
	}

	status.MetricsConverted = len(sarData)
	status.Phase = "GeneratingOutput"

	if err := os.MkdirAll(c.spec.OutputPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	outputFormat := c.spec.Format
	if outputFormat == "" {
		outputFormat = v1alpha1.SarFormatText
	}

	var outputFiles []string

	switch outputFormat {
	case v1alpha1.SarFormatText:
		outputFiles, err = c.generateTextOutput(sarData, startTime, endTime)
	case v1alpha1.SarFormatBinary:
		return nil, fmt.Errorf("binary format not yet implemented - use text format")
	default:
		return nil, fmt.Errorf("unsupported format: %s", outputFormat)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to generate output: %w", err)
	}

	status.SarFilesGenerated = len(outputFiles)
	status.OutputLocation = c.spec.OutputPath
	status.Phase = "Completed"
	status.TimeRangeCovered = &v1alpha1.TimeRange{
		Start: toMetaTime(startTime),
		End:   toMetaTime(endTime),
	}

	return status, nil
}

// generateTextOutput generates text format sar output
func (c *Converter) generateTextOutput(data []*SarData, startTime, endTime time.Time) ([]string, error) {
	var outputFiles []string

	hostname := "prometheus-tsdb"
	generator := NewGenerator(hostname)

	dateStr := startTime.Format("20060102")
	reportFile := filepath.Join(c.spec.OutputPath, fmt.Sprintf("sar-%s.txt", dateStr))

	f, err := os.Create(reportFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create report file: %w", err)
	}
	defer f.Close()

	if err := generator.GenerateTextReport(data, f); err != nil {
		return nil, fmt.Errorf("failed to generate text report: %w", err)
	}

	outputFiles = append(outputFiles, reportFile)

	summaryFile := filepath.Join(c.spec.OutputPath, fmt.Sprintf("sar-summary-%s.txt", dateStr))
	sf, err := os.Create(summaryFile)
	if err != nil {
		return nil, fmt.Errorf("failed to create summary file: %w", err)
	}
	defer sf.Close()

	if err := generator.GenerateSummary(data, sf); err != nil {
		return nil, fmt.Errorf("failed to generate summary: %w", err)
	}

	if err := generator.GenerateCompactReport(data, sf); err != nil {
		return nil, fmt.Errorf("failed to generate compact report: %w", err)
	}

	outputFiles = append(outputFiles, summaryFile)

	return outputFiles, nil
}

func toMetaTime(t time.Time) metav1.Time {
	return metav1.Time{Time: t}
}
