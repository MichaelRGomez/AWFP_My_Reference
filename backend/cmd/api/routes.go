// Filename: MyReference/backend/cmd/api/routes.go
package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	//default endpoints
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	//user related endpoints
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activationUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	//MyReference related endpoints
	router.HandlerFunc(http.MethodPost, "/v1/references", app.requirePermission("reference:write", app.createdReferenceHandler))
	router.HandlerFunc(http.MethodGet, "/v1/references/:id", app.requirePermission("reference:read", app.showReferenceHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/references/:id", app.requirePermission("reference:write", app.updateReferenceHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/references/:id", app.requirePermission("reference:write", app.deleteReferenceHandler))

	//middleware chain
	return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))
}
