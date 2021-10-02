package main

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/nhamlh/webguard/pkg/config"
	models "github.com/nhamlh/webguard/pkg/db"
	"github.com/nhamlh/webguard/pkg/sso"
	"github.com/nhamlh/webguard/pkg/web"
	"github.com/nhamlh/webguard/pkg/wg"
	"github.com/spf13/cobra"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
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

		var peers []models.Device
		db.Select(&peers, "SELECT * FROM devices")

		for _, peer := range peers {
			peerIP, err := wgInterface.AllocateIP(peer.Num)
			if err != nil {
				log.Println(fmt.Errorf("Cannot add peer %s: Failed to allocate IP from device number: %v", peer.PrivateKey.PublicKey().String(), err))
				break
			}

			peerCfg := wgtypes.PeerConfig{
				PublicKey:         peer.PrivateKey.PublicKey(),
				AllowedIPs:        []net.IPNet{peerIP},
				ReplaceAllowedIPs: true,
			}
			if err = wgInterface.AddPeer(peerCfg); err != nil {
				log.Println(fmt.Errorf("Cannot add peer %s: %v", peerCfg.PublicKey.String(), err))
				break
			} else {
				// log.Println(fmt.Sprintf("Added peer %s", peer.PrivateKey.PublicKey()))
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
