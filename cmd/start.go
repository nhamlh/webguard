package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nhamlh/wg-dash/pkg/config"
	"github.com/nhamlh/wg-dash/pkg/db"
	"github.com/nhamlh/wg-dash/pkg/web"
	"github.com/nhamlh/wg-dash/pkg/wg"
	"github.com/spf13/cobra"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"net"
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

			wgInterface := wg.LoadDevice(cfg.Wireguard)

			var peers []db.Device
			db.DB.Select(&peers, "SELECT * FROM devices")

			for _, p := range peers {
				prikey, err := wgtypes.ParseKey(p.PrivateKey)
				if err != nil {
					log.Fatal(err)
				}

				peerIP, err := wgInterface.AllocateIP(p.Num)
				if err != nil {
					log.Fatal(err)
				}

				peer := wgtypes.PeerConfig{
					PublicKey:         prikey.PublicKey(),
					AllowedIPs:        []net.IPNet{peerIP},
					ReplaceAllowedIPs: true,
				}
				if added := wgInterface.AddPeer(peer); added {
					fmt.Println("Added peer")
				}
			}

			router := web.NewRouterFor(wgInterface)

			srv := &http.Server{
				Handler: router,
				Addr:    fmt.Sprintf("%s:%d", cfg.Hostname, cfg.Web.ListenPort),
				// Good practice: enforce timeouts for servers you create!
				WriteTimeout: 15 * time.Second,
				ReadTimeout:  15 * time.Second,
			}

			fmt.Println(fmt.Sprintf("Web server is listening at %s:%d", cfg.Hostname, cfg.Web.ListenPort))
			srv.ListenAndServe()
		},
	}

	return cmd
}
