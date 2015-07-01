package todo

import (
	"fmt"
	"testing"

	"github.com/benschw/dns-clb-go/dns"
	. "gopkg.in/check.v1"
)

type StubAddressGetter struct {
	Val dns.Address
}

func (lb *StubAddressGetter) GetAddress(address string) (dns.Address, error) {
	if address == ServiceAddress {
		return lb.Val, nil
	}
	return dns.Address{}, fmt.Errorf("invalid service name")
}
func Test(t *testing.T) { TestingT(t) }
