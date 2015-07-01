package main

import (
	"flag"
	"log"
	"os"

	"github.com/benschw/opin-go/ophttp"
	"github.com/benschw/vault-todo/todo"
)

func main() {
	bind := flag.String("bind", "0.0.0.0:8080", "address to bind http server to")
	flag.Parse()

	server := ophttp.NewServer(*bind)

	svc, err := todo.NewTodoService(server)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if err = svc.Migrate(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	if err := svc.Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
