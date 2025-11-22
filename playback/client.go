package playback

import (
	"fmt"
	"io"
	"playback/playback/protocols"
	"playback/playback/readers"
	"time"

	"github.com/ebitengine/oto/v3"
)

type Config struct {
	SampleRate   int
	Channels     int
	BufferSize   int
	EnableStats  bool
	EnableBuffer bool
}

type AudioClient struct {
	config        Config
	context       *oto.Context
	player        *oto.Player
	stream        io.ReadCloser
	protocols     *protocols.ProtocolFactory
	readerFactory *readers.ReaderFactory
}

func NewAudioClient(config Config) (*AudioClient, error) {
	op := &oto.NewContextOptions{
		SampleRate:   config.SampleRate,
		ChannelCount: config.Channels,
		Format:       oto.FormatSignedInt16LE,
	}

	ctx, ready, err := oto.NewContext(op)
	if err != nil {
		return nil, err
	}
	<-ready

	// 配置读取器工厂
	readerConfig := readers.ReaderConfig{
		EnableBuffer: config.EnableBuffer,
		EnableStats:  config.EnableStats,
		BufferSize:   config.BufferSize,
	}

	return &AudioClient{
		config:        config,
		context:       ctx,
		protocols:     protocols.NewProtocolFactory(),
		readerFactory: readers.NewFactory(readerConfig),
	}, nil
}

func (ac *AudioClient) ConnectAndPlay(server string) error {
	ac.cleanup()

	protocol := ac.protocols.GetProtocol(server)
	stream, err := protocol.Connect(server)
	if err != nil {
		return err
	}

	ac.stream = ac.readerFactory.Create(stream)

	bufferStatus := "disabled"
	if ac.config.EnableBuffer {
		bufferStatus = fmt.Sprintf("enabled (%dKB)", ac.config.BufferSize/1024)
	}

	statsStatus := "disabled"
	if ac.config.EnableStats {
		statsStatus = "enabled"
	}

	fmt.Printf("Starting playback via %s | Buffer: %s | Stats: %s\n",
		protocol.Name(), bufferStatus, statsStatus)

	ac.player = ac.context.NewPlayer(ac.stream)
	ac.player.Play()

	return nil
}
func (ac *AudioClient) IsPlaying() bool {
	return ac.player != nil && ac.player.IsPlaying()
}

func (ac *AudioClient) cleanup() {
	if ac.player != nil {
		ac.player.Pause()
		time.Sleep(50 * time.Millisecond)
		ac.player.Close()
		ac.player = nil
	}
	if ac.stream != nil {
		ac.stream.Close()
		ac.stream = nil
	}
}

func (ac *AudioClient) Close() {
	ac.cleanup()
}

// RegisterProtocol 注册新协议
func (ac *AudioClient) RegisterProtocol(name string, protocol protocols.Protocol) {
	ac.protocols.Register(name, protocol)
}
