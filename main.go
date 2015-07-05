package main

import (
	"flag"
	"log"
	"log/syslog"
	"os"

	"github.com/benschw/vault-todo/todo"
)

func main() {

	bind := flag.String("bind", "0.0.0.0:8080", "address to bind http server to")
	useSyslog := flag.Bool("syslog", false, "log to syslog")
	flag.Parse()

	if *useSyslog {
		logwriter, err := syslog.New(syslog.LOG_NOTICE, "todo")
		if err == nil {
			log.SetOutput(logwriter)
		}
	}

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
