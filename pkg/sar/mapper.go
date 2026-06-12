package sar

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/prometheus/prometheus-dump-operator/pkg/tsdb"
)

// MetricsMapper maps Prometheus metrics to sar data structures
type MetricsMapper struct {
	reader *tsdb.Reader
}

// NewMetricsMapper creates a new metrics mapper
func NewMetricsMapper(reader *tsdb.Reader) *MetricsMapper {
	return &MetricsMapper{reader: reader}
}

// SarData represents aggregated system metrics in sar format
type SarData struct {
	Timestamp int64
	CPU       *CPUStats
	Memory    *MemoryStats
	Disk      *DiskStats
	Network   *NetworkStats
}

// CPUStats represents CPU utilization (sar -u)
type CPUStats struct {
	User    float64 // %user
	Nice    float64 // %nice
	System  float64 // %system
	IOWait  float64 // %iowait
	Steal   float64 // %steal
	Idle    float64 // %idle
	NumCPUs int     // number of CPUs
}

// MemoryStats represents memory utilization (sar -r)
type MemoryStats struct {
	MemTotal     float64 // kbmemfree
	MemUsed      float64 // kbmemused
	MemFree      float64 // kbmemfree
	MemCached    float64 // kbcached
	MemBuffers   float64 // kbbuffers
	SwapTotal    float64 // kbswpfree
	SwapUsed     float64 // kbswpused
	SwapFree     float64 // kbswpfree
	MemPercent   float64 // %memused
	SwapPercent  float64 // %swpused
}

// DiskStats represents disk I/O (sar -d)
type DiskStats struct {
	Devices []DiskDevice
}

type DiskDevice struct {
	Device string
	TPS    float64 // transfers per second
	RdKBs  float64 // read KB/s
	WrKBs  float64 // write KB/s
	AvgRqSz float64 // average request size
	AvgQuSz float64 // average queue size
	Await   float64 // average wait time
	Util    float64 // utilization %
}

// NetworkStats represents network I/O (sar -n DEV)
type NetworkStats struct {
	Interfaces []NetworkInterface
}

type NetworkInterface struct {
	Interface string
	RxPckS    float64 // received packets/s
	TxPckS    float64 // transmitted packets/s
	RxKBS     float64 // received KB/s
	TxKBS     float64 // transmitted KB/s
	RxErrS    float64 // received errors/s
	TxErrS    float64 // transmitted errors/s
}

// ExtractSarData extracts sar-compatible data for a time range
func (m *MetricsMapper) ExtractSarData(ctx context.Context, startTime, endTime int64, interval int64, profile string) ([]*SarData, error) {
	var data []*SarData

	for ts := startTime; ts <= endTime; ts += interval * 1000 {
		sarData := &SarData{
			Timestamp: ts,
		}

		var err error

		if profile == "all" || profile == "cpu" {
			sarData.CPU, err = m.extractCPUStats(ctx, ts, ts+interval*1000)
			if err != nil {
				return nil, fmt.Errorf("failed to extract CPU stats: %w", err)
			}
		}

		if profile == "all" || profile == "memory" {
			sarData.Memory, err = m.extractMemoryStats(ctx, ts, ts+interval*1000)
			if err != nil {
				return nil, fmt.Errorf("failed to extract memory stats: %w", err)
			}
		}

		if profile == "all" || profile == "disk" {
			sarData.Disk, err = m.extractDiskStats(ctx, ts, ts+interval*1000)
			if err != nil {
				return nil, fmt.Errorf("failed to extract disk stats: %w", err)
			}
		}

		if profile == "all" || profile == "network" {
			sarData.Network, err = m.extractNetworkStats(ctx, ts, ts+interval*1000)
			if err != nil {
				return nil, fmt.Errorf("failed to extract network stats: %w", err)
			}
		}

		data = append(data, sarData)
	}

	return data, nil
}

