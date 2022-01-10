package irc

import (
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/yinheli/wgc/wg"
)

func Test_reg_receive(t *testing.T) {
	reg := NewRegister()
	time.Sleep(time.Second * 5)

	key := wg.NewKeyFromStr("oKe/b5TU5rUoHBsUywRZUVdIZnmh0qt9sqWP/rzL1EQ=")

	go func() {
		for i := 0; i < 10; i++ {
			err := reg.Register(key, net.IPv4(127, 0, 0, 1), 1234)
			require.NoError(t, err)
			time.Sleep(time.Second)
		}

		reg.Close()
	}()

	ch, err := reg.Receive()
	require.NoError(t, err)
	for peer := range ch {
		t.Log("peer:", peer)
		key, addr := peer.Parse()
		assert.NotEmpty(t, key)
		assert.NotEmpty(t, addr)
	}
}
