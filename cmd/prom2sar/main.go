package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/prometheus/prometheus-dump-operator/pkg/sar"
	"github.com/prometheus/prometheus-dump-operator/pkg/tsdb"
)

const version = "1.0.0"

func main() {
	var (
		tsdbPath    string
		outputPath  string
		startTime   string
		endTime     string
		interval    int
		profile     string
		showVersion bool
		verbose     bool
		summaryOnly bool
	)

	flag.StringVar(&tsdbPath, "tsdb", "", "Path to Prometheus TSDB directory (required)")
	flag.StringVar(&outputPath, "output", "./sar-output", "Output directory for SAR files")
	flag.StringVar(&startTime, "start", "", "Start time (RFC3339 format, e.g., 2026-06-12T00:00:00Z)")
	flag.StringVar(&endTime, "end", "", "End time (RFC3339 format, e.g., 2026-06-12T23:59:59Z)")
	flag.IntVar(&interval, "interval", 60, "Sampling interval in seconds")
	flag.StringVar(&profile, "profile", "all", "Metrics profile: all, cpu, memory, disk, network")
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.BoolVar(&verbose, "verbose", false, "Verbose output")
	flag.BoolVar(&summaryOnly, "summary", false, "Generate summary only (no full report)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "prom2sar - Prometheus TSDB to SAR Converter v%s\n\n", version)
		fmt.Fprintf(os.Stderr, "Usage: prom2sar [options]\n\n")
		fmt.Fprintf(os.Stderr, "Convert Prometheus TSDB dumps to SAR format for analysis with standard Unix tools.\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  # Basic conversion (last 24 hours)\n")
		fmt.Fprintf(os.Stderr, "  prom2sar -tsdb /prometheus\n\n")
		fmt.Fprintf(os.Stderr, "  # Specific time range\n")
		fmt.Fprintf(os.Stderr, "  prom2sar -tsdb /prometheus -start 2026-06-12T00:00:00Z -end 2026-06-12T23:59:59Z\n\n")
		fmt.Fprintf(os.Stderr, "  # CPU metrics only\n")
		fmt.Fprintf(os.Stderr, "  prom2sar -tsdb /prometheus -profile cpu -output ./cpu-analysis\n\n")
		fmt.Fprintf(os.Stderr, "  # Quick summary\n")
		fmt.Fprintf(os.Stderr, "  prom2sar -tsdb /prometheus -summary\n\n")
		fmt.Fprintf(os.Stderr, "Output:\n")
		fmt.Fprintf(os.Stderr, "  SAR files are written to the output directory:\n")
		fmt.Fprintf(os.Stderr, "  - sar-YYYYMMDD.txt         Main SAR report (use this for analysis)\n")
		fmt.Fprintf(os.Stderr, "  - sar-summary-YYYYMMDD.txt Summary and quick stats\n\n")
		fmt.Fprintf(os.Stderr, "Profiles:\n")
		fmt.Fprintf(os.Stderr, "  all     - CPU, memory, disk, network (default)\n")
		fmt.Fprintf(os.Stderr, "  cpu     - CPU utilization only (sar -u)\n")
		fmt.Fprintf(os.Stderr, "  memory  - Memory utilization only (sar -r)\n")
		fmt.Fprintf(os.Stderr, "  disk    - Disk I/O only (sar -d)\n")
		fmt.Fprintf(os.Stderr, "  network - Network I/O only (sar -n DEV)\n\n")
	}

	flag.Parse()

	if showVersion {
		fmt.Printf("prom2sar version %s\n", version)
		os.Exit(0)
	}

	if tsdbPath == "" {
		fmt.Fprintf(os.Stderr, "Error: -tsdb is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Validate TSDB path exists
	if _, err := os.Stat(tsdbPath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: TSDB path does not exist: %s\n", tsdbPath)
		os.Exit(1)
	}

	// Parse time range
	var start, end time.Time
	var err error

	if startTime != "" {
		start, err = time.Parse(time.RFC3339, startTime)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid start time format. Use RFC3339 (e.g., 2026-06-12T00:00:00Z)\n")
			os.Exit(1)
		}
	} else {
		start = time.Now().Add(-24 * time.Hour)
		if verbose {
			fmt.Printf("No start time specified, using: %s\n", start.Format(time.RFC3339))
		}
	}

	if endTime != "" {
		end, err = time.Parse(time.RFC3339, endTime)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid end time format. Use RFC3339 (e.g., 2026-06-12T23:59:59Z)\n")
			os.Exit(1)
		}
	} else {
		end = time.Now()
		if verbose {
			fmt.Printf("No end time specified, using: %s\n", end.Format(time.RFC3339))
		}
	}

	if start.After(end) {
		fmt.Fprintf(os.Stderr, "Error: Start time must be before end time\n")
		os.Exit(1)
	}

	// Validate profile
	validProfiles := map[string]bool{
		"all": true, "cpu": true, "memory": true, "disk": true, "network": true,
	}
	if !validProfiles[profile] {
		fmt.Fprintf(os.Stderr, "Error: Invalid profile '%s'. Valid: all, cpu, memory, disk, network\n", profile)
		os.Exit(1)
	}

	// Create output directory
	if err := os.MkdirAll(outputPath, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to create output directory: %v\n", err)
		os.Exit(1)
	}

	// Run conversion
	if verbose {
		fmt.Println("=== Prometheus to SAR Conversion ===")
		fmt.Printf("TSDB Path:    %s\n", tsdbPath)
		fmt.Printf("Output Path:  %s\n", outputPath)
		fmt.Printf("Time Range:   %s to %s\n", start.Format(time.RFC3339), end.Format(time.RFC3339))
		fmt.Printf("Interval:     %d seconds\n", interval)
		fmt.Printf("Profile:      %s\n", profile)
		fmt.Println()
	}

	ctx := context.Background()

	if err := convertToSar(ctx, tsdbPath, outputPath, start, end, int64(interval), profile, verbose, summaryOnly); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Conversion failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Conversion completed successfully!")
}

