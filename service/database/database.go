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
	// Users
	CreateUser(userName string) (string, error)
	GetUserIdByName(userName string) (string, error)
	GetUserNameById(userId string) (string, error)
	SetUserName(userId string, newName string) error
	GenerateUserSessionToken(userId string) (string, error)
	GetUserIDByToken(token string) (string, error)
	SetUserAvatarPath(userId string, imagePath string) error
	GetUserAvatarPath(userId string) (string, error)
	SearchUsers(query string) ([]User, error)

	// Chats
	CreateChat(isGroup bool, groupName string, groupImagePath string) (string, error)
	GetChatById(chatId string) (Chat, error)
	IsGroupChat(chatId string) (bool, error)
	FindPrivateChatBetween(userId1 string, userId2 string) (string, error)
	SetGroupName(chatId string, name string) error
	SetGroupAvatarPath(chatId string, path string) error
	GetGroupAvatarPath(chatId string) (string, error)

	// Members
	AddChatMember(chatId string, userId string) error
	GetChatMembers(chatId string) ([]Member, error)
	GetOtherMemberId(chatId string, userId string) (string, error)
	SetLeaveTime(chatId string, userId string) error
	IsActiveChatMember(chatId string, userId string) (bool, error)

	// Messages
	CreateMessage(chatId, senderId, sendTime, text, imageId, replyTo string, isInit, isForward bool, forwardFromChat, forwardFromMsg string) (string, error)
	SaveMessageImage(path string) (string, error)
	DeleteMessageImage(imageId string) error

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

			"is_group_chat"     BOOLEAN NOT NULL,
			"group_name"        TEXT,    -- Optional (NULL for private chats)
			"group_image_path"  TEXT     -- Optional (NULL for private chats or for group chats without an image)

			CONSTRAINT chk_group_chat_fields CHECK (
				(is_group_chat = 1 AND group_name IS NOT NULL) OR
				(is_group_chat = 0 AND group_name IS NULL AND group_image_path IS NULL)
			)
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
			"image_id"  TEXT,

			"reply_to_msg_id"         TEXT,

			"is_forward_message"      BOOLEAN NOT NULL,
			"forwarded_from_chat_id"  TEXT,
			"forwarded_from_msg_id"   TEXT,

			UNIQUE(chat_id, id),
			FOREIGN KEY (chat_id) REFERENCES chats(id),
			FOREIGN KEY (sender_id) REFERENCES users(id),
			FOREIGN KEY (image_id) REFERENCES images(id),
			FOREIGN KEY (chat_id, reply_to_msg_id) REFERENCES messages(chat_id, id),
			FOREIGN KEY (forwarded_from_chat_id) REFERENCES chats(id),
			FOREIGN KEY (forwarded_from_msg_id) REFERENCES messages(id)

			-- Constraint 1.1: Minimum Content Requirement
			-- Must have text or image, unless it's an init message, deleted, or a forward
			CONSTRAINT chk_msg_minimum_content CHECK (
				(is_init_message = 1 OR is_deleted = 1 OR is_forward_message = 1) OR
				(text IS NOT NULL OR image_id IS NOT NULL)
			),

			-- Constraint 1.2: Initial Message Constraints
			-- If it's an init message, it cannot have text, images, be deleted, or be a forward
			CONSTRAINT chk_init_msg_fields CHECK (
				NOT (is_init_message = 1) OR
				(text IS NULL AND image_id IS NULL AND is_deleted = 0 AND is_forward_message = 0)
			),

			-- Constraint 1.3: Forwarded Message Parameters
			-- Either all 3 forward fields are active/filled, or all 3 are inactive/NULL
			CONSTRAINT chk_forward_pointers CHECK (
				(is_forward_message = 1 AND forwarded_from_chat_id IS NOT NULL AND forwarded_from_msg_id IS NOT NULL) OR
				(is_forward_message = 0 AND forwarded_from_chat_id IS NULL AND forwarded_from_msg_id IS NULL)
			),

			-- Constraint 1.4: Content of a Forwarded Message
			-- Forwarded messages cannot contain raw text or images directly, and cannot be an init message
			CONSTRAINT chk_forward_content_empty CHECK (
				NOT (is_forward_message = 1) OR
				(text IS NULL AND image_id IS NULL AND is_init_message = 0)
			),

			-- Constraint 1.5: Logical Deletion State Enforcement
			-- When deleted, all user data and forward links must be cleared. Init messages can't be deleted.
			CONSTRAINT chk_logical_deletion_state CHECK (
				NOT (is_deleted = 1) OR (
					text IS NULL AND
					image_id IS NULL AND
					is_forward_message = 0 AND
					forwarded_from_chat_id IS NULL AND
					forwarded_from_msg_id IS NULL AND
					is_init_message = 0
				)
			)

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

			-- Constraint 4.2: Message Timeline Coherence
			-- Chronological order: recv_time <= read_time (safely handling NULLs)
			CONSTRAINT chk_delivery_timeline CHECK (
				read_time IS NULL OR (recv_time IS NOT NULL AND recv_time <= read_time)
			)
		);`
		_, err = db.Exec(receiverStatusesStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating receiver_statuses table: %w", err)
		}

		// --- REACTIONS TABLE (for message reactions)
		reactionsStmt := `CREATE TABLE reactions (
			"message_id" TEXT NOT NULL,
			"user_id"    TEXT NOT NULL,

			"reac_time"  DATETIME NOT NULL,
			"emoji_id"	 INTEGER NOT NULL CHECK (emoji_id >= 0 AND emoji_id <= 9),

			PRIMARY KEY (message_id, user_id),
			FOREIGN KEY (message_id) REFERENCES messages(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		);`
		_, err = db.Exec(reactionsStmt)
		if err != nil {
			return nil, fmt.Errorf("error creating reactions table: %w", err)
		}

		// --- DATABASE TRIGGERS (cross-table constraints)
		if err = createTriggers(db); err != nil {
			return nil, err
		}
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
