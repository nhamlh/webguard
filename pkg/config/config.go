package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// ClientId:     "client id",
// ClientSecret: "client secret",

// github ClientId:     "e52c25ec8117d00eacfa",
// github ClientSecret: "6d79ea3b83da22578ce52fc234954234916e7fea",
// gitlab ClientId:     "6a9c2ebe3cf9ff880b1fa85ed4fbdd37322116c8f528940ab3641660c7ce063c",
// gitlab ClientSecret: "dd46ba1c366986726795e41bd0e4663d6503369c24f7df137c289c5b4f70ae7a",
// Provider:     "okta",
// ClientId:     "0oa1pdcc3haxpJuGr5d7",
// ClientSecret: "_8uAc-mhWf1xIopyy15m85-_osDo8_i2eFUIG7le",
// ProviderOpts: map[string]string{"domain": "dev-30186313.okta.com"},

var (
	DefaultConfig = Config{
		DbPath:   "./webguard.db",
		Hostname: "localhost",
		Web: WebConfig{
			Scheme:     "http",
			ListenPort: 8080,
			SSO: SSOConfig{
				Provider: "gitlab",
				// gitlab
				ClientId:     "6a9c2ebe3cf9ff880b1fa85ed4fbdd37322116c8f528940ab3641660c7ce063c",
				ClientSecret: "dd46ba1c366986726795e41bd0e4663d6503369c24f7df137c289c5b4f70ae7a",
				// github
				// ClientId:     "e52c25ec8117d00eacfa",
				// ClientSecret: "6d79ea3b83da22578ce52fc234954234916e7fea",
				// google
				// ClientId:     "95881753105-41nshnkv2b5mi0s4gbi2tvios135fk29.apps.googleusercontent.com",
				// ClientSecret: "e_3Ylg1jRhI0ooq6-xErj__q",
			},
		},
		Wireguard: WireguardConfig{
			Name:       "webguard",
			ListenPort: 51820,
			PrivateKey: "ANbdTCP22uZP3AzTdan2v6qXGRcdZRngkno0PnCPlkg=",
			Cidr:       "192.168.0.0/24",
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
	Scheme     string    `json:"scheme"`
	ListenPort int       `json:"listen_port"`
	SSO        SSOConfig `json:"sso"`
}

type SSOConfig struct {
	Provider     string            `json:"provider"`
	ClientId     string            `json:"client_id"`
	ClientSecret string            `json:"client_secret"`
	ProviderOpts map[string]string `json:"provider_options"`
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
