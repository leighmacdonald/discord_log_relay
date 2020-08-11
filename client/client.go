package client

import (
	"bytes"
	"context"
	"fmt"
	"github.com/hpcloud/tail"
	"github.com/leighmacdonald/discord_log_relay/relay"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"strings"
)

// maxBufferSize specifies the size of the buffers that
// are used to temporarily hold data from the UDP packets
// that we receive.
const maxBufferSize = 1024

func fileReader(path string, messageChan chan string) {
	t, err := tail.TailFile(path, tail.Config{Follow: true})
	if err != nil {
		log.Fatalf("Invalid log path: %s", path)
	}
	for line := range t.Lines {
		m := strings.TrimRight(line.Text, "\r\n")
		fmt.Println(m)
	}
}
func New(ctx context.Context, name string, logPath string, address string) (err error) {
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Errorf("Failed to close conn: %v", err)
		}
	}()
	messageChan := make(chan string)
	errChan := make(chan error)
	go fileReader(logPath, messageChan)
	select {
	case msg := <-messageChan:
		b, err2 := relay.Encode(relay.Payload{Type: relay.TypeLog, Server: name, Message: msg})
		if err2 != nil {
			log.Errorf("Error encoding payload")
			break
		}
		_, err2 = io.Copy(conn, bytes.NewReader(b))
		if err2 != nil {
			return
		}
	case <-ctx.Done():
		fmt.Println("cancelled")
		err = ctx.Err()
	case err = <-errChan:
		log.Fatalf("Fatal error occurred: %v", err)
	}
	return
}
