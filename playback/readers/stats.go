package readers

import (
	"fmt"
	"io"
	"time"
)

type StatsReader struct {
	reader    io.ReadCloser
	stats     *StreamStats
	startTime time.Time
	lastPrint time.Time
}

type StreamStats struct {
	BytesTransferred int64
	ReadOperations   int64
	Errors           int64
	CurrentBitrate   float64
	AverageBitrate   float64
}

func NewStats(reader io.ReadCloser) *StatsReader {
	return &StatsReader{
		reader:    reader,
		stats:     &StreamStats{},
		startTime: time.Now(),
		lastPrint: time.Now(),
	}
}

func (sr *StatsReader) Read(p []byte) (int, error) {
	start := time.Now()
	n, err := sr.reader.Read(p)

	sr.stats.BytesTransferred += int64(n)
	sr.stats.ReadOperations++

	if err != nil {
		sr.stats.Errors++
	}

	// 计算实时码率
	duration := time.Since(start)
	if duration > 0 {
		sr.stats.CurrentBitrate = float64(n*8) / duration.Seconds() / 1000 // kbps
	}

	// 计算平均码率
	totalDuration := time.Since(sr.startTime)
	if totalDuration > 0 {
		sr.stats.AverageBitrate = float64(sr.stats.BytesTransferred*8) / totalDuration.Seconds() / 1000
	}

	// 打印统计信息
	if time.Since(sr.lastPrint) > 5*time.Second {
		sr.PrintStats()
		sr.lastPrint = time.Now()
	}

	return n, err
}

func (sr *StatsReader) PrintStats() {
	fmt.Printf("Stats: %d reads, %s transferred, %.1f kbps avg\n",
		sr.stats.ReadOperations,
		formatBytes(sr.stats.BytesTransferred),
		sr.stats.AverageBitrate)
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func (sr *StatsReader) Close() error {
	// 关闭前打印最终统计
	sr.PrintStats()
	if sr.reader != nil {
		return sr.reader.Close()
	}
	return nil
}

func (sr *StatsReader) GetStats() *StreamStats {
	return sr.stats
}
