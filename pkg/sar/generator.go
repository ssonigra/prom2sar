package sar

import (
	"fmt"
	"io"
	"strings"
	"time"
)

// Generator generates sar-compatible output
type Generator struct {
	hostname string
}

// NewGenerator creates a new sar generator
func NewGenerator(hostname string) *Generator {
	if hostname == "" {
		hostname = "localhost"
	}
	return &Generator{hostname: hostname}
}

// GenerateTextReport generates a text report similar to sar output
func (g *Generator) GenerateTextReport(data []*SarData, writer io.Writer) error {
	if len(data) == 0 {
		return fmt.Errorf("no data to generate report")
	}

	fmt.Fprintf(writer, "Linux %s\t%s\n\n", "TSDB", time.Now().Format("01/02/2006"))

	if data[0].CPU != nil {
		if err := g.generateCPUReport(data, writer); err != nil {
			return err
		}
	}

	if data[0].Memory != nil {
		if err := g.generateMemoryReport(data, writer); err != nil {
			return err
		}
	}

	if data[0].Disk != nil {
		if err := g.generateDiskReport(data, writer); err != nil {
			return err
		}
	}

	if data[0].Network != nil {
		if err := g.generateNetworkReport(data, writer); err != nil {
			return err
		}
	}

	return nil
}

// generateCPUReport generates CPU utilization report (sar -u)
func (g *Generator) generateCPUReport(data []*SarData, writer io.Writer) error {
	fmt.Fprintf(writer, "%s\n", strings.Repeat("-", 80))
	fmt.Fprintf(writer, "CPU Utilization (sar -u)\n")
	fmt.Fprintf(writer, "%s\n\n", strings.Repeat("-", 80))

	fmt.Fprintf(writer, "%-12s %-8s %7s %7s %7s %7s %7s %7s\n",
		"Timestamp", "CPU", "%user", "%nice", "%system", "%iowait", "%steal", "%idle")

	for _, d := range data {
		if d.CPU == nil {
			continue
		}

		ts := time.Unix(d.Timestamp/1000, 0).Format("15:04:05")

		fmt.Fprintf(writer, "%-12s %-8s %7.2f %7.2f %7.2f %7.2f %7.2f %7.2f\n",
			ts, "all",
			d.CPU.User,
			d.CPU.Nice,
			d.CPU.System,
			d.CPU.IOWait,
			d.CPU.Steal,
			d.CPU.Idle)
	}

	fmt.Fprintf(writer, "\n")
	if len(data) > 0 && data[0].CPU != nil {
		fmt.Fprintf(writer, "Average CPU count: %d\n\n", data[0].CPU.NumCPUs)
	}

	return nil
}

// generateMemoryReport generates memory utilization report (sar -r)
func (g *Generator) generateMemoryReport(data []*SarData, writer io.Writer) error {
	fmt.Fprintf(writer, "%s\n", strings.Repeat("-", 100))
	fmt.Fprintf(writer, "Memory Utilization (sar -r)\n")
	fmt.Fprintf(writer, "%s\n\n", strings.Repeat("-", 100))

	fmt.Fprintf(writer, "%-12s %12s %12s %12s %12s %12s %8s\n",
		"Timestamp", "kbmemfree", "kbmemused", "kbbuffers", "kbcached", "kbswpfree", "%memused")

	for _, d := range data {
		if d.Memory == nil {
			continue
		}

		ts := time.Unix(d.Timestamp/1000, 0).Format("15:04:05")

		fmt.Fprintf(writer, "%-12s %12.0f %12.0f %12.0f %12.0f %12.0f %8.2f\n",
			ts,
			d.Memory.MemFree,
			d.Memory.MemUsed,
			d.Memory.MemBuffers,
			d.Memory.MemCached,
			d.Memory.SwapFree,
			d.Memory.MemPercent)
	}

	fmt.Fprintf(writer, "\n")
	return nil
}

// generateDiskReport generates disk I/O report (sar -d)
func (g *Generator) generateDiskReport(data []*SarData, writer io.Writer) error {
	fmt.Fprintf(writer, "%s\n", strings.Repeat("-", 100))
	fmt.Fprintf(writer, "Disk I/O Statistics (sar -d)\n")
	fmt.Fprintf(writer, "%s\n\n", strings.Repeat("-", 100))

	fmt.Fprintf(writer, "%-12s %-12s %8s %12s %12s %10s %10s %8s\n",
		"Timestamp", "DEV", "tps", "rd_sec/s", "wr_sec/s", "avgrq-sz", "avgqu-sz", "%util")

	for _, d := range data {
		if d.Disk == nil || len(d.Disk.Devices) == 0 {
			continue
		}

		ts := time.Unix(d.Timestamp/1000, 0).Format("15:04:05")

		for _, dev := range d.Disk.Devices {
			fmt.Fprintf(writer, "%-12s %-12s %8.2f %12.2f %12.2f %10.2f %10.2f %8.2f\n",
				ts,
				dev.Device,
				dev.TPS,
				dev.RdKBs*2,
				dev.WrKBs*2,
				dev.AvgRqSz,
				dev.AvgQuSz,
				dev.Util)
		}
	}

	fmt.Fprintf(writer, "\n")
	return nil
}

