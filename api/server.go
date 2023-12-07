package api

import (
	"encoding/json"
	"log"
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

	router.HandleFunc("/users", server.createUser).Methods("POST")
	router.HandleFunc("/users/{id}", server.getUser).Methods("GET")

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return http.ListenAndServe(address, server.router)
}

func errorResponse(err1 error) []byte {
	type res struct {
		Error string `json:"error"`
	}

	newRes, err2 := json.Marshal(res{
		Error: err1.Error(),
	})
	if err2 != nil {
		log.Fatal("Cannot marshal the error response")
	}

	return []byte(newRes)
}

func jsonResponse(res any) []byte {
	newRes, err := json.Marshal(res)
	if err != nil {
		log.Fatal("Cannot marshal the error response")
	}

	return newRes
}
