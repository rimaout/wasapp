package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
)

type MessageStatus string
const (
	StatusDelivered MessageStatus = "delivered"
	StatusReceived  MessageStatus = "received"
	StatusRead      MessageStatus = "read"
)

type EmojiReaction struct {
	UserID  string `json:"userId"`
	EmojiID int32  `json:"emojiId"`
}

type ForwardedFromInfo struct {
	ChatID    string `json:"chatId"`
	MessageID string `json:"messageId"`
}

type MessageContent struct {
	Text        *string `json:"text,omitempty"`
	MsgImageId  *string `json:"msgImageId,omitempty"`
}

type Message struct {
	ID            string             `json:"id"`
	ChatID        string             `json:"chatId"`
	SendTime      time.Time          `json:"sendTime"`
	Sender        User               `json:"sender"`
	Status        MessageStatus      `json:"status"`
	IsDeleted     bool               `json:"isDeleted"`
	IsInitMessage bool               `json:"isInitMessage"`
	ReactionsList []EmojiReaction    `json:"reactionsList"`
	ForwardedFrom *ForwardedFromInfo `json:"forwardedFrom,omitempty"`
	ReplyTo       *string            `json:"replyTo,omitempty"`
	Content       *MessageContent    `json:"content,omitempty"`
}

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

// SaveMessageImagePath inserts an image path into the images table and returns its generated UUID.
func (db *appdbimpl) SaveMessageImagePath(path string) (string, error) {
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

// DeleteMessageImagePath removes an image record from the database.
func (db *appdbimpl) DeleteMessageImagePath(imageId string) error {
	_, err := db.c.Exec(`DELETE FROM images WHERE id = ?`, imageId)
	if err != nil {
		return fmt.Errorf("deleting image from database: %w", err)
	}
	return nil
}

var ErrMessageNotFound = errors.New("message not found")

// GetMessageById retrieves a message with its sender and reactions.
// Returns ErrMessageNotFound if the message does not exist in the specified chat.
func (db *appdbimpl) GetMessageById(chatId, messageId string) (Message, error) {

	// Query the message and its sender from the database
	query := `SELECT m.id, m.chat_id, m.send_time, m.is_deleted, m.is_init_message,
		m.text, m.image_id, m.reply_to_msg_id,
		m.is_forward_message, m.forwarded_from_chat_id, m.forwarded_from_msg_id,
		u.id, u.name
	FROM messages m
	JOIN users u ON m.sender_id = u.id
	WHERE m.id = ? AND m.chat_id = ?`

	// Save the results into variables
	var (
		id, targetChatId, sendTimeStr, senderId, senderName string
		isDeleted, isInit, isForward                        bool
		dbText, dbImageId, dbReplyTo                        *string
		dbFwdChatId, dbFwdMsgId                             *string
	)
	err := db.c.QueryRow(query, messageId, chatId).Scan(
		&id, &targetChatId, &sendTimeStr, &isDeleted, &isInit,
		&dbText, &dbImageId, &dbReplyTo,
		&isForward, &dbFwdChatId, &dbFwdMsgId,
		&senderId, &senderName,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return Message{}, ErrMessageNotFound
	}
	if err != nil {
		return Message{}, fmt.Errorf("querying message by id: %w", err)
	}

	sendTime, err := time.Parse(time.RFC3339, sendTimeStr)
	if err != nil {
		return Message{}, fmt.Errorf("parsing message send_time: %w", err)
	}

	// Construct the Message object
	msg := Message{
		ID:            id,
		ChatID:        targetChatId,
		SendTime:      sendTime,
		Sender:        User{Id: senderId, Name: senderName},
		Status:        StatusDelivered, //TODO: Update status based on delivery/receipt/read events
		IsDeleted:     isDeleted,
		IsInitMessage: isInit,
	}

	// Only include content if the message is not deleted and not an init message
	if !isDeleted && !isInit {
		content := MessageContent{}
		hasContent := false

		if dbText != nil {
			content.Text = dbText
			hasContent = true
		}
		if dbImageId != nil {
			content.MsgImageId = dbImageId
			hasContent = true
		}

		if hasContent {
			msg.Content = &content
		}
	}

	// Include forwarded message info if present
	if isForward && (dbFwdChatId != nil) && (dbFwdMsgId != nil) {
		msg.ForwardedFrom = &ForwardedFromInfo{
			ChatID:    *dbFwdChatId,
			MessageID: *dbFwdMsgId,
		}
	}

	// Include reply-to message ID if present
	if dbReplyTo != nil {
		msg.ReplyTo = dbReplyTo
	}

	// Get and include reactions for the message
	reactions, err := db.getReactionsForMessage(messageId)
	if err != nil {
		return Message{}, fmt.Errorf("querying reactions: %w", err)
	}
	msg.ReactionsList = reactions

	return msg, nil
}

//TODO: In the future also delete the message text and image from the database, but for now we just mark it as deleted,
//		since for the user is still impossible to see the message content, because GetMessageById will return a Message with Content = nil if IsDeleted = true.
func (db *appdbimpl) SetMessageAsDeleted(chatId, messageId string) (Message, error) {

	// Mark the message as deleted in the database
	_, err := db.c.Exec(
		`UPDATE messages SET is_deleted = 1 WHERE id = ? AND chat_id = ?`,
		messageId, chatId,
	)
	if err != nil {
		return Message{}, fmt.Errorf("updating message as deleted: %w", err)
	}

	// Fetch the updated message to return
	msg, err := db.GetMessageById(chatId, messageId)
	if err != nil {
		return Message{}, fmt.Errorf("fetching updated message: %w", err)
	}

	return msg, nil
}
