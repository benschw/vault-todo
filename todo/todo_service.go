package todo

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/benschw/opin-go/ophttp"
	"github.com/benschw/opin-go/vault"
	"github.com/gorilla/mux"
)

func NewTodoService(bind string) (*TodoService, error) {
	server := ophttp.NewServer(bind)

	//dbStr := "root:@tcp(localhost:3306)/Todo?charset=utf8&parseTime=True"

	db, err := vault.NewDbProvider("Todo", "mysql.service.consul")
	if err != nil {
		return nil, err
	}
	return &TodoService{
		Server: server,
		Db:     db,
	}, nil
}

type TodoService struct {
	Server *ophttp.Server
	Db     vault.DbProvider
}

func (s *TodoService) Migrate() error {
	// Migrate
	db, err := s.Db.Get()
	if err != nil {
		return err
	}
	db.AutoMigrate(&Todo{})

	return nil
}

// Configure and start http server
func (s *TodoService) Run() error {
	// Build Resource
	resource := &TodoResource{Db: s.Db}

	// Wire Routes
	r := mux.NewRouter()
	r.HandleFunc("/todo", resource.Add).Methods("POST")
	r.HandleFunc("/todo", resource.GetAll).Methods("GET")
	r.HandleFunc("/todo/{id}", resource.Get).Methods("GET")

	http.Handle("/", r)

	// Start Server
	err := s.Server.Start()
	log.Println("Server Stopped")
	return err
}

func (s *TodoService) Stop() {
	log.Println("Stopping Server...")
	s.Server.Stop()
}
