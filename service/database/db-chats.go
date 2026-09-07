package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

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
// 		 because `*string` can be null.

// CreateChat inserts a new chat row and returns the chat ID.
// For group chats, groupName must be non-empty.
// For private chats use empty string for groupName.
func (db *appdbimpl) CreateChat(isGroup bool, groupNameIn string) (string, error) {

	// Generate a new UUID for the chat ID
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

	// Insert the new chat into the database (the constrains are defined in the database schema)
	_, err = db.c.Exec(
		`INSERT INTO chats (id, is_group_chat, group_name) VALUES (?, ?, ?)`,
		chatId, isGroup, groupNameDB,
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

// FindPrivateChatBetween returns the chat ID of a private chat between two users, if it exists.
// If no such chat exists, it returns an empty string and no error.
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

type ChatPreview struct {
	ChatID      string
	IsGroupChat bool
	DisplayName string
	MessageID   string
	SendTime    time.Time
	SenderID    string
	SenderName  string
	IsDeleted   bool
	IsInitMsg   bool
	IsJoinMsg   bool
	IsLeaveMsg  bool
	IsForward   bool
	Text        *string
	ImageID     *string
	ReplyTo     *string
	Status      string
	UnreadCount int
}

// GetMyChats returns a list of chats previews for the given user, including the last message in each chat.
func (db *appdbimpl) GetMyChats(userId string) ([]ChatPreview, error) {
	rows, err := db.c.Query(
		`SELECT
			c.id, c.is_group_chat, c.group_name,

			-- Determine the display name based on if it's a group chat or a private chat
			--  - For group chats, the group name is used.
			--  - For private chats, the name of the other user is used.

			CASE WHEN c.is_group_chat = 1 THEN c.group_name
				 ELSE (SELECT u2.name
					   FROM members m2 JOIN users u2 ON m2.user_id = u2.id
					   WHERE m2.chat_id = c.id AND m2.user_id != ? LIMIT 1)
			END AS display_name,

			m.id, m.send_time, m.sender_id, u.name AS sender_name,
			m.is_deleted, m.is_init_message, m.is_join_message, m.is_leave_message,
			m.is_forward_message,
			m.text, m.image_id, m.reply_to_msg_id,

			-- Count unread messages for this user in this chat
			(SELECT COUNT(*)
			 FROM receiver_statuses rs
			 JOIN messages msg ON rs.message_id = msg.id
			 WHERE msg.chat_id = c.id AND rs.user_id = ? AND rs.read_time IS NULL
			) AS unread_count,

			-- Compute the status of the last message (messageStatusCase is defined in db-chats.go)
			(SELECT `+
			messageStatusCase+
			`
			FROM receiver_statuses rs WHERE rs.message_id = m.id) AS message_status

		-- Get the last message for each chat
		FROM members mem
		JOIN chats c ON mem.chat_id = c.id
		JOIN messages m ON m.id = (
			-- send_time has second precision, so use rowid as a tiebreaker
			-- to reliably pick the most recently sent message.
			SELECT id FROM messages WHERE chat_id = c.id ORDER BY send_time DESC, rowid DESC LIMIT 1
		)

		-- Get the sender's name for the last message
		JOIN users u ON m.sender_id = u.id

		-- Only include chats the user is currently a member of
		WHERE mem.user_id = ? AND mem.group_leave_time IS NULL

		-- Order chats by the last message's send time
		ORDER BY m.send_time DESC, m.rowid DESC

		-- Limit the number of chats returned (in future, I would like implement pagination)
		LIMIT 50`,
		userId, userId, userId,
	)
	if err != nil {
		return nil, fmt.Errorf("querying chats: %w", err)
	}
	defer rows.Close()

	// Extract the chat previews from the rows
	chats := make([]ChatPreview, 0)
	for rows.Next() {
		var cp ChatPreview
		var groupName *string
		err := rows.Scan(
			&cp.ChatID, &cp.IsGroupChat, &groupName,
			&cp.DisplayName,
			&cp.MessageID, &cp.SendTime, &cp.SenderID, &cp.SenderName,
			&cp.IsDeleted, &cp.IsInitMsg, &cp.IsJoinMsg, &cp.IsLeaveMsg,
			&cp.IsForward,
			&cp.Text, &cp.ImageID, &cp.ReplyTo,
			&cp.UnreadCount, &cp.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning chat row: %w", err)
		}
		chats = append(chats, cp)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating chat rows: %w", err)
	}

	return chats, nil
}
