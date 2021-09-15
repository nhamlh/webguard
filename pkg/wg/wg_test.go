package wg

import (
	"net"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestWg(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "wg suite")
}

var _ = Describe("eqIps", func() {
	Context("Two IP arrays are equal", func() {
		It("returns true", func() {
			a := []net.IPNet{
				net.IPNet{IP: net.IP{}, Mask: []byte{255, 255, 255, 0}},
			}
			b := []net.IPNet{
				net.IPNet{IP: net.IP{}, Mask: []byte{255, 255, 255, 0}},
			}

			Expect(eqIps(a, b)).Should(BeTrue())
			Expect(eqIps(b, a)).Should(BeTrue())
		})
	})

	Context("Two IP arrays are not equal", func() {
		It("returns false", func() {
			a := []net.IPNet{
				net.IPNet{IP: net.IP{}, Mask: []byte{255, 255, 255, 0}},
			}
			b := []net.IPNet{
				net.IPNet{IP: net.IP{}, Mask: []byte{255, 255, 255, 128}},
			}
			Expect(eqIps(a, b)).Should(BeFalse())
			Expect(eqIps(b, a)).Should(BeFalse())

			b = append(b, net.IPNet{IP: net.IP{}, Mask: []byte{255, 255, 255, 0}})
			Expect(eqIps(a, b)).Should(BeFalse())
			Expect(eqIps(b, a)).Should(BeFalse())
		})
	})
})
