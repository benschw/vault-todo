package todo

import (
	"fmt"

	"github.com/benschw/dns-clb-go/dns"
	"github.com/benschw/opin-go/ophttp"
	"github.com/benschw/opin-go/rando"
	"github.com/benschw/opin-go/vault"
	. "gopkg.in/check.v1"
)

type TestSuite struct {
	s       *TodoService
	server  *ophttp.Server
	address dns.Address
}

var _ = Suite(&TestSuite{})

func (s *TestSuite) SetUpSuite(c *C) {
	s.address = dns.Address{Address: "localhost", Port: uint16(rando.Port())}

	dbStr := "root:@tcp(localhost:3306)/Todo?charset=utf8&parseTime=True"

	db, err := vault.NewStaticDbProvider(dbStr)
	if err != nil {
		panic(err)
	}
	s.server = ophttp.NewServer(fmt.Sprintf("%s:%d", s.address.Address, s.address.Port))
	s.s = &TodoService{
		Server: s.server,
		Db:     db,
	}
	go s.s.Run()
}
func (s *TestSuite) TearDownSuite(c *C) {
	s.server.Stop()
}
func (s *TestSuite) SetUpTest(c *C) {
	s.s.Migrate()
}
func (s *TestSuite) TearDownTest(c *C) {
	db, err := s.s.Db.Get()
	if err != nil {
		panic(err)
	}
	db.DropTable(Todo{})
}

func (s *TestSuite) TestAdd(c *C) {
	// given
	expected := &Todo{Id: 1, Content: "hello world", Status: "new"}

	client := TodoClient{
		Lb:      &StubAddressGetter{Val: s.address},
		Address: ServiceAddress,
	}

	// when
	created, err := client.Add("hello world")

	// then
	c.Assert(err, Equals, nil)

	found, _ := client.Find(created.Id)

	c.Assert(created, DeepEquals, expected)
	c.Assert(found, DeepEquals, expected)
}
func (s *TestSuite) TestGet(c *C) {
	// given
	expected := &Todo{Id: 1, Content: "hello world", Status: "new"}

	client := TodoClient{
		Lb:      &StubAddressGetter{Val: s.address},
		Address: ServiceAddress,
	}

	client.Add("hello world")

	// when
	found, err := client.Find(1)

	// then
	c.Assert(err, Equals, nil)

	c.Assert(found, DeepEquals, expected)
}
func (s *TestSuite) TestFindAll(c *C) {
	// given
	expected := []*Todo{
		&Todo{Id: 1, Content: "hello world", Status: "new"},
		&Todo{Id: 2, Content: "hello universe", Status: "new"},
		&Todo{Id: 3, Content: "hello galaxy", Status: "new"},
	}
	client := TodoClient{
		Lb:      &StubAddressGetter{Val: s.address},
		Address: ServiceAddress,
	}

	client.Add("hello world")
	client.Add("hello universe")
	client.Add("hello galaxy")

	// when
	found, err := client.FindAll()

	// then
	c.Assert(err, Equals, nil)

	c.Assert(found, DeepEquals, expected)
}
func (s *TestSuite) TestSave(c *C) {
	// given
	expected := &Todo{Id: 1, Content: "foo", Status: "in-progress"}

	client := TodoClient{
		Lb:      &StubAddressGetter{Val: s.address},
		Address: ServiceAddress,
	}

	todo, _ := client.Add("hello world")

	// when
	todo.Content = "foo"
	todo.Status = "in-progress"
	saved, err := client.Save(todo)

	// then
	c.Assert(err, Equals, nil)

	found, _ := client.Find(todo.Id)

	c.Assert(saved, DeepEquals, expected)
	c.Assert(found, DeepEquals, expected)
}
func (s *TestSuite) TestDelete(c *C) {
	// given
	client := TodoClient{
		Lb:      &StubAddressGetter{Val: s.address},
		Address: ServiceAddress,
	}

	todo, _ := client.Add("hello world")

	// when
	err := client.Delete(todo.Id)

	// then
	c.Assert(err, Equals, nil)

	_, err = client.Find(todo.Id)

	c.Assert(err, DeepEquals, fmt.Errorf("404: Not Found"))
}
