package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/nhamlh/webguard/pkg/config"
	models "github.com/nhamlh/webguard/pkg/db"
	"github.com/nhamlh/webguard/pkg/sso"
	"github.com/nhamlh/webguard/pkg/web"
	"github.com/nhamlh/webguard/pkg/wg"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Webguard server",
	Run: func(cmd *cobra.Command, args []string) {
		cfgFile, err := cmd.Flags().GetString("config")
		if err != nil {
			log.Fatal(err)
		}

		var cfg *config.Config
		if cfgFile == "" {
			cfg = &config.DefaultConfig
		} else {
			cfg = config.Load(cfgFile)
		}

		cfg.Wireguard.Host = cfg.Hostname

		db := models.InitDb(cfg.DbPath)

		wgInterface, err := wg.LoadInterface(cfg.Wireguard)
		if err != nil {
			log.Fatal(fmt.Errorf("Cannot load Wireguard interface: %v", err))
		}

		var devices []models.Device
		db.Select(&devices, "SELECT * FROM devices")

		for _, dev := range devices {
			err := dev.AddTo(wgInterface)
			if err != nil {
				log.Println(fmt.Errorf("Cannot add peer %s: %v", dev.PrivateKey.PublicKey().String(), err))
				break
			}

		}

		log.Println("Configure SSO provider:", cfg.Web.SSO.Provider)

		pc, err := buildProviderConfig(cfg.Web.SSO)
		if err != nil {
			log.Fatal(fmt.Errorf("Cannot configure SSO provider %s: %v", cfg.Web.SSO.Provider, err))
		}

		redirectURL := fmt.Sprintf("%s://%s:%d/login/oauth/callback",
			cfg.Web.Scheme,
			cfg.Hostname,
			cfg.Web.Port)

		op, err := sso.NewOauth2Provider(
			cfg.Web.SSO.ClientId,
			cfg.Web.SSO.ClientSecret,
			redirectURL,
			pc)

		if err != nil {
			log.Fatal(fmt.Errorf("Cannot configure SSO provider: %v", err))
		}

		svc := web.NewServer(*db, *wgInterface, *op)
		svc.StartAt(cfg.Web.Address, cfg.Web.Port)
	},
}

func buildProviderConfig(cfg config.SSOConfig) (sso.ProviderConfig, error) {
	var pc sso.ProviderConfig
	var err error
	switch cfg.Provider {
	case "github":
		pc = sso.GithubProvider
	case "gitlab":
		pc = sso.GitlabProvider
	case "google":
		pc = sso.GoogleProvider
	case "okta":
		pc = sso.NewOktaProvider(cfg.ProviderOpts["domain"])
	default:
		err = errors.New("Unsupported provider")
	}

	return pc, err
}
