package api

import (
	"encoding/json"
	"log"
	"net/http"

	db "github.com/DMV-Nicolas/DevoraTasks/db/sqlc"
	"github.com/DMV-Nicolas/DevoraTasks/token"
	"github.com/DMV-Nicolas/DevoraTasks/util"
	"github.com/gorilla/mux"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *mux.Router
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := mux.NewRouter()
	router.HandleFunc("/", Home).Methods("GET")

	router.HandleFunc("/users", server.createUser).Methods("POST")
	router.HandleFunc("/users/{id}", server.getUser).Methods("GET")

	router.HandleFunc("/tasks", server.createTask).Methods("POST")
	router.HandleFunc("/tasks", server.listTasks).Methods("GET")
	router.HandleFunc("/tasks/{id}", server.getTask).Methods("GET")
	router.HandleFunc("/tasks", server.updateTask).Methods("PUT")
	router.HandleFunc("/tasks", server.deleteTask).Methods("DELETE")

	server.router = router
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
