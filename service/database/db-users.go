package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
)

func (db *appdbimpl) CreateUser(name string) (string, error) {
    // Generate a new randomly-generated UUID
    id, err := uuid.NewV4()
    if err != nil {
        return "", fmt.Errorf("generating user UUID: %w", err)
    }

    userId := id.String()

    // Insert the new user. If a UUID collision happens (virtually impossible)
    // or if the DB connection fails, return a error.
    _, err = db.c.Exec("INSERT INTO users (id, name) VALUES (?, ?)", userId, name)
    if err != nil {
        return "", fmt.Errorf("inserting user into database: %w", err)
    }

    return userId, nil
}

func (db *appdbimpl) GetUserByName(name string) (string, error) {
	var userId string
	err := db.c.QueryRow("SELECT id FROM users WHERE name = ?", name).Scan(&userId)

	// If no user is found with the given name, return an empty string and no error
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}

	// If any other error occurs, return it
	if err != nil {
		return "", fmt.Errorf("querying user: %w", err)
	}

	// If a user is found, return the user ID
	return userId, nil
}
