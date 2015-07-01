package todo

import (
	"fmt"
	"log"
	"net/http"

	"github.com/benschw/opin-go/rest"
	"github.com/benschw/opin-go/vaultdb"
)

type TodoResource struct {
	Db vaultdb.DbProvider
}

func (r *TodoResource) Add(res http.ResponseWriter, req *http.Request) {
	var todo Todo

	db, err := r.Db.Get()
	if err != nil {
		rest.SetInternalServerErrorResponse(res, err)
	}

	if err := rest.Bind(req, &todo); err != nil {
		log.Print(err)
		rest.SetBadRequestResponse(res)
		return
	}

	db.Create(&todo)

	if err := rest.SetCreatedResponse(res, todo, fmt.Sprintf("todo/%d", todo.Id)); err != nil {
		rest.SetInternalServerErrorResponse(res, err)
		return
	}
}

func (r *TodoResource) Get(res http.ResponseWriter, req *http.Request) {
	id, err := rest.PathInt(req, "id")
	if err != nil {
		rest.SetBadRequestResponse(res)
		return
	}
	var todo Todo

	db, err := r.Db.Get()
	if err != nil {
		rest.SetInternalServerErrorResponse(res, err)
	}
	if db.First(&todo, id).RecordNotFound() {
		rest.SetNotFoundResponse(res)
		return
	}

	if err := rest.SetOKResponse(res, todo); err != nil {
		rest.SetInternalServerErrorResponse(res, err)
		return
	}
}

func (r *TodoResource) GetAll(resp http.ResponseWriter, req *http.Request) {
}
func (r *TodoResource) Update(resp http.ResponseWriter, req *http.Request) {
}
