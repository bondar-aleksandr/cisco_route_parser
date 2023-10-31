package parser

import (
	"net/netip"
	"testing"
	"github.com/stretchr/testify/assert"
)


func Test_NH(t *testing.T) {
	t.Run("new NH", func(t *testing.T) {
		addr, _ := netip.ParseAddr("1.2.3.4")
		nh := newNextHop("1.2.3.4")
		assert.Equal(t, addr, nh.Addr)
		assert.Equal(t, true, nh.IsIP)
		assert.Equal(t, "", nh.Intf)
		nh.setIntf("Loop0")
		assert.Equal(t, "Loop0", nh.Intf)
		assert.Equal(t, false, nh.IsIP)
		assert.Equal(t, netip.Addr{}, nh.Addr)
		nh = newNextHop("Vlan10")
		assert.Equal(t, netip.Addr{}, nh.Addr)
		assert.Equal(t, false, nh.IsIP)
		assert.Equal(t, "Vlan10", nh.Intf)
	})
}