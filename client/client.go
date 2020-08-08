package client

import (
	"bytes"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"time"
)

var messageQueue = make(chan []byte)

// maxBufferSize specifies the size of the buffers that
// are used to temporarily hold data from the UDP packets
// that we receive.
const maxBufferSize = 1024

func New(ctx context.Context, address string, reader io.Reader) (err error) {
	raddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return
	}
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Errorf("Failed to close conn: %v", err)
		}
	}()
	doneChan := make(chan error, 1)
	go func() {
		n, err := io.Copy(conn, reader)
		if err != nil {
			doneChan <- err
			return
		}
		fmt.Printf("packet-written: bytes=%d\n", n)
		buffer := make([]byte, maxBufferSize)
		deadline := time.Now().Add(5 * time.Second)
		err = conn.SetReadDeadline(deadline)
		if err != nil {
			doneChan <- err
			return
		}

		nRead, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			doneChan <- err
			return
		}
		fmt.Printf("packet-received: bytes=%d from=%s\n", nRead, addr.String())
		doneChan <- nil
	}()

	select {
	case msg := <-messageQueue:
		n, err := io.Copy(conn, bytes.NewReader(msg))
		if err != nil {
			doneChan <- err
			return
		}
	case <-ctx.Done():
		fmt.Println("cancelled")
		err = ctx.Err()
	case err = <-doneChan:
	}

	return
}
