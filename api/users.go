package api

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	ctx "github.com/gorilla/context"

	db "github.com/DMV-Nicolas/DevoraTasks/db/sqlc"
	"github.com/DMV-Nicolas/DevoraTasks/token"
	"github.com/DMV-Nicolas/DevoraTasks/util"
)

type createUserRequest struct {
	Username string `json:"username" requirements:"required"`
	Email    string `json:"email" requirements:"required;email"`
	Password string `json:"password" requirements:"required;min=8"`
}

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
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
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}

	user, err := server.store.CreateUser(context.Background(), arg)
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

	res := newUserResponse(user)

	w.WriteHeader(http.StatusCreated)
	w.Write(jsonResponse(res))
}

func (server *Server) getUser(w http.ResponseWriter, r *http.Request) {
	payload := ctx.Get(r, authorizationPayloadKey).(*token.Payload)
	user, err := server.store.GetUser(context.Background(), payload.Username)
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

	res := newUserResponse(user)

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse(res))
}

type loginUserRequest struct {
	Username string `json:"username" requirements:"required"`
	Password string `json:"password" requirements:"required;min=8"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) loginUser(w http.ResponseWriter, r *http.Request) {
	var req loginUserRequest
	err := util.ShouldBindJSON(r, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errorResponse(err))
		return
	}

	user, err := server.store.GetUser(context.Background(), req.Username)
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

	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write(errorResponse(err))
		return
	}

	token, err := server.tokenMaker.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorResponse(err))
		return
	}

	res := loginUserResponse{
		AccessToken: token,
		User:        newUserResponse(user),
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse(res))
}
