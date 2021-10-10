package db

import (
	"bytes"
	"database/sql/driver"
	"errors"
	"fmt"
	"net"
	"strings"
	"text/template"

	"encoding/base64"

	"github.com/jmoiron/sqlx"
	"github.com/nhamlh/webguard/pkg/wg"
	qrcode "github.com/skip2/go-qrcode"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func NewDevice(uid int, name string, num int, allowedIps []net.IPNet) (*Device, error) {
	prikey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return &Device{}, err
	}

	var ips []string
	for _, pr := range allowedIps {
		ips = append(ips, pr.String())
	}

	return &Device{
		UserId:     uid,
		Name:       name,
		PrivateKey: PrivateKey{prikey},
		AllowedIps: strings.Join(ips, ","),
		Num:        num,
	}, nil

}

type Device struct {
	Id         int        `db:"id"`
	UserId     int        `db:"user_id"`
	Name       string     `db:"name"`
	PrivateKey PrivateKey `db:"private_key"`
	AllowedIps string     `db:"allowed_ips"`
	Num        int        `db:"num"` // used to generate device IP
}

func (d *Device) Save(db sqlx.DB) error {
	_, err := db.Exec(`
INSERT INTO
devices(user_id, name, private_key, num, allowed_ips)
values ($1,$2,$3,$4,$5)
`,
		d.UserId,
		d.Name,
		d.PrivateKey.String(),
		d.Num,
		d.AllowedIps)

	return err
}

func (d *Device) AddTo(wgInf wg.Interface) error {

	peerCfg, err := d.peerConfig(wgInf)
	if err != nil {
		return fmt.Errorf("Cannot add device to wireguard interface: %v", err)
	}

	err = wgInf.AddPeer(peerCfg)
	if err != nil {
		return fmt.Errorf("Cannot add advice to wireguard interface: %v", err)
	}

	return nil
}

func (d *Device) IsAddedTo(wgInf wg.Interface) bool {
	if _, found := wgInf.GetPeer(d.PrivateKey.PublicKey()); found {
		return true
	}

	return false
}

func (d *Device) RemoveFrom(wgInf wg.Interface) error {
	if _, found := wgInf.GetPeer(d.PrivateKey.PublicKey()); !found {
		return errors.New("Peer not found")
	}

	peerCfg, err := d.peerConfig(wgInf)
	if err != nil {
		return fmt.Errorf("Cannot remove device from wireguard interface: %v", err)
	}

	isRemoved := wgInf.RemovePeer(peerCfg)
	if !isRemoved {
		return fmt.Errorf("Cannot remove device from wireguard interface: %v", err)
	}

	return nil
}

// GenQRCode returns base64 encoded qrcode of client config of this device
func (d *Device) GenQRCode(wgInf wg.Interface) string {
	png, err := qrcode.Encode(d.GenClientConfig(wgInf), qrcode.Medium, 256)
	if err != nil {
		return ""
	}

	return base64.StdEncoding.EncodeToString(png)
}

func (d *Device) GenClientConfig(wgInf wg.Interface) string {
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
	for _, pr := range wgInf.PeerRoutes {
		peerRoutes = append(peerRoutes, pr.String())
	}

	peerIP, _ := wgInf.AllocateIP(d.Num)
	pubkey := wgInf.Publickey()

	clientConfig := bytes.NewBufferString("")
	t.Execute(clientConfig, map[string]string{
		"PrivateKey":  d.PrivateKey.String(),
		"PeerIP":      peerIP.String(),
		"WgPublicKey": pubkey.String(),
		"WgEndpoint":  wgInf.Endpoint,
		"PeerRoutes":  strings.Join(peerRoutes, ","),
	})

	return clientConfig.String()
}

func (d *Device) peerConfig(wgInf wg.Interface) (wgtypes.PeerConfig, error) {
	peerIp, err := wgInf.AllocateIP(d.Num)
	if err != nil {
		return wgtypes.PeerConfig{}, fmt.Errorf("Cannot generate peer config: %v", err)
	}

	return wgtypes.PeerConfig{
		PublicKey:         d.PrivateKey.PublicKey(),
		AllowedIPs:        []net.IPNet{peerIp},
		ReplaceAllowedIPs: false,
	}, nil
}

func (d *Device) Status(wgInf wg.Interface) (Status, wgtypes.Peer) {
	foundPeer, found := wgInf.GetPeer(d.PrivateKey.PublicKey())
	if !found {
		return StatusNotFound, wgtypes.Peer{}
	}

	thisPeer, _ := d.peerConfig(wgInf)

	if !wg.IpsEqual(thisPeer.AllowedIPs, foundPeer.AllowedIPs) {
		return StatusConflict, wgtypes.Peer{}
	}

	return StatusOK, *foundPeer
}

type PrivateKey struct {
	wgtypes.Key
}

// Value implements Valuer interface
func (p *PrivateKey) Value() (driver.Value, error) {
	return driver.Value(p.String()), nil
}

// Scan implements Scanner interface
func (p *PrivateKey) Scan(src interface{}) error {
	key, err := wgtypes.ParseKey(src.(string))
	if err != nil {
		return err
	}

	*p = PrivateKey{key}

	return nil
}

type Status int

const (
	StatusOK       = iota // Device is added to wg
	StatusNotFound        // Device is not added to wg
	StatusConflict        // Device is added to wg but has config conflicting
)

func (s Status) String() string {
	switch s {
	case StatusOK:
		return "OK"
	case StatusNotFound:
		return "Not Found"
	default:
		return "Conflict"
	}
}
