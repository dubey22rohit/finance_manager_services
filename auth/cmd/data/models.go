package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const dbTimeout = 3 * time.Second

var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool
	return Models{User: User{}}
}

type Models struct {
	User User
}

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// * Gets all the users registered with the app
func (u *User) GetAll() ([]*User, error) {
	query := `SELECT id, email, first_name, last_name, created_at, updated_at FROM users ORDER BY last_name`

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

// * Gets one user by email
func (u *User) GetByEmail(email string) (*User, error) {
	query := `SELECT id, email, first_name, last_name, created_at, updated_at FROM users WHERE email = $1 ORDER BY last_name`

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var user User

	row := db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// * Gets one user by user ID
func (u *User) GetByID(id int64) (*User, error) {
	query := `SELECT id, email, first_name, last_name, created_at, updated_at FROM users WHERE id = $1 ORDER BY last_name`

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var user User

	row := db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// * Updates one user in the database, using the info stored in the receiver u
func (u *User) Update() error {
	query := `UPDATE users SET email = $1, first_name = $2, last_name = $3, updated_at = $4 WHERE id = $5`

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	_, err := db.ExecContext(ctx, query, u.Email, u.FirstName, u.LastName, time.Now(), u.ID)
	if err != nil {
		return err
	}

	return nil
}

// *Deletes one user in the database, based on the ID provided
func (u *User) Delete(id int64) error {
	query := `DELETE FROM users where id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	_, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}

// *Inserts a new user in the database
func (u *User) Insert(user User) (int64, error) {
	query := `INSERT INTO users (email, first_name, last_name, password, created_at, updated_at) VALUES
	($1, $2, $3, $4, $5, $6) returning id`

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return 0, nil
	}

	var newID int64

	err = db.QueryRowContext(ctx, query,
		user.Email,
		user.FirstName,
		user.LastName,
		hashedPassword,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}

	return newID, nil
}

func (u *User) ResetPassword(password string) error {
	query := `UPDATE users SET password = $1 WHERE id = $2`

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	_, err = db.ExecContext(ctx, query, hashedPassword, u.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) PasswordMatches(plainText string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(plainText))
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
