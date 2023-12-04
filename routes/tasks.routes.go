package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DMV-Nicolas/DevoraTasks/db"
	"github.com/DMV-Nicolas/DevoraTasks/models"
	"github.com/gorilla/mux"
)

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	json.NewDecoder(r.Body).Decode(&task)

	if createdTask := db.DB.Create(&task); createdTask.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(createdTask.Error.Error()))
		return
	}

	json.NewEncoder(w).Encode(&task)
}

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	params := mux.Vars(r)
	db.DB.First(&task, params["id"])

	if task.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Task not found"))
		return
	}

	json.NewEncoder(w).Encode(&task)
}

func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	var tasks []models.Task
	db.DB.Find(&tasks)
	json.NewEncoder(w).Encode(&tasks)
}

type deleteTaskRequest struct {
	ID uint `json:"id"`
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	var request deleteTaskRequest
	var task models.Task
	json.NewDecoder(r.Body).Decode(&request)
	db.DB.First(&task, request.ID)

	if task.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Task not found"))
		return
	}

	db.DB.Unscoped().Delete(&task)
	w.Write([]byte(fmt.Sprintf("Tarea con ID=%d, eliminada", task.ID)))
}
