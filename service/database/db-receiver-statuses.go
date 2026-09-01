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

// MarkMessagesReceivedByUser updates the receiver statuses for a given chat and user, setting the recv_time to the current time for all messages in that chat that have not yet been marked as received by that user.
// This means that the user has been notified of the messages, but has not yet read them.
func (db *appdbimpl) MarkMessagesReceivedByUser(chatId, userId string) error {
	now := time.Now().UTC().Format(time.RFC3339)

	_, err := db.c.Exec(
		`UPDATE receiver_statuses
		 SET recv_time = ?
		 WHERE message_id IN (SELECT id FROM messages WHERE chat_id = ?)
		   AND user_id = ?
		   AND recv_time IS NULL`,
		now, chatId, userId,
	)
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
		   AND read_time IS NULL`,
		now, chatId, userId,
	)
	if err != nil {
		return fmt.Errorf("marking messages read: %w", err)
	}
	return nil
}

// ComputeMessageStatus computes the overall status of a message based on the receiver statuses.
// - "delivered": Default state, at least one recipient has not received the message yet (recv_time IS NULL), or zero recipients exist.
// - "received": All recipients have received the message, but at least one has not read it yet (read_time IS NULL).
// - "read": All recipients have both received and read the message.
func (db *appdbimpl) ComputeMessageStatus(messageId string) (MessageStatus, error) {
	var status string
	err := db.c.QueryRow(
		`SELECT CASE
			WHEN COUNT(*) = 0 THEN 'delivered'
			WHEN SUM(CASE WHEN recv_time IS NULL THEN 1 ELSE 0 END) > 0 THEN 'delivered'
			WHEN SUM(CASE WHEN read_time IS NULL THEN 1 ELSE 0 END) > 0 THEN 'received'
			ELSE 'read'
		END
		FROM receiver_statuses WHERE message_id = ?`,
		messageId,
	).Scan(&status)
	if err != nil {
		return StatusDelivered, fmt.Errorf("computing message status: %w", err)
	}
	return MessageStatus(status), nil
}
