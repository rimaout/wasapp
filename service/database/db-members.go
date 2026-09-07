package database

import (
	"fmt"
	"time"
)

type Member struct {
	UserId    string  `json:"userId"`
	Name      string  `json:"name"`
	JoinTime  *string `json:"joinTime,omitempty"`  // direct chat: NULL, group chat: timestamp
	LeaveTime *string `json:"leaveTime,omitempty"` // direct chat: NULL, group chat: timestamp or NULL if still active
}

// AddChatMember adds a user to a chat. For a fresh join it inserts a new row
// (join_time set to now, leave_time NULL); for a member who had previously left,
// it updates the existing row by updating join_time and clearing leave_time.
func (db *appdbimpl) AddChatMember(chatId string, userId string) error {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := db.c.Exec(
		`INSERT INTO members (chat_id, user_id, group_join_time, group_leave_time) VALUES (?, ?, ?, ?)
		 ON CONFLICT(chat_id, user_id)
		 DO UPDATE SET group_join_time = excluded.group_join_time, group_leave_time = NULL`,
		chatId, userId, now, nil,
	)
	if err != nil {
		return fmt.Errorf("error adding member to chat: %w", err)
	}
	return nil
}

// GetChatMembers returns all members of a chat (regardless of leave status).
func (db *appdbimpl) GetChatMembers(chatId string) ([]Member, error) {

	// Query members and join with users table to get user names
	rows, err := db.c.Query(
		`SELECT m.user_id, u.name, m.group_join_time, m.group_leave_time
		 FROM members m
		 JOIN users u ON m.user_id = u.id
		 WHERE m.chat_id = ?
		 ORDER BY u.name ASC`,
		chatId,
	)

	// Check for errors in the query
	if err != nil {
		return nil, fmt.Errorf("querying chat members: %w", err)
	}
	defer rows.Close()

	// Iterate over the rows and scan into Member structs
	members := make([]Member, 0)
	for rows.Next() {
		var m Member
		if err := rows.Scan(&m.UserId, &m.Name, &m.JoinTime, &m.LeaveTime); err != nil {
			return nil, fmt.Errorf("scanning member row: %w", err)
		}
		members = append(members, m)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating member rows: %w", err)
	}

	return members, nil
}

// GetOtherMemberId returns the user ID of the other member in a private chat.
func (db *appdbimpl) GetOtherMemberId(chatId string, userId string) (string, error) {
	var otherUserId string
	err := db.c.QueryRow(
		`SELECT user_id FROM members WHERE chat_id = ? AND user_id != ? LIMIT 1`,
		chatId, userId,
	).Scan(&otherUserId)

	if err != nil {
		return "", fmt.Errorf("querying other member: %w", err)
	}

	return otherUserId, nil
}

// SetLeaveTime marks the user as having left the group chat (leave_time field is set to now).
func (db *appdbimpl) SetLeaveTime(chatId string, userId string) error {
	now := time.Now().UTC().Format(time.RFC3339)

	_, err := db.c.Exec(
		`UPDATE members SET group_leave_time = ? WHERE chat_id = ? AND user_id = ?`,
		now, chatId, userId,
	)
	if err != nil {
		return fmt.Errorf("setting leave time: %w", err)
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
