// Filename: MyReference/backend/internal/data/references.go
package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"mgomez.net/internal/validator"
)

type Reference struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Name      string    `json:"name"`
	Location  string    `json:"storage-location"`
	Version   int32     `json:"version"`
}

// validation for reference input
func ValidateReference(v *validator.Validator, reference *Reference) {
	//using check() to verify the data going into the input
	v.Check(reference.Name != "", "name", "must be provided")
	v.Check(len(reference.Name) <= 200, "name", "must no be more than 200 characters long")

	//Verification of the location will be done in the future
}

// Defining the model struct for reference
type ReferenceModel struct {
	DB *sql.DB
}

// CRUD functions
// Insert (Create)
func (m ReferenceModel) Insert(reference *Reference) error {
	query := `
		insert into reference_info (name, location)
		values ($1, $2)
		returning id, created_at, version
	`

	//preparing the arguments
	args := []interface{}{
		reference.Name, reference.Location,
	}
	//creating the context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&reference.ID, &reference.CreatedAt, &reference.Version)
}

// Get (Read)
func (m ReferenceModel) Get(id int64) (*Reference, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	//Creating the query
	query := `
		select id, created_at, name, location, version
		from reference_info
		where id = $1
	`
	//creating an instance to hold the info
	var reference Reference

	//dealing with context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&reference.ID,
		&reference.CreatedAt,
		&reference.Name,
		&reference.Location,
		&reference.Version,
	)

	//checking for errors
	if err != nil {
		//Check the type of error
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &reference, nil
}

// Update
func (m ReferenceModel) Update(reference *Reference) error {
	query := `
		update reference_info
		set name = $1, location = $2, version = version + 1
		where id = $3
		and version = $4
		returning version
	`
	args := []interface{}{
		reference.Name,
		reference.Location,
		reference.ID,
		reference.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&reference.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

// Delete
func (m ReferenceModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
		delete from reference_info
		where id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
