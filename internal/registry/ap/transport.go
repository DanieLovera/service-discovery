package ap

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"

	"google.golang.org/protobuf/proto"

	gossippb "tpiii.local/daniel-tpiii/gen/gossip"
)

const maxPacketSize = 64 * 1024

type transport struct {
	bindAddress string
	logger      *slog.Logger

	conn *net.UDPConn
	wg   sync.WaitGroup
}

func newTransport(bindAddress string, logger *slog.Logger) *transport {
	return &transport{
		bindAddress: bindAddress,
		logger:      logger,
	}
}

func (t *transport) listen(handler func(*gossippb.Message)) error {
	address, err := net.ResolveUDPAddr("udp", t.bindAddress)
	if err != nil {
		return fmt.Errorf("resolve peer bind address: %w", err)
	}

	conn, err := net.ListenUDP("udp", address)
	if err != nil {
		return fmt.Errorf("listen on peer bind address: %w", err)
	}

	t.conn = conn
	t.wg.Go(func() {
		t.read(handler)
	})

	return nil
}

func (t *transport) send(address string, message *gossippb.Message) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("encode message: %w", err)
	}

	target, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return fmt.Errorf("resolve %s: %w", address, err)
	}

	if _, err := t.conn.WriteToUDP(data, target); err != nil {
		return fmt.Errorf("send to %s: %w", address, err)
	}

	return nil
}

func (t *transport) read(handler func(*gossippb.Message)) {
	buffer := make([]byte, maxPacketSize)

	for {
		n, from, err := t.conn.ReadFromUDP(buffer)
		if errors.Is(err, net.ErrClosed) {
			return
		}

		if err != nil {
			t.logger.Warn("Failed to read peer message", "error", err)
			continue
		}

		var message gossippb.Message
		if err := proto.Unmarshal(buffer[:n], &message); err != nil {
			t.logger.Warn("Discarding malformed peer message", "from", from.String(), "error", err)
			continue
		}

		handler(&message)
	}
}

func (t *transport) close() error {
	if t.conn == nil {
		return nil
	}

	err := t.conn.Close()
	t.wg.Wait()

	return err
}
