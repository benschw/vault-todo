package todo

import (
	"fmt"
	"net/http"

	"github.com/benschw/opin-go/rest"

	"github.com/benschw/dns-clb-go/clb"
	"github.com/benschw/dns-clb-go/dns"
)

const ServiceAddress = "todo.service.consul"

// Interface for Load Balancer
type AddressGetter interface {
	GetAddress(string) (dns.Address, error)
}

// Client Factory
func NewTodoClient() *TodoClient {
	return &TodoClient{
		Lb:      clb.NewDefaultClb(clb.RoundRobin),
		Address: ServiceAddress,
	}
}

// Client
type TodoClient struct {
	Lb      AddressGetter
	Address string
}

func (c *TodoClient) Add(content string) (*Todo, error) {
	var created *Todo
	host, err := c.Lb.GetAddress(c.Address)
	if err != nil {
		return created, err
	}
	todo := &Todo{Content: content, Status: "new"}

	r, err := rest.MakeRequest("POST", fmt.Sprintf("http://%s/todo", host), todo)
	if err != nil {
		return created, err
	}
	err = rest.ProcessResponseEntity(r, &created, http.StatusCreated)
	return created, err
}

func (c *TodoClient) Find(id int) (*Todo, error) {
	var found *Todo
	host, err := c.Lb.GetAddress(c.Address)
	if err != nil {
		return found, err
	}

	r, err := rest.MakeRequest("GET", fmt.Sprintf("http://%s/todo/%d", host, id), nil)
	if err != nil {
		return found, err
	}
	err = rest.ProcessResponseEntity(r, &found, http.StatusOK)
	return found, err
}

func (c *TodoClient) FindAll() ([]*Todo, error) {
	var todos []*Todo

	host, err := c.Lb.GetAddress(c.Address)
	if err != nil {
		return todos, err
	}

	r, err := rest.MakeRequest("GET", fmt.Sprintf("http://%s/todo", host), nil)
	if err != nil {
		return todos, err
	}
	err = rest.ProcessResponseEntity(r, &todos, http.StatusOK)
	return todos, err
}
func (c *TodoClient) Save(todo *Todo) (*Todo, error) {
	var saved *Todo
	host, err := c.Lb.GetAddress(c.Address)
	if err != nil {
		return saved, err
	}

	r, err := rest.MakeRequest("PUT", fmt.Sprintf("http://%s/todo/%d", host, todo.Id), todo)
	if err != nil {
		return saved, err
	}
	err = rest.ProcessResponseEntity(r, &saved, http.StatusOK)
	return saved, err
}
func (c *TodoClient) Delete(id int) error {
	host, err := c.Lb.GetAddress(c.Address)
	if err != nil {
		return err
	}

	r, err := rest.MakeRequest("DELETE", fmt.Sprintf("http://%s/todo/%d", host, id), nil)
	if err != nil {
		return err
	}
	return rest.ProcessResponseEntity(r, nil, http.StatusNoContent)
}
