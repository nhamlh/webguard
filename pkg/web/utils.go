package web

import (
	"net"
	"sort"
	"text/template"

	"bytes"
	"github.com/nhamlh/webguard/pkg/db"
	"github.com/nhamlh/webguard/pkg/wg"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"strings"
)

// getAvailNum returns a number to which will be used allocate
// an IP for peer/device from IPs pool of wg interface
// it should return 1 as the minimum because 0 (first IP in
// the pool) is used for the wg interface itself
// When a device is deleted from the database, its number
// can be used for future devices
func getAvailNum(devices []db.Device) int {
	sort.SliceStable(devices, func(i, j int) bool {
		return devices[i].Num < devices[j].Num
	})

	num := 0
	for _, d := range devices {
		if d.Num == num+1 {
			num += 1
		} else {
			break
		}
	}

	return num + 1
}

func generatePeerConfig(d db.Device, peerIp net.IPNet) (wgtypes.PeerConfig, error) {
	return wgtypes.PeerConfig{
		PublicKey:         d.PrivateKey.PublicKey(),
		AllowedIPs:        []net.IPNet{peerIp},
		ReplaceAllowedIPs: true,
	}, nil
}

func generateClientConfig(wgInt *wg.Interface, d db.Device) string {
	t, _ := template.New("clientConfig").Parse(`
[Interface]
PrivateKey = {{ .PrivateKey }}
Address = {{ .PeerIP }}

[Peer]
PublicKey = {{ .WgPublicKey }}
Endpoint = {{ .WgEndpoint }}
AllowedIPs = {{ .PeerRoutes }}
`)

	var peerRoutes []string
	for _, pr := range wgInt.PeerRoutes {
		peerRoutes = append(peerRoutes, pr.String())
	}

	peerIP, _ := wgInt.AllocateIP(d.Num)
	pubkey := wgInt.Publickey()

	clientConfig := bytes.NewBufferString("")
	t.Execute(clientConfig, map[string]string{
		"PrivateKey":  d.PrivateKey.String(),
		"PeerIP":      peerIP.String(),
		"WgPublicKey": pubkey.String(),
		"WgEndpoint":  wgInt.Endpoint,
		"PeerRoutes":  strings.Join(peerRoutes, ","),
	})

	return clientConfig.String()
}
