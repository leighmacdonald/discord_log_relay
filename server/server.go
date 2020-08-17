package server

import (
	"context"
	"fmt"
	"github.com/leighmacdonald/discord_log_relay/consts"
	"github.com/leighmacdonald/discord_log_relay/relay"
	log "github.com/sirupsen/logrus"
	"net"
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
			var p relay.Payload
			if err := relay.Decode(buffer[:n], &p); err != nil {
				log.Errorf("failed to decode payload: %v", err)
				continue
			}
			_ = relay.SendPayload(p)
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
