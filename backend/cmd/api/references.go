// Filename: MyReference/backend/cmd/api/references.go
package main

import (
	"fmt"
	"net/http"

	"mgomez.net/internal/data"
	"mgomez.net/internal/validator"
)

// createReferenceHandler() will create a instance of a reference for the user
func (app *application) createdReferenceHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Location string `json:"storage-location"`
	}

	//pulling info from the json
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	//copying the information over from the json request
	reference := &data.Reference{
		Name:     input.Name,
		Location: input.Location,
	}

	//creating the validator
	v := validator.New()

	if data.ValidateReference(v, reference); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//Creating a location header for the new created resource / Reference
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/reference/%d", reference.ID))

	//writing the response
	err = app.writeJSON(w, http.StatusCreated, envelope{"reference": reference}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
