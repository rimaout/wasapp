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

func (db *appdbimpl) GetUserNameById(userId string) (string, error) {
	var userName string
	err := db.c.QueryRow("SELECT name FROM users WHERE id = ?", userId).Scan(&userName)

	// If no user is found with the given id, return an empty string and no error
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}

	// If any other error occurs, return it
	if err != nil {
		return "", fmt.Errorf("querying user: %w", err)
	}

	// If a user is found, return the user name
	return userName, nil
}

func (db *appdbimpl) SetUserName(userId string, newName string) error {
	_, err := db.c.Exec("UPDATE users SET name = ? WHERE id = ?", newName, userId)
	if err != nil {
		return fmt.Errorf("updating user name: %w", err)
	}
	return nil
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

// GetUserIDByToken checks if the given token represents a live session.
// Returns the associated user ID if the token exists and has not yet expired.
// Returns an empty string (and no error) if the token is expired or invalid.
func (db *appdbimpl) GetUserIDByToken(token string) (string, error) {
	var userId string
	err := db.c.QueryRow("SELECT user_id FROM tokens WHERE token = ? AND expires_at > datetime('now')", token).Scan(&userId)

	if errors.Is(err, sql.ErrNoRows) {
		return "", nil // no match: token invalid or expired
	}

	if err != nil {
		return "", fmt.Errorf("querying token: %w", err) // DB error
	}

	return userId, nil
}

func (db *appdbimpl) SetUserAvatarPath(userId string, imagePath string) error {
	_, err := db.c.Exec("UPDATE users SET avatar_image_path = ? WHERE id = ?", imagePath, userId)
	if err != nil {
		return fmt.Errorf("updating user avatar image path: %w", err)
	}
	return nil
}

func (db *appdbimpl) GetUserAvatarPath(userId string) (string, error) {
	var path sql.NullString
	err := db.c.QueryRow("SELECT avatar_image_path FROM users WHERE id = ?", userId).Scan(&path)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("querying user avatar image path: %w", err)
	}
	if !path.Valid { // check if the path is NULL in the database
		return "", nil
	}
	return path.String, nil
}
