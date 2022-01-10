package wg

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	reSpace = regexp.MustCompile(`\s+`)
)

func GetIfacePubKey(iface string) (Key, error) {
	r, err := run("wg", "show", iface, "public-key")
	if err != nil {
		return DefaultKey, err
	}
	return NewKeyFromStr(r), nil
}

func GetIfaceListenPort(iface string) (uint16, error) {
	r, err := run("wg", "show", iface, "listen-port")
	if err != nil {
		return 0, err
	}
	v, err := strconv.ParseUint(r, 10, 16)
	if err != nil {
		return 0, err
	}
	return uint16(v), nil
}

func GetEndpoints(iface string) (map[Key]string, error) {
	r, err := run("wg", "show", iface, "endpoints")
	if err != nil {
		return nil, err
	}
	peers := make(map[Key]string, 128)
	for _, it := range strings.Split(r, "\n") {
		it = strings.TrimSpace(it)
		if it == "" {
			continue
		}
		arr := reSpace.Split(it, -1)
		if len(arr) < 2 {
			continue
		}
		peer := NewKeyFromStr(arr[0])
		endpoint := arr[1]
		if endpoint == "(none)" {
			endpoint = ""
		}
		peers[peer] = endpoint
	}
	return peers, nil
}

func SetListenPort(iface string, port uint16) error {
	_, err := run(
		"wg", "set", iface,
		"listen-port", fmt.Sprint(port),
	)
	if err != nil {
		return err
	}
	return nil
}

func SetPeerEndpoint(iface string, peer Key, endpoint string) error {
	_, err := run(
		"wg", "set", iface,
		"peer", base64.StdEncoding.EncodeToString(peer[:]),
		"persistent-keepalive", "10",
		"endpoint", endpoint,
	)
	if err != nil {
		return err
	}
	return nil
}
