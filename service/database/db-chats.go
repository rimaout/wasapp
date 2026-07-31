package database

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/gofrs/uuid"
)

var ErrChatNotFound = errors.New("chat not found")

type Chat struct {
	Id             string  `json:"id"`
	IsGroupChat    bool    `json:"isGroupChat"`
	GroupName      *string `json:"groupName,omitempty"`
	GroupImagePath *string `json:"groupImagePath,omitempty"`
}
// Note: for GroupName & GroupImagePath we use `*string` (pointer to string) instad of `string`
//		 because `*string` can be null.

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

// GetChatById returns the chat struct with the given ID.
func (db *appdbimpl) GetChatById(chatId string) (Chat, error) {
	var chat Chat

	err := db.c.QueryRow(
		`SELECT id, is_group_chat, group_name, group_image_path FROM chats WHERE id = ?`,
		chatId,
	).Scan(&chat.Id, &chat.IsGroupChat, &chat.GroupName, &chat.GroupImagePath)

	if errors.Is(err, sql.ErrNoRows) {
		return Chat{}, ErrChatNotFound
	}

	if err != nil {
		return Chat{}, fmt.Errorf("querying chat: %w", err)
	}

	return chat, nil
}

// IsGroupChat returns whether the given chat is a group chat.
func (db *appdbimpl) IsGroupChat(chatId string) (bool, error) {
	var isGroup bool
	err := db.c.QueryRow(
		`SELECT is_group_chat FROM chats WHERE id = ?`,
		chatId,
	).Scan(&isGroup)

	if errors.Is(err, sql.ErrNoRows) {
		return false, ErrChatNotFound
	}
	if err != nil {
		return false, fmt.Errorf("querying chat type: %w", err)
	}

	return isGroup, nil
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

// SetGroupName updates the group name for a group chat.
func (db *appdbimpl) SetGroupName(chatId string, name string) error {
	_, err := db.c.Exec(
		`UPDATE chats SET group_name = ? WHERE id = ?`,
		name, chatId,
	)
	if err != nil {
		return fmt.Errorf("updating group name: %w", err)
	}
	return nil
}

// SetGroupAvatarPath updates the group image path for a group chat.
func (db *appdbimpl) SetGroupAvatarPath(chatId string, path string) error {
	_, err := db.c.Exec(
		`UPDATE chats SET group_image_path = ? WHERE id = ?`,
		path, chatId,
	)
	if err != nil {
		return fmt.Errorf("updating group image path: %w", err)
	}
	return nil
}

// GetGroupAvatarPath returns the group image path for a group chat.
// If the chat does not have an image, it returns an empty string and no error.
func (db *appdbimpl) GetGroupAvatarPath(chatId string) (string, error) {
	var path *string
	err := db.c.QueryRow(
		`SELECT group_image_path FROM chats WHERE id = ?`,
		chatId,
	).Scan(&path)

	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrChatNotFound
	}
	if err != nil {
		return "", fmt.Errorf("querying group image path: %w", err)
	}

	if path == nil {
		return "", nil
	}

	return *path, nil
}
