package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/DMV-Nicolas/DevoraTasks/db"
	"github.com/DMV-Nicolas/DevoraTasks/models"
	"github.com/gorilla/mux"
)

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	if createdUser := db.DB.Create(&user); createdUser.Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(createdUser.Error.Error()))
		return
	}

	json.NewEncoder(w).Encode(&user)
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	params := mux.Vars(r)
	db.DB.First(&user, params["id"])

	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}
	db.DB.Model(&user).Association("Tasks").Find(&user.Tasks)

	json.NewEncoder(w).Encode(&user)
}

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	db.DB.Find(&users)
	json.NewEncoder(w).Encode(&users)
}

type deleteUserRequest struct {
	ID uint `json:"id"`
}

func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	var request deleteUserRequest
	var user models.User
	json.NewDecoder(r.Body).Decode(&request)
	db.DB.First(&user, request.ID)

	if user.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("User not found"))
		return
	}

	db.DB.Unscoped().Delete(&user)
	w.Write([]byte(fmt.Sprintf("Usuario con ID=%d, eliminado", user.ID)))
}
