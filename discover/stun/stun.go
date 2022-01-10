package stun

import (
	"net"

	"github.com/yinheli/wgc/config/env"
	"github.com/yinheli/wgc/discover"
	"github.com/yinheli/wgc/stun"
)

const (
	envNamespace = "STUN_"

	ENVServer = envNamespace + "SERVER"
)

type Config struct {
	Server string
}

type STUNDiscover struct {
	config Config
}

func NewDiscover() discover.Discover {
	instance := &STUNDiscover{
		config: Config{
			Server: env.GetOrDefaultString(ENVServer, "stun.miwifi.com:3478"),
		},
	}
	return instance
}

func (t *STUNDiscover) Discover(listenPort uint16) (net.IP, uint16, error) {
	ip, _, _, err := stun.GetIPAndPort(t.config.Server)
	if err != nil {
		return nil, 0, err
	}

	err = stun.OpenTunnel(t.config.Server, listenPort)
	if err != nil {
		return nil, 0, err
	}
	return ip, listenPort, nil
}
