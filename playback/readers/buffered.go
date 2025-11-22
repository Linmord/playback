package readers

import (
	"fmt"
	"io"
	"time"
)

type BufferedReader struct {
	reader      io.ReadCloser
	buffer      []byte
	bufferSize  int
	pos         int
	end         int
	created     time.Time
	readCount   int
	lastLogTime time.Time
	totalFilled int64
	totalServed int64
}

func NewBuffered(reader io.ReadCloser, bufferSize int) io.ReadCloser {
	fmt.Printf("BufferedReader created with %dKB buffer\n", bufferSize/1024)
	return &BufferedReader{
		reader:      reader,
		buffer:      make([]byte, bufferSize),
		bufferSize:  bufferSize,
		created:     time.Now(),
		lastLogTime: time.Now(),
	}
}

func (br *BufferedReader) Read(p []byte) (int, error) {
	br.readCount++

	if br.pos < br.end {
		available := br.end - br.pos
		n := copy(p, br.buffer[br.pos:br.end])
		br.pos += n
		br.totalServed += int64(n)

		br.logBufferStatus("serving", n, available)
		return n, nil
	}

	n, err := br.reader.Read(br.buffer)
	if err != nil {
		fmt.Printf("BufferedReader read error: %v\n", err)
		return 0, err
	}

	br.pos = 0
	br.end = n
	br.totalFilled += int64(n)

	available := br.end - br.pos
	copied := copy(p, br.buffer[br.pos:br.end])
	br.pos += copied
	br.totalServed += int64(copied)

	br.logBufferStatus("filled", n, available)

	return copied, nil
}

func (br *BufferedReader) logBufferStatus(action string, bytes int, available int) {
	now := time.Now()

	if br.readCount%100 == 0 || now.Sub(br.lastLogTime) > 3*time.Second {
		bufferUsage := 0.0
		if action == "filled" {
			bufferUsage = float64(bytes) / float64(br.bufferSize) * 100
		} else {
			bufferUsage = float64(available) / float64(br.bufferSize) * 100
		}

		efficiency := 0.0
		if br.totalFilled > 0 {
			efficiency = float64(br.totalServed) / float64(br.totalFilled) * 100
		}

		fmt.Printf("BufferedReader: %s %dB, buffer %.1f%% used, efficiency %.1f%%, reads: %d\n",
			action, bytes, bufferUsage, efficiency, br.readCount)

		br.lastLogTime = now
	}
}

func (br *BufferedReader) Close() error {
	duration := time.Since(br.created)
	efficiency := 0.0
	if br.totalFilled > 0 {
		efficiency = float64(br.totalServed) / float64(br.totalFilled) * 100
	}

	fmt.Printf("BufferedReader closed after %v\n", duration.Round(time.Second))
	fmt.Printf("   Total: filled %s, served %s, efficiency %.1f%%, reads: %d\n",
		formatBytes(br.totalFilled), formatBytes(br.totalServed), efficiency, br.readCount)

	if br.reader != nil {
		return br.reader.Close()
	}
	return nil
}
