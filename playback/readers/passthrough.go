package readers

import "io"

type PassThroughReader struct {
	reader io.ReadCloser
}

func NewPassThrough(reader io.ReadCloser) io.ReadCloser {
	return &PassThroughReader{reader: reader}
}

func (ptr *PassThroughReader) Read(p []byte) (int, error) {
	return ptr.reader.Read(p)
}

func (ptr *PassThroughReader) Close() error {
	if ptr.reader != nil {
		return ptr.reader.Close()
	}
	return nil
}
