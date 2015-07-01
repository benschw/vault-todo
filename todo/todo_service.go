package todo // import "github.com/benschw/vault-todo/todo"

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/benschw/opin-go/ophttp"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

func NewTodoService(server *ophttp.Server) (*TodoService, error) {
	dbStr := "root:@tcp(localhost:3306)/Todo?charset=utf8&parseTime=True"
	// Connect to Databayse
	db, err := DbOpen(dbStr)
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
	Db     gorm.DB
}

func (s *TodoService) Migrate() error {
	// Migrate
	s.Db.AutoMigrate(&Todo{})

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
	s.Server.Start()
	return nil
}

func (s *TodoService) Stop() {
	s.Server.Stop()
}

func DbOpen(dbStr string) (gorm.DB, error) {
	db, err := gorm.Open("mysql", dbStr)
	if err != nil {
		return db, err
	}
	db.SingularTable(true)
	return db, nil
}