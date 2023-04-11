package protocol

import "net"

// Protocol is the interface for protocol
type Protocol interface {
	UnPack(net.Conn) error
	Pack() []byte
}
