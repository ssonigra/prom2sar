package tsdb

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/tsdb"
	"github.com/prometheus/prometheus/tsdb/chunkenc"
)

// Reader reads Prometheus TSDB blocks
type Reader struct {
	db *tsdb.DBReadOnly
}

// TimeSeriesData represents a single time series
type TimeSeriesData struct {
	Labels  labels.Labels
	Samples []Sample
}

// Sample represents a single data point
type Sample struct {
	Timestamp int64
	Value     float64
}

// NewReader creates a new TSDB reader
func NewReader(dbPath string) (*Reader, error) {
	db, err := tsdb.OpenDBReadOnly(dbPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to open TSDB: %w", err)
	}

	return &Reader{db: db}, nil
}

// Close closes the reader
func (r *Reader) Close() error {
	return r.db.Close()
}

// Query queries time series data
func (r *Reader) Query(ctx context.Context, startTime, endTime int64, matchers ...*labels.Matcher) ([]*TimeSeriesData, error) {
	querier, err := r.db.Querier(startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to create querier: %w", err)
	}
	defer querier.Close()

	ss := querier.Select(ctx, false, nil, matchers...)

	var result []*TimeSeriesData
	for ss.Next() {
		series := ss.At()

		tsData := &TimeSeriesData{
			Labels:  series.Labels(),
			Samples: make([]Sample, 0),
		}

		it := series.Iterator(nil)
		for it.Next() == chunkenc.ValFloat {
			t, v := it.At()
			tsData.Samples = append(tsData.Samples, Sample{
				Timestamp: t,
				Value:     v,
			})
		}

		if it.Err() != nil {
			return nil, fmt.Errorf("iterator error: %w", it.Err())
		}

		result = append(result, tsData)
	}

	if ss.Err() != nil {
		return nil, fmt.Errorf("series set error: %w", ss.Err())
	}

	return result, nil
}

// QueryMetric queries a specific metric name
func (r *Reader) QueryMetric(ctx context.Context, metricName string, startTime, endTime int64, labelFilters map[string]string) ([]*TimeSeriesData, error) {
	matchers := []*labels.Matcher{
		labels.MustNewMatcher(labels.MatchEqual, "__name__", metricName),
	}

	for k, v := range labelFilters {
		matchers = append(matchers, labels.MustNewMatcher(labels.MatchEqual, k, v))
	}

	return r.Query(ctx, startTime, endTime, matchers...)
}

// GetBlocks returns information about available blocks
func (r *Reader) GetBlocks() ([]BlockInfo, error) {
	blocks, err := r.db.Blocks()
	if err != nil {
		return nil, fmt.Errorf("failed to get blocks: %w", err)
	}

	result := make([]BlockInfo, len(blocks))
	for i, block := range blocks {
		meta := block.Meta()
		result[i] = BlockInfo{
			ULID:    meta.ULID.String(),
			MinTime: meta.MinTime,
			MaxTime: meta.MaxTime,
			Stats:   meta.Stats,
		}
	}

	return result, nil
}

// BlockInfo contains information about a TSDB block
type BlockInfo struct {
	ULID    string
	MinTime int64
	MaxTime int64
	Stats   tsdb.BlockStats
}

// DiscoverTSDBBlocks finds all TSDB block directories in a path
func DiscoverTSDBBlocks(basePath string) ([]string, error) {
	var blocks []string

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		blockPath := filepath.Join(basePath, entry.Name())

		metaPath := filepath.Join(blockPath, "meta.json")
		if _, err := os.Stat(metaPath); err == nil {
			blocks = append(blocks, blockPath)
		}
	}

	return blocks, nil
}

// TimeRangeToMillis converts time.Time to Prometheus timestamp (milliseconds)
func TimeRangeToMillis(t time.Time) int64 {
	return t.Unix() * 1000
}
