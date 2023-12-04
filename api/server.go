package api

import (
	"encoding/json"
	"net/http"

	db "github.com/DMV-Nicolas/DevoraTasks/db/sqlc"
	"github.com/gorilla/mux"
)

type Server struct {
	db     *db.Queries
	router *mux.Router
}

func NewServer(db *db.Queries) *Server {
	server := &Server{db: db}
	router := mux.NewRouter()

	router.HandleFunc("/", Home).Methods("GET")

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return http.ListenAndServe(address, server.router)
}

func errorResponse(err error) []byte {
	type res struct {
		Error string `json:"error"`
	}

	newRes, _ := json.Marshal(res{
		Error: err.Error(),
	})

	return []byte(newRes)

}
