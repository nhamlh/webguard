package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/nhamlh/webguard/pkg/config"
	models "github.com/nhamlh/webguard/pkg/db"
	"github.com/nhamlh/webguard/pkg/sso"
	"github.com/nhamlh/webguard/pkg/web"
	"github.com/nhamlh/webguard/pkg/wg"
	"github.com/spf13/cobra"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func newStartCmd() *cobra.Command {

	cmd := &cobra.Command{
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

			wgInterface, err := wg.LoadDevice(cfg.Wireguard)
			if err != nil {
				log.Fatal(fmt.Errorf("Cannot load Wireguard interface: %v", err))
			}

			var peers []models.Device
			db.Select(&peers, "SELECT * FROM devices")

			for _, p := range peers {
				peerIP, err := wgInterface.AllocateIP(p.Num)
				if err != nil {
					log.Fatal(err)
				}

				peer := wgtypes.PeerConfig{
					PublicKey:         p.PrivateKey.PublicKey(),
					AllowedIPs:        []net.IPNet{peerIP},
					ReplaceAllowedIPs: true,
				}
				if added := wgInterface.AddPeer(peer); added {
					log.Println("Added peer", p.PrivateKey.PublicKey())
				}
			}

			log.Println("Configure SSO provider:", cfg.Web.SSO.Provider)

			pc, err := BuildProviderConfig(cfg.Web.SSO)
			if err != nil {
				log.Fatal(fmt.Errorf("Cannot configure SSO provider %s: %v", cfg.Web.SSO.Provider, err))
			}

			redirectURL := fmt.Sprintf("%s://%s:%s/login/oauth/callback",
				cfg.Web.Scheme,
				cfg.Hostname,
				strconv.Itoa(cfg.Web.ListenPort))

			op, err := sso.NewOauth2Provider(
				cfg.Web.SSO.ClientId,
				cfg.Web.SSO.ClientSecret,
				redirectURL,
				pc)

			if err != nil {
				log.Fatal(fmt.Errorf("Cannot configure SSO provider: %v", err))
			}

			router := web.NewRouter(db, wgInterface, op)

			srv := &http.Server{
				Handler: router,
				Addr:    fmt.Sprintf("%s:%d", cfg.Hostname, cfg.Web.ListenPort),
				// Good practice: enforce timeouts for servers you create!
				WriteTimeout: 15 * time.Second,
				ReadTimeout:  15 * time.Second,
			}

			log.Println(fmt.Sprintf("Web server is listening at %s:%d", cfg.Hostname, cfg.Web.ListenPort))
			err = srv.ListenAndServe()
			if err != nil {
				log.Fatal(fmt.Errorf("Web server failed: %v", err))
			}

		},
	}

	return cmd
}

func BuildProviderConfig(cfg config.SSOConfig) (sso.ProviderConfig, error) {
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
