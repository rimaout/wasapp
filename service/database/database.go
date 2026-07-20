/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	CreateUser(userName string) (string, error)
	GetUserIdByName(userName string) (string, error)
	GetUserNameById(userId string) (string, error)
	SetUserName(userId string, newName string) error
	GenerateUserSessionToken(userId string) (string, error)
	GetUserIDByToken(token string) (string, error)
	SetUserAvatarPath(userId string, imagePath string) error
	GetUserAvatarPath(userId string) (string, error)
	Ping() error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")
	}

	// Enable Foreign Keys
	_, err := db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return nil, fmt.Errorf("error enabling foreign keys: %w", err)
	}

	// Check if a Table exists. If not, the database is empty, and we need to create the structure
	var tableName string
	err = db.QueryRow(`SELECT name FROM sqlite_master WHERE type='table' AND name='users';`).Scan(&tableName)
	if errors.Is(err, sql.ErrNoRows) {

		// --- USERS TABLE
		usersStmt := `CREATE TABLE users (
			"id"  TEXT NOT NULL PRIMARY KEY, -- UUID

			"name"               TEXT NOT NULL UNIQUE,
			"avatar_image_path"  TEXT
		);`
		_, err = db.Exec(usersStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating users table: %w", err)
		}

		// --- TOKENS TABLE
		tokensStmt := `CREATE TABLE tokens (
			"token"       TEXT NOT NULL PRIMARY KEY,
			"user_id"     TEXT NOT NULL,
			"expires_at"  DATETIME NOT NULL,

			FOREIGN KEY (user_id) REFERENCES users(id)
		)`
		_, err = db.Exec(tokensStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating tokens table: %w", err)
		}

		// ---  CHATS TABLE (for private and group chats)
		chatsStmt := `CREATE TABLE chats (
			"id"  TEXT NOT NULL PRIMARY KEY, -- UUID

			"group_name"        TEXT,    -- Optional (NULL for private chats)
			"group_image_path"  TEXT     -- Optional (NULL for private chats)
		);`
		_, err = db.Exec(chatsStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating chats table: %w", err)
		}

		// --- MEMBERS TABLE (for chat members)
		membersStmt := `CREATE TABLE members (
			"chat_id"  TEXT NOT NULL,
			"user_id"  TEXT NOT NULL,

			"group_join_time"   DATETIME,		-- Optional (NULL for private chats)
			"group_leave_time"	DATETIME,		-- Optional (NULL if the user is active or for private chats)

			PRIMARY KEY (chat_id, user_id),
			FOREIGN KEY (chat_id) REFERENCES chats(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		);`
		_, err = db.Exec(membersStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating members table: %w", err)
		}

		// --- MESSAGES TABLE
		messagesStmt := `CREATE TABLE messages (
			"id"              TEXT NOT NULL PRIMARY KEY, -- UUID

			"chat_id"         TEXT NOT NULL,
			"sender_id"       TEXT NOT NULL,
			"send_time"       DATETIME NOT NULL,
			"is_deleted"      BOOLEAN NOT NULL,
			"is_init_message" BOOLEAN NOT NULL,

			"text"      TEXT,
			"image_id"  TEXT, -- Optional, points to images(id)

			"reply_to_msg_id"         TEXT,  -- Self-referencing for replies

			"is_forward_message"      BOOLEAN NOT NULL,
			"forwarded_from_chat_id"  TEXT,
			"forwarded_from_msg_id"   TEXT,

			FOREIGN KEY (chat_id) REFERENCES chats(id),
			FOREIGN KEY (sender_id) REFERENCES users(id),
			FOREIGN KEY (image_id) REFERENCES images(id),
			FOREIGN KEY (reply_to_msg_id) REFERENCES messages(id),
			FOREIGN KEY (forwarded_from_chat_id) REFERENCES chats(id),
			FOREIGN KEY (forwarded_from_msg_id) REFERENCES messages(id)
		);`
		_, err = db.Exec(messagesStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating messages table: %w", err)
		}

		// --- IMAGES TABLE (for message images)
		imagesStmt := `CREATE TABLE images (
			"id"   TEXT NOT NULL PRIMARY KEY,
			"path" TEXT NOT NULL
		);`
		_, err = db.Exec(imagesStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating images table: %w", err)
		}

	    // --- RECEIVER_STATUSES TABLE (for message delivery and read statuses)
		receiverStatusesStmt := `CREATE TABLE receiver_statuses (
			"message_id"  TEXT NOT NULL,
			"user_id"     TEXT NOT NULL,

			"recv_time"   DATETIME,            -- Optional (until delivered)
			"read_time"   DATETIME,            -- Optional (until read)

			PRIMARY KEY (message_id, user_id),
			FOREIGN KEY (message_id) REFERENCES messages(id),
			FOREIGN KEY (user_id) REFERENCES users(id),

		);`
		_, err = db.Exec(receiverStatusesStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating receiver_statuses table: %w", err)
		}

		// --- REACTIONS TABLE (for message reactions)
		reactionsStmt := `CREATE TABLE reactions (
			"message_id" TEXT NOT NULL,
			"user_id"    TEXT NOT NULL,

			"emoji_id"   INTEGER NOT NULL CHECK (emoji_id >= 0 AND emoji_id <= 9),

			PRIMARY KEY (message_id, user_id),
			FOREIGN KEY (message_id) REFERENCES messages(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		);`
		_, err = db.Exec(reactionsStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating reactions table: %w", err)
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
