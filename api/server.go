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
	router.Use(corsMiddleware)
	router.HandleFunc("/", Home).Methods("GET")

	router.HandleFunc("/users", server.createUser).Methods("POST")
	router.HandleFunc("/users/login", server.loginUser).Methods("POST")
	router.HandleFunc("/users", authMiddleware(server.getUser, server.tokenMaker)).Methods("GET")

	router.HandleFunc("/tasks", authMiddleware(server.createTask, server.tokenMaker)).Methods("POST")
	router.HandleFunc("/tasks", authMiddleware(server.listTasks, server.tokenMaker)).Methods("GET")
	router.HandleFunc("/tasks/{id}", authMiddleware(server.getTask, server.tokenMaker)).Methods("GET")
	router.HandleFunc("/tasks", authMiddleware(server.updateTask, server.tokenMaker)).Methods("PUT")
	router.HandleFunc("/tasks", authMiddleware(server.deleteTask, server.tokenMaker)).Methods("DELETE")

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
