package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

var (
	DefaultConfig = Config{
		DbPath:   "/tmp/wg-dash.db",
		Hostname: "localhost",
		Web: WebConfig{
			ListenPort: 8080,
		},
		Wireguard: WireguardConfig{
			Name:       "wg-dash",
			ListenPort: 51820,
			PrivateKey: "ANbdTCP22uZP3AzTdan2v6qXGRcdZRngkno0PnCPlkg=",
			Cidr:       "192.168.0.0/29",
			PeerRoutes: []string{"0.0.0.0/0"},
		},
	}
)

type Config struct {
	DbPath string `json:"dbpath"`
	// Hostname will be used for both web and wireguard endpoint
	Hostname  string          `json:"hostname"`
	Web       WebConfig       `json:"web"`
	Wireguard WireguardConfig `json:"wireguard"`
}

func Load(configFile string) *Config {
	bytes, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}

	cfg := DefaultConfig
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &cfg
}

type WebConfig struct {
	ListenPort int `json:"listen_port"`
}

type WireguardConfig struct {
	Name       string `json:"name"`
	Host       string
	ListenPort int    `json:"listen_port"`
	PrivateKey string `json:"private_key"`
	// CIDR to allocate IPs for peers
	// CIDR has form of a.b.c.d/x
	Cidr string `json:"cidr"`
	// PeerRoutes will be pushed to peer's PeerRoutes directive in config file
	PeerRoutes []string `json:"peer_routes"`
}
