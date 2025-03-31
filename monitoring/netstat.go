package monitoring

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/net"
	"time"
)

// NetworkStats 表示网络带宽统计数据
type NetworkStats struct {
	BytesSent       uint64  `json:"bytes_sent"`        // 发送的总字节数
	BytesRecv       uint64  `json:"bytes_recv"`        // 接收的总字节数
	UploadRate      float64 `json:"upload_rate"`       // 上传速率 (KB/s)
	DownloadRate    float64 `json:"download_rate"`     // 下载速率 (KB/s)
	UploadRateRaw   float64 `json:"upload_rate_raw"`   // 上传速率原始值 (Bytes/s)
	DownloadRateRaw float64 `json:"download_rate_raw"` // 下载速率原始值 (Bytes/s)
	Timestamp       int64   `json:"timestamp"`         // 时间戳
}

// NetworkMonitor 网络监控器
type NetworkMonitor struct {
	prevBytesSent uint64    // 上一次采样的发送字节数
	prevBytesRecv uint64    // 上一次采样的接收字节数
	prevTime      time.Time // 上一次采样的时间
}

// NewNetworkMonitor 创建一个新的网络监控器
// interval: 采样间隔时间
func NewNetworkMonitor() *NetworkMonitor {
	return &NetworkMonitor{}
}

// sample 采样网络统计数据
func (nm *NetworkMonitor) sample() (NetworkStats, error) {
	// 获取当前网络IO计数器
	counters, err := net.IOCounters(false)
	if err != nil {
		return NetworkStats{}, fmt.Errorf("failed to get network IO counters: %v", err)
	}

	// 汇总所有网络接口的统计数据
	var totalBytesSent, totalBytesRecv uint64
	for _, counter := range counters {
		totalBytesSent += counter.BytesSent
		totalBytesRecv += counter.BytesRecv
	}

	now := time.Now()
	timeElapsed := now.Sub(nm.prevTime).Seconds()

	// 计算速率
	var uploadRate, downloadRate float64
	if timeElapsed > 0 {
		bytesSentDiff := totalBytesSent - nm.prevBytesSent
		bytesRecvDiff := totalBytesRecv - nm.prevBytesRecv

		uploadRateRaw := float64(0)
		downloadRateRaw := float64(0)
		if !nm.prevTime.IsZero() {
			uploadRateRaw = float64(bytesSentDiff) / timeElapsed
			downloadRateRaw = float64(bytesRecvDiff) / timeElapsed
			uploadRate = uploadRateRaw / 1024     // 转换为KB/s
			downloadRate = downloadRateRaw / 1024 // 转换为KB/s
		}

		// 更新上一次的采样数据
		nm.prevBytesSent = totalBytesSent
		nm.prevBytesRecv = totalBytesRecv
		nm.prevTime = now

		return NetworkStats{
			BytesSent:       totalBytesSent,
			BytesRecv:       totalBytesRecv,
			UploadRate:      uploadRate,
			DownloadRate:    downloadRate,
			UploadRateRaw:   uploadRateRaw,
			DownloadRateRaw: downloadRateRaw,
			Timestamp:       now.Unix(),
		}, nil
	}

	return NetworkStats{
		BytesSent: totalBytesSent,
		BytesRecv: totalBytesRecv,
		Timestamp: now.Unix(),
	}, nil
}

// GetCurrentStats 获取当前网络统计信息（单次采样）
func (nm *NetworkMonitor) GetCurrentStats() (NetworkStats, error) {

	stats, err := nm.sample()
	if err != nil {
		return NetworkStats{}, err
	}

	return stats, nil
}
