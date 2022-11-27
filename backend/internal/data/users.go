// Filename: MyReference/backend/internal/data/users.go
package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"mgomez.net/internal/validator"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
)

// Anonymous User with no properties available
var AnonymousUser = &User{}

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

// IsAnonymous() checks if a user is anonymous
func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

// creating a password type
type password struct {
	plaintext *string
	hash      []byte
}

// Set() method stores the hash of the plaintext password
func (p *password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}

	p.plaintext = &plaintextPassword
	p.hash = hash

	return nil
}

// Matches() checks if the supplied password is correct
func (p *password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

// ValidateEmail() validates the user email input
func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

// ValidatePasswordPlaintext() validates if the password follows the rules
func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be atleast 8 characters long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 characters long")
}

// ValidateUser() validates user input
func ValidateUser(v *validator.Validator, user *User) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 characters long")

	//Checking the email
	ValidateEmail(v, user.Email)

	//Checking the password
	if user.Password.plaintext != nil {
		ValidatePasswordPlaintext(v, *user.Password.plaintext)
	}

	//Ensuring a hash was created for the password
	if user.Password.hash == nil {
		panic("missing password hash for the user")
	}
}

// Creating our DB model
type UserModel struct {
	DB *sql.DB
}

// Insert() for user
func (m UserModel) Insert(user *User) error {
	query := `
		insert into users (name, email, password_hash, activated)
		values ($1, $2, $3, $4)
		returning id, created_at, version
	`
	args := []interface{}{
		user.Name,
		user.Email,
		user.Password.hash,
		user.Activated,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

// GetByEmail() retrieves a user based on email
func (m UserModel) GetByEmail(email string) (*User, error) {
	query := `
		select id, created_at, name, email, password_hash, activated, version
		from users
		where email = $1
	`
	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

// Update()
func (m UserModel) Update(user *User) error {
	query := `
		update users
		set name = $1, email = $2, password_hash = $3, activated = $4, version + 1
		where id = $5 and version = $6
		returning version
	`

	args := []interface{}{
		user.Name,
		user.Email,
		user.Password.hash,
		user.Activated,
		user.ID,
		user.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

// GetForToken()
func (m UserModel) GetForToken(tokenScope, TokenPlaintext string) (*User, error) {
	tokenHash := sha256.Sum256([]byte(TokenPlaintext))
	query := `
		select users.id, users.created_at, users.name, users.email, users.password_hash, users.activated, users.version
		from users
		inner join tokens
		on users.id = tokens.user_id
		where tokens.hash = $1
		and tokens.scope = $2
		and tokens.expiry > $3
	`
	args := []interface{}{tokenHash[:], tokenScope, time.Now()}

	var user User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}
