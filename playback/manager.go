package playback

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"
)

// Run 默认配置运行
func Run() {
	defaultConfig := Config{
		SampleRate:   48000,
		Channels:     2,
		BufferSize:   64 * 1024,
		EnableBuffer: false,
		EnableStats:  true,
	}
	RunWithConfig(defaultConfig)
}

// RunWithConfig 自定义配置
func RunWithConfig(config Config) {
	client, err := NewAudioClient(config)
	if err != nil {
		fmt.Printf("Audio init failed: %v\n", err)
		return
	}
	defer client.Close()

	serverAddr := getValidServerAddress()
	if serverAddr == "" {
		fmt.Println("No valid address provided, exiting...")
		return
	}

	startPlaybackLoop(client, serverAddr)
}

func validateAddress(addr string) bool {
	// 空地址检查
	if addr == "" {
		return false
	}

	// 检查是否是 HTTP/HTTPS URL
	if strings.HasPrefix(addr, "http://") || strings.HasPrefix(addr, "https://") {
		_, err := url.Parse(addr)
		if err != nil {
			fmt.Printf("   Invalid URL: %v\n", err)
			return false
		}

		// 检查是否有路径（对于HTTP流）
		if !strings.Contains(addr, "/") {
			fmt.Println("   HTTP address should include path (e.g., /stream.wav)")
			return false
		}

		return true
	}

	// 检查 TCP 地址格式 (IP:端口 或 tcp://IP:端口)
	if strings.HasPrefix(addr, "tcp://") {
		addr = strings.TrimPrefix(addr, "tcp://")
	}

	// 验证 IP:端口 格式
	parts := strings.Split(addr, ":")
	if len(parts) != 2 {
		fmt.Println("   TCP address should be in format IP:PORT")
		return false
	}

	// 验证端口
	port := parts[1]
	if len(port) == 0 || len(port) > 5 {
		fmt.Println("   Port should be between 1-65535")
		return false
	}

	for _, ch := range port {
		if ch < '0' || ch > '9' {
			fmt.Println("   Port should contain only digits")
			return false
		}
	}

	portNum := 0
	fmt.Sscanf(port, "%d", &portNum)
	if portNum < 1 || portNum > 65535 {
		fmt.Println("   Port should be between 1-65535")
		return false
	}

	// 验证 IP 地址（简化验证）
	ip := parts[0]
	if ip != "localhost" {
		ipParts := strings.Split(ip, ".")
		if len(ipParts) != 4 {
			fmt.Println("   IP address should be in format X.X.X.X")
			return false
		}

		for _, part := range ipParts {
			if len(part) == 0 || len(part) > 3 {
				fmt.Println("   Each IP segment should be 1-3 digits")
				return false
			}
			for _, ch := range part {
				if ch < '0' || ch > '9' {
					fmt.Println("   IP should contain only digits and dots")
					return false
				}
			}
			num := 0
			fmt.Sscanf(part, "%d", &num)
			if num < 0 || num > 255 {
				fmt.Println("   IP segments should be between 0-255")
				return false
			}
		}
	}

	return true
}

func getValidServerAddress() string {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter server address: ")
		input, _ := reader.ReadString('\n')
		addr := strings.TrimSpace(input)

		if addr == "" {
			addr = "192.168.1.8:12345"
			fmt.Printf("Using default address: %s\n", addr)
			return addr
		}

		if validateAddress(addr) {
			return addr
		}

		fmt.Println("Invalid address format. Please try again.")
		fmt.Println("   Examples:")
		fmt.Println("   - TCP: 192.168.1.8:12345 or tcp://192.168.1.8:12345")
		fmt.Println("   - HTTP: http://192.168.1.8:8888/stream.wav")
		fmt.Println("   - HTTPS: https://example.com/audio/stream")
		fmt.Println()
	}
}

func startPlaybackLoop(client *AudioClient, serverAddr string) {
	reconnectAttempts := 0
	maxReconnectDelay := 60 * time.Second

	fmt.Printf("Starting playback: %s\n", serverAddr)
	fmt.Println("Press Ctrl+C to stop")

	for {
		fmt.Printf("Connecting to %s...\n", serverAddr)

		if err := client.ConnectAndPlay(serverAddr); err != nil {
			fmt.Printf("Connection failed: %v\n", err)
			reconnectAttempts++

			delay := time.Duration(reconnectAttempts) * 3 * time.Second
			if delay > maxReconnectDelay {
				delay = maxReconnectDelay
			}

			fmt.Printf("Retrying in %v (attempt %d)...\n", delay, reconnectAttempts)
			time.Sleep(delay)
			continue
		}

		fmt.Println("Connected successfully")
		reconnectAttempts = 0

		monitorPlayback(client)

		fmt.Println("Connection lost, reconnecting...")
		time.Sleep(2 * time.Second)
	}
}

func monitorPlayback(client *AudioClient) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if !client.IsPlaying() {
			fmt.Println("Playback stopped")
			return
		}
	}
}
