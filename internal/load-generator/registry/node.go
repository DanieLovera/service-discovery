package registry

import (
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/credentials/insecure"

	registrypb "tpiii.local/daniel-tpiii/gen/registry"
)

const (
	reconnectBaseDelay = 100 * time.Millisecond
	reconnectMaxDelay  = time.Second
)

type Node struct {
	Name   string
	conn   *grpc.ClientConn
	client registrypb.RegistryClient
}

func Connect(address string) (*Node, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithConnectParams(grpc.ConnectParams{
			Backoff: backoff.Config{
				BaseDelay:  reconnectBaseDelay,
				Multiplier: backoff.DefaultConfig.Multiplier,
				Jitter:     backoff.DefaultConfig.Jitter,
				MaxDelay:   reconnectMaxDelay,
			},
			MinConnectTimeout: reconnectMaxDelay,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("create registry connection to %s: %w", address, err)
	}

	return &Node{
		Name:   nodeName(address),
		conn:   conn,
		client: registrypb.NewRegistryClient(conn),
	}, nil
}

func (n *Node) Close() error {
	if err := n.conn.Close(); err != nil {
		return fmt.Errorf("close registry connection to %s: %w", n.Name, err)
	}

	return nil
}

func nodeName(address string) string {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return address
	}

	return host
}
