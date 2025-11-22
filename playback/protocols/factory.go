package protocols

import "strings"

// ProtocolFactory 协议工厂
type ProtocolFactory struct {
	protocols map[string]Protocol
}

func NewProtocolFactory() *ProtocolFactory {
	factory := &ProtocolFactory{
		protocols: make(map[string]Protocol),
	}

	// 注册内置协议
	factory.Register("tcp", &TCPProtocol{})
	factory.Register("http", &HTTPProtocol{})
	factory.Register("https", &HTTPProtocol{})

	return factory
}

// Register 注册新协议
func (f *ProtocolFactory) Register(name string, protocol Protocol) {
	f.protocols[name] = protocol
}

// GetProtocol 根据地址获取协议
func (f *ProtocolFactory) GetProtocol(server string) Protocol {
	if strings.HasPrefix(server, "http://") {
		return f.protocols["http"]
	}
	if strings.HasPrefix(server, "https://") {
		return f.protocols["https"]
	}
	if strings.HasPrefix(server, "tcp://") {
		return f.protocols["tcp"]
	}
	// 默认使用TCP
	return f.protocols["tcp"]
}
