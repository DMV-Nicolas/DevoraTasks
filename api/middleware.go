package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/DMV-Nicolas/DevoraTasks/token"
	ctx "github.com/gorilla/context"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(nextHandler http.HandlerFunc, tokenMaker token.Maker) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := fmt.Errorf("authorization header is not provided")
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := fmt.Errorf("invalid authorization header format")
			http.Error(w, err.Error(), http.StatusUnauthorized)
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			err := fmt.Errorf("unsupported authorization type %s", authorizationType)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx.Set(r, authorizationPayloadKey, payload)
		nextHandler.ServeHTTP(w, r)
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:52330")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			next.ServeHTTP(w, req)
		})
}