// extractCPUStats extracts CPU statistics from Prometheus metrics
func (m *MetricsMapper) extractCPUStats(ctx context.Context, startTime, endTime int64) (*CPUStats, error) {
	cpuMetrics := []string{"user", "nice", "system", "iowait", "steal", "idle"}
	stats := &CPUStats{}

	for _, mode := range cpuMetrics {
		metric := "node_cpu_seconds_total"
		filters := map[string]string{"mode": mode}

		series, err := m.reader.QueryMetric(ctx, metric, startTime, endTime, filters)
		if err != nil {
			return nil, err
		}

		total := 0.0
		count := 0

		for _, ts := range series {
			if len(ts.Samples) > 0 {
				lastIdx := len(ts.Samples) - 1
				if lastIdx > 0 {
					rate := (ts.Samples[lastIdx].Value - ts.Samples[0].Value) /
						float64(ts.Samples[lastIdx].Timestamp-ts.Samples[0].Timestamp) * 1000 * 100
					total += rate
					count++
				}
			}
		}

		stats.NumCPUs = count

		avg := 0.0
		if count > 0 {
			avg = total / float64(count)
		}

		switch mode {
		case "user":
			stats.User = avg
		case "nice":
			stats.Nice = avg
		case "system":
			stats.System = avg
		case "iowait":
			stats.IOWait = avg
		case "steal":
			stats.Steal = avg
		case "idle":
			stats.Idle = avg
		}
	}

	return stats, nil
}

// extractMemoryStats extracts memory statistics
func (m *MetricsMapper) extractMemoryStats(ctx context.Context, startTime, endTime int64) (*MemoryStats, error) {
	stats := &MemoryStats{}

	memTotal, err := m.getMetricAverage(ctx, "node_memory_MemTotal_bytes", startTime, endTime, nil)
	if err != nil {
		return nil, err
	}
	stats.MemTotal = memTotal / 1024

	memFree, err := m.getMetricAverage(ctx, "node_memory_MemFree_bytes", startTime, endTime, nil)
	if err != nil {
		return nil, err
	}
	stats.MemFree = memFree / 1024

	memCached, err := m.getMetricAverage(ctx, "node_memory_Cached_bytes", startTime, endTime, nil)
	if err != nil {
		return nil, err
	}
	stats.MemCached = memCached / 1024

	memBuffers, err := m.getMetricAverage(ctx, "node_memory_Buffers_bytes", startTime, endTime, nil)
	if err != nil {
		return nil, err
	}
	stats.MemBuffers = memBuffers / 1024

	stats.MemUsed = stats.MemTotal - stats.MemFree - stats.MemCached - stats.MemBuffers

	if stats.MemTotal > 0 {
		stats.MemPercent = (stats.MemUsed / stats.MemTotal) * 100
	}

	swapTotal, err := m.getMetricAverage(ctx, "node_memory_SwapTotal_bytes", startTime, endTime, nil)
	if err == nil {
		stats.SwapTotal = swapTotal / 1024

		swapFree, err := m.getMetricAverage(ctx, "node_memory_SwapFree_bytes", startTime, endTime, nil)
		if err == nil {
			stats.SwapFree = swapFree / 1024
			stats.SwapUsed = stats.SwapTotal - stats.SwapFree

			if stats.SwapTotal > 0 {
				stats.SwapPercent = (stats.SwapUsed / stats.SwapTotal) * 100
			}
		}
	}

	return stats, nil
}

// extractDiskStats extracts disk I/O statistics
func (m *MetricsMapper) extractDiskStats(ctx context.Context, startTime, endTime int64) (*DiskStats, error) {
	stats := &DiskStats{
		Devices: []DiskDevice{},
	}

	deviceSeries, err := m.reader.QueryMetric(ctx, "node_disk_io_time_seconds_total", startTime, endTime, nil)
	if err != nil {
		return nil, err
	}

	deviceMap := make(map[string]*DiskDevice)

	for _, ts := range deviceSeries {
		device := ""
		for _, label := range ts.Labels {
			if label.Name == "device" {
				device = label.Value
				break
			}
		}

		if device == "" {
			continue
		}

		if _, exists := deviceMap[device]; !exists {
			deviceMap[device] = &DiskDevice{Device: device}
		}
	}

	for device, diskDev := range deviceMap {
		filters := map[string]string{"device": device}

		readBytes, _ := m.getMetricRate(ctx, "node_disk_read_bytes_total", startTime, endTime, filters)
		diskDev.RdKBs = readBytes / 1024

		writeBytes, _ := m.getMetricRate(ctx, "node_disk_written_bytes_total", startTime, endTime, filters)
		diskDev.WrKBs = writeBytes / 1024

		readOps, _ := m.getMetricRate(ctx, "node_disk_reads_completed_total", startTime, endTime, filters)
		writeOps, _ := m.getMetricRate(ctx, "node_disk_writes_completed_total", startTime, endTime, filters)
		diskDev.TPS = readOps + writeOps

		if diskDev.TPS > 0 {
			diskDev.AvgRqSz = ((diskDev.RdKBs + diskDev.WrKBs) * 2) / diskDev.TPS
		}

		ioTime, _ := m.getMetricRate(ctx, "node_disk_io_time_seconds_total", startTime, endTime, filters)
		diskDev.Util = math.Min(ioTime*100, 100)

		stats.Devices = append(stats.Devices, *diskDev)
	}

	sort.Slice(stats.Devices, func(i, j int) bool {
		return stats.Devices[i].Device < stats.Devices[j].Device
	})

	return stats, nil
}

