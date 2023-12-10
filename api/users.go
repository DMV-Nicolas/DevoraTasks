package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	db "github.com/DMV-Nicolas/DevoraTasks/db/sqlc"
	"github.com/DMV-Nicolas/DevoraTasks/util"
)

type createUserRequest struct {
	Username string `json:"username" requirements:"required"`
	Email    string `json:"email" requirements:"required;email"`
	Password string `json:"password" requirements:"required;min=8"`
}

func (server *Server) createUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	err := util.ShouldBindJSON(r, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorResponse(err))
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}

	user, err := server.db.CreateUser(context.Background(), arg)
	fmt.Println(err)
	if err != nil {
		if db.ErrorCode(err) == db.UniqueViolation {
			w.WriteHeader(http.StatusForbidden)
			w.Write(errorResponse(err))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorResponse(err))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse(user))
}

type getUserRequest struct {
	ID int64 `uri:"id" requirements:"required;min=1"`
}

func (server *Server) getUser(w http.ResponseWriter, r *http.Request) {
	var req getUserRequest
	err := util.ShouldBindUri(r, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponse(err))
		return
	}

	user, err := server.db.GetUser(context.Background(), req.ID)
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
	w.Write(jsonResponse(user))
}
