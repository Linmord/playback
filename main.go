package main

import "playback/playback"

func main() {
	config := playback.Config{
		SampleRate:   48000,
		Channels:     2,
		BufferSize:   10 * 1024,
		EnableBuffer: true,  // 禁用缓冲
		EnableStats:  false, // 禁用统计
	}
	playback.RunWithConfig(config)

}
