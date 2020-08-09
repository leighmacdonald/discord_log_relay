package relay

import (
	"bytes"
	"encoding/gob"
)

type MessageType int

const (
	TypeLog MessageType = iota
	TypeStartup
	TypeShutdown
)

type Payload struct {
	Type    MessageType
	Server  string
	Message string
}

func Encode(payload Payload) ([]byte, error) {
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(payload)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func Decode(b []byte, payload *Payload) error {
	return gob.NewDecoder(bytes.NewReader(b)).Decode(payload)
}

func init() {
	gob.Register(Payload{})
}
