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
)

type createTaskRequest struct {
	UserID      int64  `json:"user_id" requirements:"min=1"`
	Title       string `json:"title" requirements:"required"`
	Description string `json:"description"`
}

func (server *Server) createTask(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponse(err))
		return
	}

	err = util.VerifyRequirements(req)
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
		if db.ErrorCode(err) == db.ForeignKeyViolation {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errorResponse(err))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorResponse(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse(task))
}

type listTasksRequest struct {
	Offset int32 `requirements:"min=0"`
	Limit  int32 `requirements:"min=1"`
}

func (server *Server) listTasks(w http.ResponseWriter, r *http.Request) {
	var req listTasksRequest

	v := r.URL.Query()
	offset, _ := strconv.Atoi(v.Get("offset"))
	limit, _ := strconv.Atoi(v.Get("limit"))

	req.Offset = int32(offset)
	req.Limit = int32(limit)
	err := util.VerifyRequirements(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponse(err))
		return
	}

	arg := db.ListTasksParams{
		Offset: req.Offset,
		Limit:  req.Limit,
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

type getTaskRequest struct {
	ID int64 `requirements:"min=1"`
}

func (server *Server) getTask(w http.ResponseWriter, r *http.Request) {
	var req getTaskRequest
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	req.ID = int64(id)
	err := util.VerifyRequirements(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponse(err))
		return
	}

	task, err := server.db.GetTask(context.Background(), int64(id))
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
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

type updateTaskRequest struct {
	ID          int64  `json:"id" requirements:"min=1"`
	Title       string `json:"title" requirements:"required"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

func (server *Server) updateTask(w http.ResponseWriter, r *http.Request) {
	var req updateTaskRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponse(err))
		return
	}

	err = util.VerifyRequirements(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponse(err))
		return
	}

	arg := db.UpdateTaskParams{
		ID:          req.ID,
		Title:       req.Title,
		Description: req.Description,
		Done:        req.Done,
	}

	task, err := server.db.UpdateTask(context.Background(), arg)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
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

type deleteTaskRequest struct {
	ID int64 `json:"id" requirements:"min=1"`
}

func (server *Server) deleteTask(w http.ResponseWriter, r *http.Request) {
	var req deleteTaskRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponse(err))
		return
	}

	err = util.VerifyRequirements(req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponse(err))
		return
	}

	err = server.db.DeleteTask(context.Background(), req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write(errorResponse(err))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorResponse(err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
