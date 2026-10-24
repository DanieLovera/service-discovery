package ap

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"

	"google.golang.org/protobuf/proto"

	appb "tpiii.local/daniel-tpiii/gen/ap"
)

const maxPacketSize = 64 * 1024

type udpTransport struct {
	bindAddress string
	logger      *slog.Logger

	conn *net.UDPConn
	wg   sync.WaitGroup
}

func newUDPTransport(bindAddress string, logger *slog.Logger) *udpTransport {
	return &udpTransport{
		bindAddress: bindAddress,
		logger:      logger,
	}
}

func (t *udpTransport) listen(handler func(*appb.Packet)) error {
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

func (t *udpTransport) send(address string, packet *appb.Packet) error {
	data, err := proto.Marshal(packet)
	if err != nil {
		return fmt.Errorf("encode packet: %w", err)
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

func (t *udpTransport) read(handler func(*appb.Packet)) {
	buffer := make([]byte, maxPacketSize)

	for {
		n, from, err := t.conn.ReadFromUDP(buffer)
		if errors.Is(err, net.ErrClosed) {
			return
		}

		if err != nil {
			t.logger.Warn("Failed to read peer packet", "error", err)
			continue
		}

		var packet appb.Packet
		if err := proto.Unmarshal(buffer[:n], &packet); err != nil {
			t.logger.Warn("Discarding malformed peer packet", "from", from.String(), "error", err)
			continue
		}

		handler(&packet)
	}
}

func (t *udpTransport) close() error {
	if t.conn == nil {
		return nil
	}

	err := t.conn.Close()
	t.wg.Wait()

	return err
}
