package api

import (
	"context"
	"encoding/json"
	"net/http"

	db "github.com/DMV-Nicolas/DevoraTasks/db/sqlc"
	"github.com/DMV-Nicolas/DevoraTasks/util"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type createUserRequest struct {
	Username string `json:"username" requirements:"required"`
	Email    string `json:"email" requirements:"required;email"`
	Password string `json:"password" requirements:"required;min=8"`
}

func (server *Server) createUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
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

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		if err == bcrypt.ErrPasswordTooLong {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(errorResponse(err))
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorResponse(err))
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	}

	user, err := server.db.CreateUser(context.Background(), arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				w.WriteHeader(http.StatusForbidden)
				w.Write(errorResponse(err))
				return
			}
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorResponse(err))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse(user))
}
