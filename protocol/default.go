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
	messageID    uint64 // 8 bytes
	serviceName  string // 16 bytes
	functionName string // 16 bytes
	bodyLength   uint32 // 4 bytes
	body         []byte // bodyLength bytes
}

var _ Protocol = &DefaultProtocol{}

func (p *DefaultProtocol) Pack() []byte {
	buffer := make([]byte, MessageIDLength+ServiceNameLength+FunctionNameLength+BodyLengthLength+len(p.body))

	binary.BigEndian.PutUint64(buffer[0:MessageIDLength], p.messageID)
	copy(buffer[MessageIDLength:MessageIDLength+ServiceNameLength], p.serviceName)
	copy(buffer[MessageIDLength+ServiceNameLength:MessageIDLength+ServiceNameLength+FunctionNameLength], p.functionName)
	binary.BigEndian.PutUint32(buffer[MessageIDLength+ServiceNameLength+FunctionNameLength:MessageIDLength+
		ServiceNameLength+FunctionNameLength+BodyLengthLength], p.bodyLength)
	copy(buffer[MessageIDLength+ServiceNameLength+FunctionNameLength+BodyLengthLength:], p.body)

	return buffer
}

func (p *DefaultProtocol) UnPack(c net.Conn) error {
	var messageID = make([]byte, MessageIDLength)
	_, err := io.ReadFull(c, messageID)
	if err != nil {
		return err
	}
	p.messageID = binary.BigEndian.Uint64(messageID)

	var serviceName = make([]byte, ServiceNameLength)
	_, err = io.ReadFull(c, serviceName)
	if err != nil {
		return err
	}
	p.serviceName = string(serviceName)

	var functionName = make([]byte, FunctionNameLength)
	_, err = io.ReadFull(c, functionName)
	if err != nil {
		return err
	}
	p.functionName = string(functionName)

	var bodyLength = make([]byte, BodyLengthLength)
	_, err = io.ReadFull(c, bodyLength)
	if err != nil {
		return err
	}
	p.bodyLength = binary.BigEndian.Uint32(bodyLength)

	p.body = make([]byte, p.bodyLength)
	_, err = io.ReadFull(c, p.body)
	if err != nil {
		return err
	}
	
	// p.messageID = binary.BigEndian.Uint64(data[0:MessageIDLength])
	// p.serviceName = string(data[MessageIDLength : MessageIDLength+ServiceNameLength])
	// p.functionName = string(data[MessageIDLength+ServiceNameLength : MessageIDLength+ServiceNameLength+FunctionNameLength])
	// p.bodyLength = binary.BigEndian.Uint32(data[MessageIDLength+ServiceNameLength+
	// 	FunctionNameLength : MessageIDLength+ServiceNameLength+FunctionNameLength+BodyLengthLength])
	// p.body = data[MessageIDLength+ServiceNameLength+FunctionNameLength+BodyLengthLength:]
	return nil
}
