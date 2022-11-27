// Filename: MyReference/backend/cmd/api/context.go
package main

import (
	"context"
	"net/http"

	"mgomez.net/internal/data"
)

// Define a custome contextKey type
type contextKey string

//make user a key

const userContextkey = contextKey("user")

//Method to add user to the context

func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextkey, user)
	return r.WithContext(ctx)
}

// Retrieve the User struct
func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextkey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
