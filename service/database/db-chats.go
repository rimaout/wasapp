package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
)

// CreateChat inserts a new chat row and returns the chat ID.
// For group chats, groupName must be non-empty and groupImagePath is optional.
// For private chats use empty string for groupName and groupImagePath.
func (db *appdbimpl) CreateChat(isGroup bool, groupNameIn string, groupImagePathIn string) (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("error generating chat UUID: %w", err)
	}
	chatId := id.String()

	// For group chats, groupNameDB must be non-empty
	var groupNameDB *string
	if groupNameIn != "" {
		groupNameDB = &groupNameIn
	}

	// For group chats, groupImagePath is optional
	var groupImagePathDB *string
	if groupImagePathIn != "" {
		groupImagePathDB = &groupImagePathIn
	}

	// Insert the new chat into the database (the constrains group constraints are defined in the database schema)
	_, err = db.c.Exec(
		`INSERT INTO chats (id, is_group_chat, group_name, group_image_path) VALUES (?, ?, ?, ?)`,
		chatId, isGroup, groupNameDB, groupImagePathDB,
	)
	if err != nil {
		return "", fmt.Errorf("error inserting chat into database: %w", err)
	}

	return chatId, nil
}

func (db *appdbimpl) FindPrivateChatBetween(userId1 string, userId2 string) (string, error) {
	var chatId string
	err := db.c.QueryRow(
		`SELECT c.id FROM chats c
			JOIN members m1 ON c.id = m1.chat_id
			JOIN members m2 ON c.id = m2.chat_id
		 WHERE c.is_group_chat = 0
		   AND m1.user_id = ?
		   AND m2.user_id = ?
		 LIMIT 1`,
		userId1, userId2,
	).Scan(&chatId)

	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("finding private chat: %w", err)
	}

	return chatId, nil
}
