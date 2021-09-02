package web

import (
	"net"
	"sort"

	"github.com/nhamlh/wg-dash/pkg/db"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
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