// extractNetworkStats extracts network I/O statistics
func (m *MetricsMapper) extractNetworkStats(ctx context.Context, startTime, endTime int64) (*NetworkStats, error) {
	stats := &NetworkStats{
		Interfaces: []NetworkInterface{},
	}

	ifaceSeries, err := m.reader.QueryMetric(ctx, "node_network_receive_bytes_total", startTime, endTime, nil)
	if err != nil {
		return nil, err
	}

	ifaceMap := make(map[string]*NetworkInterface)

	for _, ts := range ifaceSeries {
		iface := ""
		for _, label := range ts.Labels {
			if label.Name == "device" {
				iface = label.Value
				break
			}
		}

		if iface == "" || iface == "lo" {
			continue
		}

		if _, exists := ifaceMap[iface]; !exists {
			ifaceMap[iface] = &NetworkInterface{Interface: iface}
		}
	}

	for iface, netIface := range ifaceMap {
		filters := map[string]string{"device": iface}

		rxBytes, _ := m.getMetricRate(ctx, "node_network_receive_bytes_total", startTime, endTime, filters)
		netIface.RxKBS = rxBytes / 1024

		txBytes, _ := m.getMetricRate(ctx, "node_network_transmit_bytes_total", startTime, endTime, filters)
		netIface.TxKBS = txBytes / 1024

		rxPackets, _ := m.getMetricRate(ctx, "node_network_receive_packets_total", startTime, endTime, filters)
		netIface.RxPckS = rxPackets

		txPackets, _ := m.getMetricRate(ctx, "node_network_transmit_packets_total", startTime, endTime, filters)
		netIface.TxPckS = txPackets

		rxErrs, _ := m.getMetricRate(ctx, "node_network_receive_errs_total", startTime, endTime, filters)
		netIface.RxErrS = rxErrs

		txErrs, _ := m.getMetricRate(ctx, "node_network_transmit_errs_total", startTime, endTime, filters)
		netIface.TxErrS = txErrs

		stats.Interfaces = append(stats.Interfaces, *netIface)
	}

	sort.Slice(stats.Interfaces, func(i, j int) bool {
		return stats.Interfaces[i].Interface < stats.Interfaces[j].Interface
	})

	return stats, nil
}

// getMetricAverage gets the average value of a metric
func (m *MetricsMapper) getMetricAverage(ctx context.Context, metric string, startTime, endTime int64, filters map[string]string) (float64, error) {
	series, err := m.reader.QueryMetric(ctx, metric, startTime, endTime, filters)
	if err != nil {
		return 0, err
	}

	total := 0.0
	count := 0

	for _, ts := range series {
		for _, sample := range ts.Samples {
			total += sample.Value
			count++
		}
	}

	if count == 0 {
		return 0, nil
	}

	return total / float64(count), nil
}

// getMetricRate calculates the rate of a counter metric
func (m *MetricsMapper) getMetricRate(ctx context.Context, metric string, startTime, endTime int64, filters map[string]string) (float64, error) {
	series, err := m.reader.QueryMetric(ctx, metric, startTime, endTime, filters)
	if err != nil {
		return 0, err
	}

	totalRate := 0.0
	count := 0

	for _, ts := range series {
		if len(ts.Samples) < 2 {
			continue
		}

		first := ts.Samples[0]
		last := ts.Samples[len(ts.Samples)-1]

		timeDiff := float64(last.Timestamp-first.Timestamp) / 1000.0
		if timeDiff > 0 {
			rate := (last.Value - first.Value) / timeDiff
			totalRate += rate
			count++
		}
	}

	if count == 0 {
		return 0, nil
	}

	return totalRate / float64(count), nil
}
