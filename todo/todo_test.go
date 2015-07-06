package todo

import (
	"fmt"
	"net/http"

	"github.com/benschw/opin-go/rest"

	. "gopkg.in/check.v1"
)

var _ = Suite(&TestSuite{})

type TestSuite struct {
	svc    *TodoService
	client *TodoClient
}

func (s *TestSuite) SetUpSuite(c *C) {
	s.client, s.svc = GetClientAndService()

	go s.svc.Run()
}
func (s *TestSuite) TearDownSuite(c *C) {
	s.svc.Stop()
}
func (s *TestSuite) SetUpTest(c *C) {
	s.svc.Migrate()
}
func (s *TestSuite) TearDownTest(c *C) {
	db, err := s.svc.Db.Get()
	if err != nil {
		panic(err)
	}
	db.DropTable(Todo{})
}

func (s *TestSuite) TestHealth(c *C) {
	// when
	r, err := rest.MakeRequest("GET", fmt.Sprintf("http://%s/health", s.svc.Server.Bind), nil)

	// then
	c.Assert(err, Equals, nil)
	c.Assert(r.StatusCode, Equals, http.StatusOK)
}
func (s *TestSuite) TestAdd(c *C) {
	// given
	expected := &Todo{Id: 1, Content: "hello world", Status: "new"}

	// when
	created, err := s.client.Add("hello world")

	// then
	c.Assert(err, Equals, nil)

	found, _ := s.client.Find(created.Id)

	c.Assert(created, DeepEquals, expected)
	c.Assert(found, DeepEquals, expected)
}
func (s *TestSuite) TestGet(c *C) {
	// given
	expected := &Todo{Id: 1, Content: "hello world", Status: "new"}

	s.client.Add("hello world")

	// when
	found, err := s.client.Find(1)

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

	s.client.Add("hello world")
	s.client.Add("hello universe")
	s.client.Add("hello galaxy")

	// when
	found, err := s.client.FindAll()

	// then
	c.Assert(err, Equals, nil)

	c.Assert(found, DeepEquals, expected)
}
func (s *TestSuite) TestSave(c *C) {
	// given
	expected := &Todo{Id: 1, Content: "foo", Status: "in-progress"}

	todo, _ := s.client.Add("hello world")

	// when
	todo.Content = "foo"
	todo.Status = "in-progress"
	saved, err := s.client.Save(todo)

	// then
	c.Assert(err, Equals, nil)

	found, _ := s.client.Find(todo.Id)

	c.Assert(saved, DeepEquals, expected)
	c.Assert(found, DeepEquals, expected)
}
func (s *TestSuite) TestDelete(c *C) {
	// given
	todo, _ := s.client.Add("hello world")

	// when
	err := s.client.Delete(todo.Id)

	// then
	c.Assert(err, Equals, nil)

	_, err = s.client.Find(todo.Id)

	c.Assert(err, DeepEquals, fmt.Errorf("404: Not Found"))
}
