package loader

import (
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/prometheus/prometheus-dump-operator/pkg/apis/prometheus/v1alpha1"
)

// DumpLoader handles copying Prometheus dumps
type DumpLoader struct {
	sourcePath string
	targetPath string
	compress   bool
}

// NewDumpLoader creates a new loader instance
func NewDumpLoader(sourcePath, targetPath string, compress bool) *DumpLoader {
	return &DumpLoader{
		sourcePath: sourcePath,
		targetPath: targetPath,
		compress:   compress,
	}
}

// LoadResult contains the results of a load operation
type LoadResult struct {
	FilesCopied int
	BytesCopied int64
	Errors      []error
}

// Load copies Prometheus dumps from source to target
func (l *DumpLoader) Load(ctx context.Context, spec *v1alpha1.PrometheusDumpLoaderSpec) (*LoadResult, error) {
	result := &LoadResult{
		Errors: make([]error, 0),
	}

	// Ensure target directory exists
	if err := os.MkdirAll(spec.TargetPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create target directory: %w", err)
	}

	// Walk the source directory
	err := filepath.Walk(spec.SourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("error accessing %s: %w", path, err))
			return nil // Continue walking
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Filter by time range if specified
		if spec.TimeRange != nil {
			if !l.isInTimeRange(info.ModTime(), spec.TimeRange) {
				return nil
			}
		}

		// Calculate relative path
		relPath, err := filepath.Rel(spec.SourcePath, path)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("failed to get relative path for %s: %w", path, err))
			return nil
		}

		// Copy the file
		targetPath := filepath.Join(spec.TargetPath, relPath)
		bytes, err := l.copyFile(path, targetPath, spec.Compression)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("failed to copy %s: %w", path, err))
			return nil
		}

		result.FilesCopied++
		result.BytesCopied += bytes

		return nil
	})

	if err != nil && err != context.Canceled {
		return nil, fmt.Errorf("failed to walk source directory: %w", err)
	}

	return result, nil
}

// copyFile copies a single file, optionally compressing it
func (l *DumpLoader) copyFile(src, dst string, compress bool) (int64, error) {
	// Ensure parent directory exists
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return 0, err
	}

	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer srcFile.Close()

	// Adjust destination path if compressing
	if compress && filepath.Ext(dst) != ".gz" {
		dst = dst + ".gz"
	}

	// Create destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer dstFile.Close()

	// Copy with optional compression
	var written int64
	if compress {
		gzWriter := gzip.NewWriter(dstFile)
		defer gzWriter.Close()
		written, err = io.Copy(gzWriter, srcFile)
	} else {
		written, err = io.Copy(dstFile, srcFile)
	}

	if err != nil {
		return 0, err
	}

	return written, nil
}

// isInTimeRange checks if a time falls within the specified range
func (l *DumpLoader) isInTimeRange(t time.Time, timeRange *v1alpha1.TimeRange) bool {
	if !timeRange.Start.IsZero() && t.Before(timeRange.Start.Time) {
		return false
	}
	if !timeRange.End.IsZero() && t.After(timeRange.End.Time) {
		return false
	}
	return true
}
