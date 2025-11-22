package protocols

import "io"

// Protocol 协议接口
type Protocol interface {
	Connect(server string) (io.ReadCloser, error)
	Name() string
}
