package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/DMV-Nicolas/DevoraTasks/token"
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

		fmt.Println(payload)
		nextHandler.ServeHTTP(w, r)
	})
}
