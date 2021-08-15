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
			PrivateKey: "abc",
			Cidr:       "",
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
	ListenPort int    `json:"listen_port"`
	PrivateKey string `json:"private_key"`
	Cidr       string `json:"cidr"`
}
