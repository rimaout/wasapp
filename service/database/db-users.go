package database

import (
	"database/sql"
	"errors"
	"time"
	"fmt"

	"github.com/gofrs/uuid"
)

func (db *appdbimpl) CreateUser(userName string) (string, error) {
    // Generate a new randomly-generated UUID
    id, err := uuid.NewV4()
    if err != nil {
        return "", fmt.Errorf("generating user UUID: %w", err)
    }

    userId := id.String()

    // Insert the new user. If a UUID collision happens (virtually impossible)
    // or if the DB connection fails, return a error.
    _, err = db.c.Exec("INSERT INTO users (id, name) VALUES (?, ?)", userId, userName)
    if err != nil {
        return "", fmt.Errorf("inserting user into database: %w", err)
    }

    return userId, nil
}

func (db *appdbimpl) GetUserIdByName(userName string) (string, error) {
	var userId string
	err := db.c.QueryRow("SELECT id FROM users WHERE name = ?", userName).Scan(&userId)

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

// GenerateUserSessionToken creates a new session token for the given user,
// saves it to the database, and returns the plain token string.
//
// NOTE: This implementation supports multi-device sessions. A user
// can log in from multiple devices simultaneously, and each device
// will receive and maintain its own valid session token.
func (db *appdbimpl) GenerateUserSessionToken(userId string) (string, error) {

	// Generate a new random UUID for the token
    tokenUUID, err := uuid.NewV4()
    if err != nil {
        return "", fmt.Errorf("generating session token: %w", err)
    }
    token := tokenUUID.String()

    // Set an expiration time (24 hours from now)
    expiresAt := time.Now().Add(24 * time.Hour)

    // Insert the token into the database
    _, err = db.c.Exec(`INSERT INTO tokens (token, user_id, expires_at) VALUES (?, ?, ?)`,
        token, userId, expiresAt)
    if err != nil {
        return "", fmt.Errorf("inserting session token into database: %w", err)
    }

    return token, nil
}
