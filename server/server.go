package server

import (
	"context"
	"fmt"
	"github.com/leighmacdonald/discord_log_relay/consts"
	"github.com/leighmacdonald/discord_log_relay/relay"
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

func Server(ctx context.Context, address string) (err error) {
	pc, err := net.ListenPacket("udp", address)
	if err != nil {
		return
	}
	defer func() {
		if err := pc.Close(); err != nil {
			log.Errorf("Failed to close client conn: %v", err)
		}
	}()
	doneChan := make(chan error, 1)
	buffer := make([]byte, consts.MaxBufferSize)
	go func() {
		for {
			n, addr, err := pc.ReadFrom(buffer)
			if err != nil {
				doneChan <- err
				return
			}

			fmt.Printf("packet-received: bytes=%d from=%s\n",
				n, addr.String())
			deadline := time.Now().Add(consts.Timeout)
			err = pc.SetWriteDeadline(deadline)
			if err != nil {
				doneChan <- err
				return
			}
			var p relay.Payload
			relay.Decode(buffer[:n], &p)
			// Write the packet's contents back to the client.
			n, err = pc.WriteTo(buffer[:n], addr)
			if err != nil {
				log.Errorf("failed to write packed to client: %v", err)
				// doneChan <- err
				return
			}

			fmt.Printf("packet-written: bytes=%d to=%s\n", n, addr.String())
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("cancelled")
		err = ctx.Err()
	case err = <-doneChan:
	}

	return
}
