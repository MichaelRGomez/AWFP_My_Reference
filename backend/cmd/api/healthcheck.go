//Filename: MyReferences/backend/cmd/api/healthcheck.go

package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	//create a map that'll hold healtcheck data
	data := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	//covering map -> JSON object
	err := app.writeJSON(w, http.StatusOK, data, nil)

	//will print error if any
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
