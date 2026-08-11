package database

import (
	"fmt"
	"time"
)

// InsertReceiverStatuses inserts receiver statuses for a given message and chat, for each active member of the chat except the sender.
//It sets the recv_time and read_time to NULL, this equates to the message being delivered but not yet received or read by the recipients.
func (db *appdbimpl) InsertReceiverStatuses(messageId, chatId, senderId string) error {
	now := time.Now().UTC().Format(time.RFC3339)

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

	_ = now
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

// ComputeMessageStatus computes the overall status of a message based on the receiver statuses.
// If all receivers have recv_time set, the message is considered "received".
// If at least one receiver has recv_time set, but not all, the message is considered "delivered".
// If no receivers have recv_time set, the message is considered "delivered".
func (db *appdbimpl) ComputeMessageStatus(messageId string) (MessageStatus, error) {
	var status string
	err := db.c.QueryRow(
		`SELECT CASE
			WHEN COUNT(*) = 0 THEN 'delivered'
			WHEN SUM(CASE WHEN recv_time IS NULL THEN 1 ELSE 0 END) > 0 THEN 'delivered'
			ELSE 'received'
		END
		FROM receiver_statuses WHERE message_id = ?`,
		messageId,
	).Scan(&status)
	if err != nil {
		return StatusDelivered, fmt.Errorf("computing message status: %w", err)
	}
	return MessageStatus(status), nil
}
