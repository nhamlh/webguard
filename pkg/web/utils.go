package web

import (
	"sort"

	"github.com/nhamlh/webguard/pkg/db"
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

