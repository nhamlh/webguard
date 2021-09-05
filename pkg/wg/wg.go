package wg

import (
	"errors"
	"fmt"
	"github.com/nhamlh/wg-dash/pkg/config"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"net"
	"strconv"
)

type Device struct {
	c        wgctrl.Client
	dev      wgtypes.Device
	Endpoint string
	CIDR     net.IPNet
	// Routes to push to peers
	PeerRoutes []net.IPNet
	peerIps    []net.IP // a cache of allocatable IPs for peers
}

func LoadDevice(cfg config.WireguardConfig) (*Device, error) {
	client, err := wgctrl.New()
	if err != nil {
		return nil, fmt.Errorf("Cannot initialize client: %v", err)
	}

	dev, err := client.Device(cfg.Name)
	if err != nil {
		return nil, fmt.Errorf("Client cannot get device %s: %v", cfg.Name, err)
	}

	key, err := wgtypes.ParseKey(cfg.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse private key: %v", err)
	}

	ip, ipnet, err := net.ParseCIDR(cfg.Cidr)
	if err != nil {
		return nil, fmt.Errorf("Cannot parse CIDR: %v", err)
	}

	ips, _ := allIPs(ip, *ipnet)
	if len(ips) < 2 {
		return nil, errors.New("Not enough allocatable IPs for wireguard to run. It needs at least 2 IPs")
	}

	wgCfg := wgtypes.Config{
		PrivateKey: &key,
		ListenPort: &cfg.ListenPort,
	}

	err = client.ConfigureDevice(cfg.Name, wgCfg)
	if err != nil {
		return nil, err
	}

	var peerRoutes []net.IPNet
	for _, r := range cfg.PeerRoutes {
		_, ipnet, err := net.ParseCIDR(r)
		if err != nil {
			continue
		}

		peerRoutes = append(peerRoutes, *ipnet)
	}

	if len(peerRoutes) == 0 {
		return nil, errors.New("Cannot configure wireguard interface: there's no routes to push to peers")
	}

	return &Device{
		c:          *client,
		Endpoint:   cfg.Host + ":" + strconv.Itoa(cfg.ListenPort),
		dev:        *dev,
		CIDR:       *ipnet,
		PeerRoutes: peerRoutes,
		peerIps:    ips,
	}, nil
}

func (d *Device) GetPeer(pubkey wgtypes.Key) (*wgtypes.Peer, bool) {

	refreshedDev, _ := d.c.Device(d.dev.Name)
	d.dev = *refreshedDev

	for _, p := range d.dev.Peers {
		if p.PublicKey == pubkey {
			return &p, true
		}
	}
	return &wgtypes.Peer{}, false
}

func (d *Device) AddPeer(peer wgtypes.PeerConfig) bool {
	cfg := wgtypes.Config{
		Peers: []wgtypes.PeerConfig{peer},
	}

	err := d.c.ConfigureDevice(d.dev.Name, cfg)
	if err != nil {
		return false
	}

	return true
}

func (d *Device) RemovePeer(pubkey wgtypes.Key) bool {
	peer, found := d.GetPeer(pubkey)
	if !found {
		return false
	}

	peerCfg := wgtypes.PeerConfig{
		Remove:    true,
		PublicKey: peer.PublicKey,
	}

	cfg := wgtypes.Config{
		Peers: []wgtypes.PeerConfig{peerCfg},
	}

	err := d.c.ConfigureDevice(d.dev.Name, cfg)
	if err != nil {
		return false
	}

	return true
}

func (d *Device) Publickey() wgtypes.Key {
	return d.dev.PublicKey
}

func (d *Device) AllocateIP(num int) (net.IPNet, error) {
	if num < 0 || num > len(d.peerIps) {
		return net.IPNet{}, errors.New("Cannot allocate IP: Out of bound")
	}

	return net.IPNet{
		IP:   d.peerIps[num],
		Mask: net.IPv4Mask(255, 255, 255, 255),
	}, nil
}

func allIPs(ip net.IP, ipnet net.IPNet) ([]net.IP, error) {
	ips := []net.IP{}
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		// Mimic deep clone operation, otherwise ips
		// adding ip into ips directly will cause ips to contain all same ip
		tmp := net.ParseIP(ip.String())
		ips = append(ips, tmp)
	}
	// remove network address and broadcast address
	return ips[1 : len(ips)-1], nil
}

func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
