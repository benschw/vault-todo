package main

import (
	"flag"

	"github.com/benschw/opin-go/ophttp"
)

func main() {
	bind := flag.String("bind", "0.0.0.0:8080", "address to bind http server to")
	flag.Parse()

	server := ophttp.NewServer(*bind)

	RunServer(server)
}
