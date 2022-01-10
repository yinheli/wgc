package register

import (
	"net"

	"github.com/yinheli/wgc/wg"
)

// Register
type Register interface {
	Register(key wg.Key, ip net.IP, port uint16) error
	Receive() (<-chan wg.Peer, error)
	Close() error
}
