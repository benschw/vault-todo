package main

import (
	"flag"
	"log"
	"os"

	"github.com/benschw/vault-todo/todo"
)

func main() {
	bind := flag.String("bind", "0.0.0.0:8080", "address to bind http server to")
	flag.Parse()

	log.Print("constructing service")
	svc, err := todo.NewTodoService(*bind)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Print("migrating")

	if err = svc.Migrate(); err != nil {
		log.Println(err)
		os.Exit(1)
	}

	log.Print("running service")
	if err := svc.Run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
