package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	db "github.com/DMV-Nicolas/DevoraTasks/db/sqlc"
	"github.com/DMV-Nicolas/DevoraTasks/util"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

type createTaskRequest struct {
	UserID      int64  `json:"user_id" requirements:"required"`
	Title       string `json:"title" requirements:"required"`
	Description string `json:"description"`
}

func (server *Server) createTask(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest
	json.NewDecoder(r.Body).Decode(&req)

	err := util.VerifyRequirements(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponse(err))
		return
	}

	arg := db.CreateTaskParams{
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
	}

	task, err := server.db.CreateTask(context.Background(), arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation":
				w.WriteHeader(http.StatusBadRequest)
				w.Write(errorResponse(err))
				return
			}
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorResponse(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse(task))
}

func (server *Server) listTasks(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()

	offset, err := strconv.Atoi(v.Get("offset"))
	limit, err := strconv.Atoi(v.Get("limit"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponse(err))
		return
	}

	arg := db.ListTasksParams{
		Offset: int32(offset),
		Limit:  int32(limit),
	}

	tasks, err := server.db.ListTasks(context.Background(), arg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse(tasks))
}

func (server *Server) getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponse(err))
		return
	}

	task, err := server.db.GetTask(context.Background(), int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errorResponse(err))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse(task))
}

func (server *Server) updateTask(w http.ResponseWriter, r *http.Request) {
}
func (server *Server) deleteTask(w http.ResponseWriter, r *http.Request) {
}
