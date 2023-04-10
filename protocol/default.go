package protocol

import (
	"encoding/binary"
	"io"
	"net"
)

const (
	MessageIDLength    = 8
	ServiceNameLength  = 16
	FunctionNameLength = 16
	BodyLengthLength   = 4
)

type DefaultProtocol struct {
	MessageID    uint64 // 8 bytes
	ServiceName  string // 16 bytes
	FunctionName string // 16 bytes
	BodyLength   uint32 // 4 bytes
	Body         []byte // BodyLength bytes
}

var _ Protocol = &DefaultProtocol{}

func NewDefaultProtocol() *DefaultProtocol {
	return &DefaultProtocol{}
}

func (p *DefaultProtocol) Pack() []byte {
	bodyLen := len(p.Body)
	buffer := make([]byte, MessageIDLength+ServiceNameLength+FunctionNameLength+BodyLengthLength+bodyLen)

	p.BodyLength = uint32(bodyLen)

	binary.BigEndian.PutUint64(buffer[0:MessageIDLength], p.MessageID)
	copy(buffer[MessageIDLength:MessageIDLength+ServiceNameLength], p.ServiceName)
	copy(buffer[MessageIDLength+ServiceNameLength:MessageIDLength+ServiceNameLength+FunctionNameLength], p.FunctionName)
	binary.BigEndian.PutUint32(buffer[MessageIDLength+ServiceNameLength+FunctionNameLength:MessageIDLength+
		ServiceNameLength+FunctionNameLength+BodyLengthLength], p.BodyLength)
	copy(buffer[MessageIDLength+ServiceNameLength+FunctionNameLength+BodyLengthLength:], p.Body)

	return buffer
}

func (p *DefaultProtocol) UnPack(c net.Conn) error {
	var messageID = make([]byte, MessageIDLength)
	_, err := io.ReadFull(c, messageID)
	if err != nil {
		return err
	}
	p.MessageID = binary.BigEndian.Uint64(messageID)

	var serviceName = make([]byte, ServiceNameLength)
	_, err = io.ReadFull(c, serviceName)
	if err != nil {
		return err
	}
	p.ServiceName = string(serviceName)

	var functionName = make([]byte, FunctionNameLength)
	_, err = io.ReadFull(c, functionName)
	if err != nil {
		return err
	}
	p.FunctionName = string(functionName)

	var bodyLength = make([]byte, BodyLengthLength)
	_, err = io.ReadFull(c, bodyLength)
	if err != nil {
		return err
	}
	p.BodyLength = binary.BigEndian.Uint32(bodyLength)

	p.Body = make([]byte, p.BodyLength)
	_, err = io.ReadFull(c, p.Body)
	if err != nil {
		return err
	}

	return nil
}
