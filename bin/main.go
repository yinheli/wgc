package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/yinheli/wgc"
	"github.com/yinheli/wgc/discover"
	"github.com/yinheli/wgc/discover/stun"
	"github.com/yinheli/wgc/register"
	"github.com/yinheli/wgc/register/irc"
	"github.com/yinheli/wgc/wg"
)

var (
	l = log.New(os.Stdout, "", log.LstdFlags)

	// flags
	iface        = flag.String("iface", "", "wireguard interface")
	listenPort   = flag.Uint("listen", 0, "wireguard listen port")
	publicKey    = flag.String("pubKey", "", "wireguard listen port")
	discoverType = flag.String("discover", "stun", "discover type")
	registerType = flag.String("register", "irc", "register type")
	version      = flag.Bool("version", false, "show version")
)

func main() {
	if flag.Parse(); !flag.Parsed() {
		flag.Usage()
		os.Exit(1)
	}

	if *version {
		fmt.Println(wgc.Version)
		os.Exit(0)
	}

	if *iface == "" {
		if *publicKey == "" || *listenPort == 0 {
			l.Fatal("publicKey & listenPort are required")
		}
	}

	dis := createDiscovery(*discoverType)
	reg := createRegister(*registerType)

	go discoverAndUpdateRegister(dis, reg)

	// register only, passivity connect
	if *iface == "" {
		select {}
	}

	// auto update peer endpoint

	ch, err := reg.Receive()
	if err != nil {
		l.Fatal(err)
	}

	// batch update
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	peers := make(map[wg.Key]string, 8)

	for {
		var (
			update bool
		)
		select {
		case <-ticker.C:
			if len(peers) > 0 {
				update = true
			}
		case peer := <-ch:
			l.Print("peer:", peer)
			key, addr := peer.Parse()
			peers[key] = addr

			if len(peers) > 8 {
				update = true
			}
		}

		if !update {
			continue
		}

		if len(peers) == 0 {
			continue
		}

		localPeers, err := wg.GetEndpoints(*iface)
		if err != nil {
			l.Print("get endpoints failed:", err)
			time.Sleep(time.Second * 5)
			continue
		}

		for key, addr := range peers {
			origin, ok := localPeers[key]

			delete(peers, key)

			if !ok {
				continue
			} else {
				if origin == addr {
					continue
				}
			}

			err = wg.SetPeerEndpoint(*iface, key, addr)
			if err != nil {
				l.Printf("set peer (%v) endpoint failed: %v", key, err)
				continue
			}
			l.Print("resolve peer endpoint:", key, addr)
		}

		update = false
	}
}

func createDiscovery(discover string) discover.Discover {
	switch strings.ToLower(discover) {
	case "stun":
		return stun.NewDiscover()
	default:
		l.Fatalf("unknown discovery type: %s", discover)
	}
	return nil
}

func createRegister(register string) register.Register {
	switch strings.ToLower(register) {
	case "irc":
		return irc.NewRegister()
	default:
		l.Fatalf("unknown register type: %s", register)
	}
	return nil
}

func discoverAndUpdateRegister(dis discover.Discover, reg register.Register) {
	defer func() {
		if err := recover(); err != nil {
			l.Println(err)
			time.Sleep(time.Second)

			// retry
			go discoverAndUpdateRegister(dis, reg)
		}
	}()

	key := wg.NewKeyFromStr(*publicKey)
	var err error
	for {
		if key == wg.DefaultKey {
			key, err = wg.GetIfacePubKey(*iface)
			if err != nil {
				l.Println("get iface pubkey failed:", err)
				time.Sleep(time.Second * 5)
				continue
			}
		}

		lPort := uint16(*listenPort)
		if lPort == 0 {
			lPort, err = wg.GetIfaceListenPort(*iface)
			if err != nil {
				time.Sleep(time.Second * 5)
				continue
			}
		}

		ip, port, err := dis.Discover(lPort)
		if err != nil {
			l.Println("discover failed:", err)
			time.Sleep(time.Second)
			continue
		}
		err = reg.Register(key, ip, port)

		if err != nil {
			l.Println("register failed:", err)
			time.Sleep(time.Second)
			continue
		}

		time.Sleep(time.Second * 10)
	}
}
