package main

import (
	"net/http"

	"github.com/DMV-Nicolas/DevoraTasks/db"
	"github.com/DMV-Nicolas/DevoraTasks/models"
	"github.com/DMV-Nicolas/DevoraTasks/routes"
	"github.com/gorilla/mux"
)

func main() {
	db.DBConnection()
	db.DB.AutoMigrate(models.User{})
	db.DB.AutoMigrate(models.Task{})

	router := mux.NewRouter()

	router.HandleFunc("/", routes.HomeHandler)

	router.HandleFunc("/users", routes.CreateUserHandler).Methods("POST")
	router.HandleFunc("/users/{id}", routes.GetUserHandler).Methods("GET")
	router.HandleFunc("/users", routes.GetUsersHandler).Methods("GET")
	router.HandleFunc("/users", routes.DeleteUserHandler).Methods("DELETE")

	router.HandleFunc("/tasks", routes.CreateTaskHandler).Methods("POST")
	router.HandleFunc("/tasks/{id}", routes.GetTaskHandler).Methods("GET")
	router.HandleFunc("/tasks", routes.GetTasksHandler).Methods("GET")
	router.HandleFunc("/tasks", routes.DeleteTaskHandler).Methods("DELETE")

	http.ListenAndServe(":5000", router)
}
