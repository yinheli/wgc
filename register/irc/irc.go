package irc

import (
	"encoding/base64"
	"fmt"
	"net"
	"strings"

	"github.com/google/uuid"
	irc "github.com/thoj/go-ircevent"
	"github.com/yinheli/wgc/config/env"
	"github.com/yinheli/wgc/register"
	"github.com/yinheli/wgc/wg"
)

const (
	envNamespace = "IRC_"

	ENVServer  = envNamespace + "SERVER"
	ENVTls     = envNamespace + "TLS"
	ENVChannel = envNamespace + "CHANNEL"
)

type Config struct {
	Server  string
	Tls     bool
	Channel string
}

// IRCRegister implements register.Register
type IRCRegister struct {
	config Config
	c      *irc.Connection
	ch     chan wg.Peer
}

func NewRegister() register.Register {
	instance := &IRCRegister{
		config: Config{
			Server:  env.GetOrDefaultString(ENVServer, "irc.freenode.net:6667"),
			Tls:     env.GetOrDefaultBool(ENVServer, false),
			Channel: "#" + strings.TrimPrefix(env.GetOrDefaultString(ENVChannel, "#pub-wgc"), "#"),
		},
		ch: make(chan wg.Peer, 128),
	}

	username := strings.ReplaceAll(uuid.New().String(), "-", "")
	nickname := fmt.Sprint("wgc", username[:8])
	instance.c = irc.IRC(nickname, username)
	// instance.c.Debug = true
	instance.c.UseTLS = instance.config.Tls
	instance.c.Connect(instance.config.Server)

	instance.c.AddCallback("001", func(e *irc.Event) {
		instance.c.Join(instance.config.Channel)
	})

	instance.c.AddCallback("PRIVMSG", func(e *irc.Event) {
		data, err := base64.StdEncoding.DecodeString(e.Message())
		if err != nil {
			return
		}

		arr := strings.Split(string(data), "|")
		if len(arr) != 2 {
			return
		}

		addr, err := net.ResolveUDPAddr("udp", arr[1])
		if err != nil {
			return
		}

		key := wg.NewKeyFromStr(arr[0])

		peer := wg.NewPeerFromAddr(key, addr)

		instance.ch <- peer
	})

	go func() {
		instance.c.Loop()
	}()

	return instance
}

func (t *IRCRegister) Register(key wg.Key, ip net.IP, port uint16) error {
	if t.c.Connected() {
		info := fmt.Sprint(key.String(), "|", ip.String(), ":", port)
		t.c.Privmsg(t.config.Channel, base64.StdEncoding.EncodeToString([]byte(info)))
	}

	return nil
}

func (t *IRCRegister) Receive() (<-chan wg.Peer, error) {
	return t.ch, nil
}

func (t *IRCRegister) Close() error {
	t.c.Quit()
	t.c.Disconnect()
	close(t.ch)
	return nil
}
