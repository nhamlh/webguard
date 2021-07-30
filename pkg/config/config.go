package config

type Config struct {
	DataDir string
	Web WebConfig
	Wg  WgConfig
}

type WebConfig struct {
	ListenPort string
}

type WgConfig struct {
	Name       string
	ListenPort string
	Pubkey     string
}

func Load() {

}
