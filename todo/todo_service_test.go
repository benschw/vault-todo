package todo

import (
	"fmt"

	"github.com/benschw/dns-clb-go/dns"
	"github.com/benschw/opin-go/ophttp"
	"github.com/benschw/opin-go/rando"
	. "gopkg.in/check.v1"
)

type TestSuite struct {
	s       *TodoService
	server  *ophttp.Server
	address dns.Address
}

var _ = Suite(&TestSuite{})

func (s *TestSuite) SetUpTest(c *C) {
	s.address = dns.Address{Address: "localhost", Port: uint16(rando.Port())}

	s.server = ophttp.NewServer(fmt.Sprintf("%s:%d", s.address.Address, s.address.Port))
	svc, _ := NewTodoService(s.server)
	s.s = svc
}

func (s *TestSuite) TearDownTest(c *C) {
	s.s.Db.DropTable(Todo{})
	s.server.Stop()
}

func (s *TestSuite) TestTodo(c *C) {
	// given

	s.s.Migrate()
	go s.s.Run()

	expected := &Todo{Content: "hello world", Status: "new"}

	client := TodoClient{
		Lb:      &StubAddressGetter{Val: s.address},
		Address: ServiceAddress,
	}

	// when
	todo, err := client.Add("hello world")
	found, err2 := client.Find(todo.Id)

	// then
	c.Assert(err, Equals, nil)
	c.Assert(err2, Equals, nil)

	c.Assert(expected.Content, Equals, todo.Content)
	c.Assert(todo, Equals, found)
}
