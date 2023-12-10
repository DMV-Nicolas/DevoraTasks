package api

import (
	"context"
	"database/sql"
	"net/http"

	db "github.com/DMV-Nicolas/DevoraTasks/db/sqlc"
	"github.com/DMV-Nicolas/DevoraTasks/util"
)

type createTaskRequest struct {
	UserID      int64  `json:"user_id" requirements:"min=1"`
	Title       string `json:"title" requirements:"required"`
	Description string `json:"description"`
}

func (server *Server) createTask(w http.ResponseWriter, r *http.Request) {
	var req createTaskRequest
	err := util.ShouldBindJSON(r, &req)
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

	task, err := server.store.CreateTask(context.Background(), arg)
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
	Offset int32 `form:"offset" requirements:"min=0"`
	Limit  int32 `form:"limit" requirements:"min=1"`
}

func (server *Server) listTasks(w http.ResponseWriter, r *http.Request) {
	var req listTasksRequest
	err := util.ShouldBindQuery(r, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponse(err))
		return
	}

	arg := db.ListTasksParams{
		Offset: req.Offset,
		Limit:  req.Limit,
	}

	tasks, err := server.store.ListTasks(context.Background(), arg)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse(tasks))
}

type getTaskRequest struct {
	ID int64 `uri:"id" requirements:"min=1"`
}

func (server *Server) getTask(w http.ResponseWriter, r *http.Request) {
	var req getTaskRequest
	err := util.ShouldBindUri(r, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponse(err))
		return
	}

	task, err := server.store.GetTask(context.Background(), req.ID)
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
	err := util.ShouldBindJSON(r, &req)
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

	task, err := server.store.UpdateTask(context.Background(), arg)
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
	err := util.ShouldBindJSON(r, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponse(err))
		return
	}

	err = server.store.DeleteTask(context.Background(), req.ID)
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
