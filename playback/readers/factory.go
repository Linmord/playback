package readers

import (
	"fmt"
	"io"
)

type ReaderFactory struct {
	config ReaderConfig
}

func NewFactory(config ReaderConfig) *ReaderFactory {
	fmt.Printf("Creating ReaderFactory: Buffer=%v, Stats=%v, Size=%dKB\n",
		config.EnableBuffer, config.EnableStats, config.BufferSize/1024)
	return &ReaderFactory{config: config}
}

func (rf *ReaderFactory) Create(reader io.ReadCloser) io.ReadCloser {
	result := reader

	fmt.Printf("Building reader chain: ")
	components := []string{}

	if rf.config.EnableBuffer {
		components = append(components, fmt.Sprintf("Buffered(%dKB)", rf.config.BufferSize/1024))
		result = NewBuffered(result, rf.config.BufferSize)
	} else {
		components = append(components, "PassThrough")
	}

	if rf.config.EnableStats {
		components = append(components, "Stats")
		result = NewStats(result)
	}

	if len(components) == 0 {
		components = append(components, "PassThrough")
	}

	fmt.Printf("%s\n", joinComponents(components))
	return result
}

func joinComponents(components []string) string {
	if len(components) == 1 {
		return components[0]
	}

	result := components[0]
	for i := 1; i < len(components); i++ {
		result += " -> " + components[i]
	}
	return result
}
