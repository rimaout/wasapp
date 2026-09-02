package database

import (
	"fmt"
	"time"
)

// InsertReceiverStatuses inserts receiver statuses for a message and chat, for each active member of the chat except the sender.
// It sets the recv_time and read_time to NULL, this means that the message is delivered but not yet received or read by the recipients.
func (db *appdbimpl) InsertReceiverStatuses(messageId, chatId, senderId string) error {
	_, err := db.c.Exec(
		`INSERT INTO receiver_statuses (message_id, user_id, recv_time, read_time)
		 SELECT ?, m.user_id, NULL, NULL
		 FROM members m
		 WHERE m.chat_id = ? AND m.user_id != ? AND m.group_leave_time IS NULL`,
		messageId, chatId, senderId,
	)
	if err != nil {
		return fmt.Errorf("inserting receiver statuses: %w", err)
	}

	return nil
}

// MarkMessagesReceivedByUser marks all messages of the user as received, setting recv_time to the
// current time where it is still NULL. This means the user's device has received the messages, but
// has not yet read them.
// If chatId is non-empty, only messages in that chat are marked; otherwise all the user's messages
// are marked (used when the client delivers the whole chat list).
func (db *appdbimpl) MarkMessagesReceivedByUser(chatId, userId string) error {
	now := time.Now().UTC().Format(time.RFC3339)

	query := `UPDATE receiver_statuses
		SET recv_time = ?
		WHERE user_id = ?
		  AND recv_time IS NULL`
	args := []interface{}{now, userId}
	if chatId != "" {
		query += ` AND message_id IN (SELECT id FROM messages WHERE chat_id = ?)`
		args = append(args, chatId)
	}

	_, err := db.c.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("marking messages received: %w", err)
	}
	return nil
}

// MarkMessagesReadByUser marks all messages in a chat as read by the given user.
func (db *appdbimpl) MarkMessagesReadByUser(chatId, userId string) error {
	now := time.Now().UTC().Format(time.RFC3339)

	_, err := db.c.Exec(
		`UPDATE receiver_statuses
		 SET read_time = ?
		 WHERE message_id IN (SELECT id FROM messages WHERE chat_id = ?)
		   AND user_id = ?
		   AND read_time IS NULL
		   AND recv_time IS NOT NULL`,
		now, chatId, userId,
	)
	if err != nil {
		return fmt.Errorf("marking messages read: %w", err)
	}
	return nil
}

// messageStatusCase is the SQL CASE expression that maps a message's receiver statuses to
// DELIVERED/RECEIVED/READ. It is shared by ComputeMessageStatus and GetMyChats so the two stay in sync.
const messageStatusCase = `CASE
			WHEN COUNT(*) > 0 AND SUM(CASE WHEN read_time IS NULL THEN 1 ELSE 0 END) = 0 THEN 'READ'
			WHEN COUNT(*) > 0 AND SUM(CASE WHEN recv_time IS NULL THEN 1 ELSE 0 END) = 0 THEN 'RECEIVED'
			ELSE 'DELIVERED'
		END`

// ComputeMessageStatus computes the overall status of a message based on the receiver statuses.
// - "READ": All recipients have both received and read the message.
// - "RECEIVED": All recipients have received the message, but at least one has not read it yet (read_time IS NULL).
// - "DELIVERED": At least one recipient has not received the message yet (recv_time IS NULL), or zero recipients exist.
func (db *appdbimpl) ComputeMessageStatus(messageId string) (MessageStatus, error) {
	var status string
	err := db.c.QueryRow(
		`SELECT `+messageStatusCase+` FROM receiver_statuses WHERE message_id = ?`,
		messageId,
	).Scan(&status)
	if err != nil {
		return StatusDelivered, fmt.Errorf("computing message status: %w", err)
	}
	return MessageStatus(status), nil
}
