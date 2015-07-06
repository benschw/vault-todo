package todo

import (
	"fmt"
	"testing"

	"github.com/benschw/dns-clb-go/dns"
	"github.com/benschw/opin-go/ophttp"
	"github.com/benschw/opin-go/rando"
	"github.com/benschw/opin-go/vault"
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

func GetClientAndService() (*TodoClient, *TodoService) {
	address := dns.Address{Address: "localhost", Port: uint16(rando.Port())}

	dbStr := "root:@tcp(localhost:3306)/Todo?charset=utf8&parseTime=True"

	db, err := vault.NewStaticDbProvider(dbStr)
	if err != nil {
		panic(err)
	}
	server := ophttp.NewServer(fmt.Sprintf("%s:%d", address.Address, address.Port))
	svc := &TodoService{
		Server: server,
		Db:     db,
	}

	client := &TodoClient{
		Lb:      &StubAddressGetter{Val: address},
		Address: ServiceAddress,
	}
	return client, svc
}
