package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/nhamlh/wg-dash/pkg/config"
	"github.com/nhamlh/wg-dash/pkg/db"
	"github.com/nhamlh/wg-dash/pkg/sso"
	"github.com/nhamlh/wg-dash/pkg/web"
	"github.com/nhamlh/wg-dash/pkg/wg"
	"github.com/spf13/cobra"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func newStartCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start wg-dash server",
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

			wgInterface := wg.LoadDevice(cfg.Wireguard)

			var peers []db.Device
			db.DB.Select(&peers, "SELECT * FROM devices")

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

			redirectURL := fmt.Sprintf("%s://%s:%s/login/oauth/callback",
				cfg.Web.Scheme,
				cfg.Hostname,
				strconv.Itoa(cfg.Web.ListenPort))

			op, err := sso.NewOauth2Provider(
				cfg.Web.SSO.ClientId,
				cfg.Web.SSO.ClientSecret,
				redirectURL,
				cfg.Web.SSO.Provider)

			if err != nil {
				log.Fatal(err)
			}

			router := web.NewRouter(wgInterface, op)

			srv := &http.Server{
				Handler: router,
				Addr:    fmt.Sprintf("%s:%d", cfg.Hostname, cfg.Web.ListenPort),
				// Good practice: enforce timeouts for servers you create!
				WriteTimeout: 15 * time.Second,
				ReadTimeout:  15 * time.Second,
			}

			log.Println(fmt.Sprintf("Web server is listening at %s:%d", cfg.Hostname, cfg.Web.ListenPort))
			srv.ListenAndServe()
		},
	}

	return cmd
}
