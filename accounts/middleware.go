package accounts

import (
	"context"
	"net/http"
)

type MyToken struct {
	token string
}

// This middleware authenticates a user to access resources on a protected route
func Authenticator(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// get the token
		token := r.Header.Get("Authorization")

		ctx := r.Context()

		ctx = context.WithValue(ctx, "token", token)

		myReq := r.WithContext(ctx)

		// if token is not there just pass the request to the next middleware
		if token == "" {
			next.ServeHTTP(w, r)
		}

		next.ServeHTTP(w, myReq)
	})
}
