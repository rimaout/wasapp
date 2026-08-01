package database

import (
	"fmt"

	"github.com/gofrs/uuid"
)

// CreateMessage inserts a new message and returns the message ID.
// Pass empty strings for optional fields that should be NULL.
func (db *appdbimpl) CreateMessage(
	chatId, senderId, sendTime, inText, inImageId, inReplyTo string,
	isInit, isForward bool,
	inForwardFromChat, inForwardFromMsg string,
) (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("error generating message UUID: %w", err)
	}
	msgId := id.String()

	var dbText *string
	if inText != "" {
		dbText = &inText
	}
	var dbImageId *string
	if inImageId != "" {
		dbImageId = &inImageId
	}
	var dbReplyTo *string
	if inReplyTo != "" {
		dbReplyTo = &inReplyTo
	}
	var dbForwardFromChat *string
	if inForwardFromChat != "" {
		dbForwardFromChat = &inForwardFromChat
	}
	var dbForwardFromMsg *string
	if inForwardFromMsg != "" {
		dbForwardFromMsg = &inForwardFromMsg
	}

	_, err = db.c.Exec(
		`INSERT INTO messages
			(id, chat_id, sender_id, send_time, is_deleted, is_init_message,
			 text, image_id, reply_to_msg_id,
			 is_forward_message, forwarded_from_chat_id, forwarded_from_msg_id)
		 VALUES (?, ?, ?, ?, 0, ?, ?, ?, ?, ?, ?, ?)`,
		msgId, chatId, senderId, sendTime, isInit,
		dbText, dbImageId, dbReplyTo,
		isForward, dbForwardFromChat, dbForwardFromMsg,
	)
	if err != nil {
		return "", fmt.Errorf("inserting message into database: %w", err)
	}

	return msgId, nil
}

// SaveMessageImage inserts an image path into the images table and returns its generated UUID.
func (db *appdbimpl) SaveMessageImage(path string) (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("error generating image UUID: %w", err)
	}
	imageId := id.String()

	_, err = db.c.Exec(
		`INSERT INTO images (id, path) VALUES (?, ?)`,
		imageId, path,
	)
	if err != nil {
		return "", fmt.Errorf("inserting image into database: %w", err)
	}

	return imageId, nil
}

// DeleteMessageImage removes an image record from the database.
func (db *appdbimpl) DeleteMessageImage(imageId string) error {
	_, err := db.c.Exec(`DELETE FROM images WHERE id = ?`, imageId)
	if err != nil {
		return fmt.Errorf("deleting image from database: %w", err)
	}
	return nil
}