// generateNetworkReport generates network I/O report (sar -n DEV)
func (g *Generator) generateNetworkReport(data []*SarData, writer io.Writer) error {
	fmt.Fprintf(writer, "%s\n", strings.Repeat("-", 110))
	fmt.Fprintf(writer, "Network Statistics (sar -n DEV)\n")
	fmt.Fprintf(writer, "%s\n\n", strings.Repeat("-", 110))

	fmt.Fprintf(writer, "%-12s %-12s %10s %10s %12s %12s %10s %10s\n",
		"Timestamp", "IFACE", "rxpck/s", "txpck/s", "rxkB/s", "txkB/s", "rxerr/s", "txerr/s")

	for _, d := range data {
		if d.Network == nil || len(d.Network.Interfaces) == 0 {
			continue
		}

		ts := time.Unix(d.Timestamp/1000, 0).Format("15:04:05")

		for _, iface := range d.Network.Interfaces {
			fmt.Fprintf(writer, "%-12s %-12s %10.2f %10.2f %12.2f %12.2f %10.2f %10.2f\n",
				ts,
				iface.Interface,
				iface.RxPckS,
				iface.TxPckS,
				iface.RxKBS,
				iface.TxKBS,
				iface.RxErrS,
				iface.TxErrS)
		}
	}

	fmt.Fprintf(writer, "\n")
	return nil
}

// GenerateSummary generates a summary report
func (g *Generator) GenerateSummary(data []*SarData, writer io.Writer) error {
	if len(data) == 0 {
		return fmt.Errorf("no data to generate summary")
	}

	fmt.Fprintf(writer, "=== Prometheus to SAR Conversion Summary ===\n\n")
	fmt.Fprintf(writer, "Hostname: %s\n", g.hostname)
	fmt.Fprintf(writer, "Total samples: %d\n", len(data))

	if len(data) > 0 {
		startTime := time.Unix(data[0].Timestamp/1000, 0)
		endTime := time.Unix(data[len(data)-1].Timestamp/1000, 0)
		fmt.Fprintf(writer, "Time range: %s to %s\n", startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))
		fmt.Fprintf(writer, "Duration: %s\n\n", endTime.Sub(startTime).String())
	}

	hasData := false
	if data[0].CPU != nil {
		fmt.Fprintf(writer, "✓ CPU statistics available\n")
		hasData = true
	}
	if data[0].Memory != nil {
		fmt.Fprintf(writer, "✓ Memory statistics available\n")
		hasData = true
	}
	if data[0].Disk != nil {
		fmt.Fprintf(writer, "✓ Disk I/O statistics available\n")
		hasData = true
	}
	if data[0].Network != nil {
		fmt.Fprintf(writer, "✓ Network statistics available\n")
		hasData = true
	}

	if !hasData {
		fmt.Fprintf(writer, "⚠ No statistics data available\n")
	}

	fmt.Fprintf(writer, "\n")
	return nil
}

// GenerateCompactReport generates a compact summary for quick viewing
func (g *Generator) GenerateCompactReport(data []*SarData, writer io.Writer) error {
	if len(data) == 0 {
		return fmt.Errorf("no data available")
	}

	avgCPU := calculateAverageCPU(data)
	avgMem := calculateAverageMemory(data)

	fmt.Fprintf(writer, "=== Quick Summary ===\n")
	fmt.Fprintf(writer, "Samples: %d\n", len(data))

	if avgCPU != nil {
		fmt.Fprintf(writer, "\nCPU (avg):\n")
		fmt.Fprintf(writer, "  User:   %6.2f%%\n", avgCPU.User)
		fmt.Fprintf(writer, "  System: %6.2f%%\n", avgCPU.System)
		fmt.Fprintf(writer, "  IOWait: %6.2f%%\n", avgCPU.IOWait)
		fmt.Fprintf(writer, "  Idle:   %6.2f%%\n", avgCPU.Idle)
	}

	if avgMem != nil {
		fmt.Fprintf(writer, "\nMemory (avg):\n")
		fmt.Fprintf(writer, "  Used:   %6.2f%% (%12.0f KB)\n", avgMem.MemPercent, avgMem.MemUsed)
		fmt.Fprintf(writer, "  Free:   %12.0f KB\n", avgMem.MemFree)
		fmt.Fprintf(writer, "  Cached: %12.0f KB\n", avgMem.MemCached)
	}

	return nil
}

func calculateAverageCPU(data []*SarData) *CPUStats {
	var total CPUStats
	count := 0

	for _, d := range data {
		if d.CPU != nil {
			total.User += d.CPU.User
			total.Nice += d.CPU.Nice
			total.System += d.CPU.System
			total.IOWait += d.CPU.IOWait
			total.Steal += d.CPU.Steal
			total.Idle += d.CPU.Idle
			total.NumCPUs = d.CPU.NumCPUs
			count++
		}
	}

	if count == 0 {
		return nil
	}

	return &CPUStats{
		User:    total.User / float64(count),
		Nice:    total.Nice / float64(count),
		System:  total.System / float64(count),
		IOWait:  total.IOWait / float64(count),
		Steal:   total.Steal / float64(count),
		Idle:    total.Idle / float64(count),
		NumCPUs: total.NumCPUs,
	}
}

func calculateAverageMemory(data []*SarData) *MemoryStats {
	var total MemoryStats
	count := 0

	for _, d := range data {
		if d.Memory != nil {
			total.MemTotal += d.Memory.MemTotal
			total.MemUsed += d.Memory.MemUsed
			total.MemFree += d.Memory.MemFree
			total.MemCached += d.Memory.MemCached
			total.MemBuffers += d.Memory.MemBuffers
			total.MemPercent += d.Memory.MemPercent
			count++
		}
	}

	if count == 0 {
		return nil
	}

	return &MemoryStats{
		MemTotal:   total.MemTotal / float64(count),
		MemUsed:    total.MemUsed / float64(count),
		MemFree:    total.MemFree / float64(count),
		MemCached:  total.MemCached / float64(count),
		MemBuffers: total.MemBuffers / float64(count),
		MemPercent: total.MemPercent / float64(count),
	}
}
