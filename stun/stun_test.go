package stun

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_GetIPAndPort(t *testing.T) {
	ip, port, srcPort, err := GetIPAndPort("stun.miwifi.com:3478")
	require.NoError(t, err)
	fmt.Println(ip, port, srcPort)
}
