package discover

import (
	"net"
)

type Discover interface {
	Discover(listenPort uint16) (net.IP, uint16, error)
}
