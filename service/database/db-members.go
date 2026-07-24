package database

import (
	"database/sql"
	"fmt"
	"time"
)

// AddChatMember inserts a new row into the members table.
// For private chats, join_time and leave_time are NULL.
// For group chats, join_time is set to now, leave_time is NULL.
func (db *appdbimpl) AddChatMember(chatId string, userId string) error {
	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	_, err := db.c.Exec(
		`INSERT INTO members (chat_id, user_id, group_join_time, group_leave_time) VALUES (?, ?, ?, ?)`,
		chatId, userId, now, nil,
	)
	if err != nil {
		return fmt.Errorf("error adding member to chat: %w", err)
	}
	return nil
}


// IsActiveChatMember checks if the user is an active member of the chat.
func (db *appdbimpl) IsActiveChatMember(chatId string, userId string) (bool, error) {
	var exists int
	err := db.c.QueryRow(
		`SELECT COUNT(*) FROM members
		 WHERE chat_id = ? AND user_id = ? AND group_leave_time IS NULL`,
		chatId, userId,
	).Scan(&exists)

	if err != nil {
		return false, fmt.Errorf("checking chat membership: %w", err)
	}

	return exists > 0, nil
}

