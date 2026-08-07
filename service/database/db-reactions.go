package database

import (
	"fmt"
	"time"
)

// getReactionsForMessage fetches all reactions for a message.
func (db *appdbimpl) getReactionsForMessage(messageId string) ([]EmojiReaction, error) {

	// Query the reactions for the given message ID
	rows, err := db.c.Query("SELECT user_id, emoji_id FROM reactions WHERE message_id = ?", messageId)
	if err != nil {
		return nil, fmt.Errorf("querying reactions: %w", err)
	}
	defer rows.Close()

	// Initialize a slice to hold the reactions
	reactions := []EmojiReaction{}

	// Iterate over the rows and scan each reaction into the slice
	for rows.Next() {
		var r EmojiReaction
		if err := rows.Scan(&r.UserID, &r.EmojiID); err != nil {
			return nil, fmt.Errorf("scanning reaction row: %w", err)
		}
		reactions = append(reactions, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating reaction rows: %w", err)
	}

	return reactions, nil
}

func (db *appdbimpl) CreateReaction(messageId, userId string, emojiId int32) error {
	_, err := db.c.Exec(
		`INSERT OR REPLACE INTO reactions (message_id, user_id, reac_time, emoji_id)
		 VALUES (?, ?, ?, ?)`,
		messageId, userId, time.Now().UTC().Format(time.RFC3339), emojiId,
	)
	if err != nil {
		return fmt.Errorf("creating reaction: %w", err)
	}
	return nil
}
