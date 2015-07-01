package main

import (
	"net/http"
	"os"

	"github.com/benschw/opin-go/ophttp"
	"github.com/benschw/opin-go/rando"
	"github.com/benschw/opin-go/rest"
)

type Resource struct {
	Client *GreetingClient
}

func (r *Resource) StatusHandler(resp http.ResponseWriter, req *http.Request) {
	rest.SetOKResponse(resp, HealthStatus{Status: "UP"})
}

func (r *Resource) DemoHandler(resp http.ResponseWriter, req *http.Request) {
	host, _ := os.Hostname()

	greeting, _ := r.Client.GetGreeting()

	rest.SetOKResponse(resp, &DemoGreeting{
		Message:  "hello from myapp on " + host + "/" + rando.MyIp(),
		Greeting: greeting,
	})
}

func (r *Resource) GreetingHandler(resp http.ResponseWriter, req *http.Request) {
	host, _ := os.Hostname()

	rest.SetOKResponse(resp, Greeting{Message: "hello from mysvc on " + host + "/" + rando.MyIp()})
}

// Wire and start http server
func RunServer(server *ophttp.Server) {

	r := &Resource{Client: NewGreetingClient()}

	http.Handle("/status", http.HandlerFunc(r.StatusHandler))

	http.Handle("/greeting", http.HandlerFunc(r.GreetingHandler))
	http.Handle("/demo", http.HandlerFunc(r.DemoHandler))
	server.Start()
}
