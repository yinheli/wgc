package stun

import (
	"fmt"
	"net"
	"strconv"

	"github.com/pion/stun"
	"github.com/yinheli/wgc/netx"
)

// GetIPAndPort returns the public IP and port (and local source port) of the STUN server.
func GetIPAndPort(addr string) (net.IP, uint16, uint16, error) {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return nil, 0, 0, err
	}

	laddr := conn.LocalAddr()
	_, port, _ := net.SplitHostPort(laddr.String())
	sourcePort, _ := strconv.ParseUint(port, 10, 16)

	req := stun.MustBuild(stun.TransactionID, stun.BindingRequest)

	_, err = conn.Write(req.Raw)
	if err != nil {
		return nil, 0, 0, err
	}

	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		return nil, 0, 0, err
	}

	if n < 20 {
		return nil, 0, 0, fmt.Errorf("STUN response too short")
	}

	rsp := stun.MustBuild()
	rsp.Raw = data[:n]
	err = rsp.Decode()

	if err != nil {
		return nil, 0, 0, err
	}

	var xorAddr stun.XORMappedAddress
	err = xorAddr.GetFrom(rsp)
	if err != nil {
		return nil, 0, 0, err
	}

	return xorAddr.IP, uint16(xorAddr.Port), uint16(sourcePort), nil
}

// OpenTunnel opens a UDP tunnel (use srcPort) to the STUN server.
func OpenTunnel(addr string, srcPort uint16) error {
	stunAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	conn, err := netx.Dial(stunAddr.IP, srcPort, uint16(stunAddr.Port))
	if err != nil {
		return err
	}

	_, err = conn.Write(stun.MustBuild(stun.TransactionID, stun.BindingRequest).Raw)
	if err != nil {
		return err
	}

	return nil
}
