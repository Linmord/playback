package readers

import "io"

type ReaderConfig struct {
	EnableBuffer bool
	EnableStats  bool
	BufferSize   int
}

type ReadCloser interface {
	io.ReadCloser
}
