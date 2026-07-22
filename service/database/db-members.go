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