func convertToSar(ctx context.Context, tsdbPath, outputPath string, start, end time.Time, interval int64, profile string, verbose, summaryOnly bool) error {
	// Open TSDB reader
	if verbose {
		fmt.Printf("Opening TSDB at %s...\n", tsdbPath)
	}

	reader, err := tsdb.NewReader(tsdbPath)
	if err != nil {
		return fmt.Errorf("failed to open TSDB: %w", err)
	}
	defer reader.Close()

	// Get block info
	blocks, err := reader.GetBlocks()
	if err != nil {
		return fmt.Errorf("failed to get blocks: %w", err)
	}

	if verbose {
		fmt.Printf("Found %d TSDB blocks\n", len(blocks))
		for _, block := range blocks {
			blockStart := time.Unix(block.MinTime/1000, 0)
			blockEnd := time.Unix(block.MaxTime/1000, 0)
			fmt.Printf("  - Block %s: %s to %s\n", block.ULID, blockStart.Format(time.RFC3339), blockEnd.Format(time.RFC3339))
		}
		fmt.Println()
	}

	// Extract metrics
	if verbose {
		fmt.Printf("Extracting %s metrics...\n", profile)
	}

	mapper := sar.NewMetricsMapper(reader)

	startMillis := tsdb.TimeRangeToMillis(start)
	endMillis := tsdb.TimeRangeToMillis(end)

	sarData, err := mapper.ExtractSarData(ctx, startMillis, endMillis, interval, profile)
	if err != nil {
		return fmt.Errorf("failed to extract SAR data: %w", err)
	}

	if len(sarData) == 0 {
		fmt.Println("Warning: No data found in the specified time range")
		return nil
	}

	if verbose {
		fmt.Printf("Extracted %d data points\n\n", len(sarData))
	}

	// Generate output
	generator := sar.NewGenerator("prometheus-tsdb")

	dateStr := start.Format("20060102")

	// Generate summary
	summaryFile := filepath.Join(outputPath, fmt.Sprintf("sar-summary-%s.txt", dateStr))
	sf, err := os.Create(summaryFile)
	if err != nil {
		return fmt.Errorf("failed to create summary file: %w", err)
	}
	defer sf.Close()

	if err := generator.GenerateSummary(sarData, sf); err != nil {
		return fmt.Errorf("failed to generate summary: %w", err)
	}

	if err := generator.GenerateCompactReport(sarData, sf); err != nil {
		return fmt.Errorf("failed to generate compact report: %w", err)
	}

	fmt.Printf("✓ Generated summary: %s\n", summaryFile)

	// Show summary on stdout if verbose
	if verbose || summaryOnly {
		fmt.Println()
		sf.Seek(0, 0)
		content, _ := os.ReadFile(summaryFile)
		fmt.Println(string(content))
	}

	// Generate full report unless summary-only
	if !summaryOnly {
		reportFile := filepath.Join(outputPath, fmt.Sprintf("sar-%s.txt", dateStr))
		rf, err := os.Create(reportFile)
		if err != nil {
			return fmt.Errorf("failed to create report file: %w", err)
		}
		defer rf.Close()

		if err := generator.GenerateTextReport(sarData, rf); err != nil {
			return fmt.Errorf("failed to generate text report: %w", err)
		}

		fmt.Printf("✓ Generated report:  %s\n", reportFile)

		if verbose {
			fmt.Printf("\nPreview (first 50 lines):\n")
			fmt.Println(string(repeatChar('=', 80)))
			rf.Seek(0, 0)
			content, _ := os.ReadFile(reportFile)
			lines := 0
			for _, line := range string(content) {
				if line == '\n' {
					lines++
					if lines >= 50 {
						fmt.Println("... (truncated, see full file)")
						break
					}
				}
				fmt.Printf("%c", line)
			}
		}
	}

	return nil
}

func repeatChar(c rune, n int) string {
	result := make([]rune, n)
	for i := range result {
		result[i] = c
	}
	return string(result)
}
