package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/benschw/dns-clb-go/dns"
	"github.com/benschw/opin-go/ophttp"
	"github.com/benschw/opin-go/rando"
)

type GreetingAddressGetter struct {
	Val dns.Address
}

func (lb *GreetingAddressGetter) GetAddress(address string) (dns.Address, error) {
	if address == ServiceAddress {
		return lb.Val, nil
	}
	return dns.Address{}, fmt.Errorf("invalid service name")
}

// Component test for greeting endpoint:
// exercise running server with the client
func TestGreetingEndpoint(t *testing.T) {
	// given
	host, _ := os.Hostname()

	expectedGreeting := &Greeting{Message: "hello from greeting on " + host + "/" + rando.MyIp()}

	address := dns.Address{Address: "localhost", Port: uint16(rando.Port())}

	server := ophttp.NewServer(fmt.Sprintf("%s:%d", address.Address, address.Port))
	go RunServer(server)

	client := GreetingClient{
		Lb:      &GreetingAddressGetter{Val: address},
		Address: ServiceAddress,
	}

	// when
	greeting, _ := client.GetGreeting()

	// then
	if expectedGreeting.Message != greeting.Message {
		t.Errorf("expected '%s', got '%s'", expectedGreeting.Message, greeting.Message)
	}

	// teardown
	server.Stop()
}
