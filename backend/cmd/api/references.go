// Filename: MyReference/backend/cmd/api/references.go
package main

import (
	"errors"
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
	headers.Set("Location", fmt.Sprintf("/v1/references/%d", reference.ID))

	//writing the response
	err = app.writeJSON(w, http.StatusCreated, envelope{"reference": reference}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

// showReferenceHandler() retrieves a given reference when supplied with a reference id
func (app *application) showReferenceHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	reference, err := app.models.Reference.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	//writing the json response
	err = app.writeJSON(w, http.StatusOK, envelope{"reference": reference}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateReferenceHandler() allows a user to edit the name of a reference; for now
func (app *application) updateReferenceHandler(w http.ResponseWriter, r *http.Request) {
	//pulling the id from the json request
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	//trying to retrieve the reference from the database
	reference, err := app.models.Reference.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	//constructing a new version of the reference
	var input struct {
		Name *string `json:"name"`
	}

	//copying over the info from the edit request
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	//checking for updates
	if input.Name != nil {
		reference.Name = *input.Name
	}

	//validating the chagnes
	v := validator.New()
	if data.ValidateReference(v, reference); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//updating the reference on the database
	err = app.models.Reference.Update(reference)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	//writing the new info to a get() response to display to the caller
	err = app.writeJSON(w, http.StatusOK, envelope{"reference": reference}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// deleteReferenceHandler() will delete a given reference depending on the id provided
func (app *application) deleteReferenceHandler(w http.ResponseWriter, r *http.Request) {
	//pulling the id from the json request
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	//attemping to delete the reference if the id exist on the database
	err = app.models.Reference.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	//providing a confirmation response to the user
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "reference sucessfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
