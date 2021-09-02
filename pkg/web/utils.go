package web

import (
	"html/template"
	"net"
	"sort"

	"bytes"
	"github.com/nhamlh/wg-dash/pkg/db"
	"github.com/nhamlh/wg-dash/pkg/wg"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"strings"
)

func getAvailNum(devices []db.Device) int {
	sort.SliceStable(devices, func(i, j int) bool {
		return devices[i].Num < devices[j].Num
	})

	num := -1
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
	prikey, err := wgtypes.ParseKey(d.PrivateKey)
	if err != nil {
		return wgtypes.PeerConfig{}, err
	}

	return wgtypes.PeerConfig{
		PublicKey:         prikey.PublicKey(),
		AllowedIPs:        []net.IPNet{peerIp},
		ReplaceAllowedIPs: true,
	}, nil
}

func generateClientConfig(wgInt *wg.Device, d db.Device) string {
	t, _ := template.New("clientConfig").Parse(`
[Interface]
PrivateKey = {{ .PrivateKey }}
Address = {{ .PeerIP }}

[Peer]
PublicKey = {{ .WgPublicKey }}
Endpoint = {{ .WgEndpoint }}
AllowedIPs = {{ .PeerRoutes }}
`)

	prikey, _ := wgtypes.ParseKey(d.PrivateKey)

	var peerRoutes []string
	for _, pr := range wgInt.PeerRoutes {
		peerRoutes = append(peerRoutes, pr.String())
	}

	peerIP, _ := wgInt.AllocateIP(d.Num)
	pubkey := wgInt.Publickey()

	clientConfig := bytes.NewBufferString("")
	t.Execute(clientConfig, map[string]string{
		"PrivateKey":  prikey.String(),
		"PeerIP":      peerIP.String(),
		"WgPublicKey": pubkey.String(),
		"WgEndpoint":  wgInt.Endpoint,
		"PeerRoutes":  strings.Join(peerRoutes, ","),
	})

	return clientConfig.String()
}
