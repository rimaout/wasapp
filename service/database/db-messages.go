package database

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/gofrs/uuid"
)

type MessageStatus string

const (
	StatusDelivered MessageStatus = "DELIVERED"
	StatusReceived  MessageStatus = "RECEIVED"
	StatusRead      MessageStatus = "READ"
)

type EmojiReaction struct {
	UserID  string `json:"userId"`
	EmojiID int32  `json:"emojiId"`
}

type ForwardedFromInfo struct {
	ChatID      string `json:"chatId"`
	MessageID   string `json:"messageId"`
	Text       *string `json:"text,omitempty"`
	MsgImageId *string `json:"msgImageId,omitempty"`
}

type RepliedToInfo struct {
	MessageID   string `json:"messageId"`
	SenderID    string `json:"senderId"`
	SenderName  string `json:"senderName"`
	Text       *string `json:"text,omitempty"`
	MsgImageId *string `json:"msgImageId,omitempty"`
}

type MessageContent struct {
	Text       *string `json:"text,omitempty"`
	MsgImageId *string `json:"msgImageId,omitempty"`
}

type Message struct {
	ID             string             `json:"id"`
	ChatID         string             `json:"chatId"`
	SendTime       time.Time          `json:"sendTime"`
	Sender         User               `json:"sender"`
	Status         MessageStatus      `json:"status"`
	IsDeleted      bool               `json:"isDeleted"`
	IsInitMessage  bool               `json:"isInitMessage"`
	IsJoinMessage  bool               `json:"isJoinMessage"`
	IsLeaveMessage bool               `json:"isLeaveMessage"`
	ReactionsList  []EmojiReaction    `json:"reactionsList"`
	ForwardedFrom  *ForwardedFromInfo `json:"forwardedFrom,omitempty"`
	ReplyTo        *string            `json:"replyTo,omitempty"`
	RepliedTo      *RepliedToInfo     `json:"repliedTo,omitempty"`
	Content        *MessageContent    `json:"content,omitempty"`
}

// CreateMessage inserts a new message and returns the message ID.
// Pass empty strings for optional fields that should be NULL.
// Optional fields: inText, inImageId, inReplyTo, inForwardFromChat, inForwardFromMsg.
func (db *appdbimpl) CreateMessage(
	chatId, senderId, sendTime, inText, inImageId, inReplyTo string,
	isInit, isForward bool,
	inForwardFromChat, inForwardFromMsg string,
) (string, error) {
	// Generate a new UUID for the message
	id, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("error generating message UUID: %w", err)
	}
	msgId := id.String()

	// Prepare optional fields for database insertion
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

	// Insert the new message into the database
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

	// Return the generated message ID
	return msgId, nil
}

// CreateSystemMessage inserts a join or leave system message (no content) and returns its ID.
// `isJoin` selects between a join message (true) and a leave message (false).
func (db *appdbimpl) CreateSystemMessage(chatId, senderId string, isJoin bool) (string, error) {
	// Generate a new UUID for the message
	id, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("error generating message UUID: %w", err)
	}
	msgId := id.String()

	// Get the current time
	sendTime := time.Now().UTC().Format(time.RFC3339)

	// Insert the system message into the database
	_, err = db.c.Exec(
		`INSERT INTO messages
			(id, chat_id, sender_id, send_time, is_deleted, is_init_message,
			 is_join_message, is_leave_message,
			 text, image_id, reply_to_msg_id,
			 is_forward_message, forwarded_from_chat_id, forwarded_from_msg_id)
		 VALUES (?, ?, ?, ?, 0, 0, ?, ?, NULL, NULL, NULL, 0, NULL, NULL)`,
		msgId, chatId, senderId, sendTime, isJoin, !isJoin,
	)
	if err != nil {
		return "", fmt.Errorf("inserting system message into database: %w", err)
	}

	// Return the generated message ID
	return msgId, nil
}

// SaveMessageImagePath inserts an image path into the images table and returns its generated UUID.
func (db *appdbimpl) SaveMessageImagePath(path string) (string, error) {
	// Generate a new UUID for the image
	id, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("error generating image UUID: %w", err)
	}
	imageId := id.String()

	// Insert the image record into the database
	_, err = db.c.Exec(
		`INSERT INTO images (id, path) VALUES (?, ?)`,
		imageId, path,
	)
	if err != nil {
		return "", fmt.Errorf("inserting image into database: %w", err)
	}

	// Return the generated image ID
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

// resolveRepliedTo returns info about the message being replied to.
// If the original message is deleted, it returns info with only the
// messageId set (sender/content left empty). If the original message is a forwarded message,
// it quotes the original forwarded-from message's content instead.
func (db *appdbimpl) resolveRepliedTo(replyToId string) (*RepliedToInfo, error) {
	// Initialize the RepliedToInfo struct with the message ID,
	// other fields will be filled if the message exists and is not deleted.
	info := &RepliedToInfo{MessageID: replyToId}

	var senderId, senderName string
	var text, imageId *string
	var isForward bool
	var fwdMsgId *string
	err := db.c.QueryRow(
		`SELECT m.sender_id, u.name, m.text, m.image_id, m.is_forward_message, m.forwarded_from_msg_id
		 FROM messages m JOIN users u ON u.id = m.sender_id
		 WHERE m.id = ? AND m.is_deleted = 0`,
		replyToId,
	).Scan(&senderId, &senderName, &text, &imageId, &isForward, &fwdMsgId)

	// If the message does not exist or is deleted, return info with only the messageId set.
	if errors.Is(err, sql.ErrNoRows) {
		return info, nil
	}

	// If there was an error other than no rows, return the error.
	if err != nil {
		return nil, fmt.Errorf("resolving replied-to content: %w", err)
	}

	// Fill in the sender info
	info.SenderID = senderId
	info.SenderName = senderName

	if isForward && fwdMsgId != nil {
		// Parent is a forwarded message: quote the original message's content.
		var fwdText, fwdImageId *string
		_ = db.c.QueryRow(
			"SELECT text, image_id FROM messages WHERE id = ? AND is_deleted = 0",
			*fwdMsgId,
		).Scan(&fwdText, &fwdImageId)
		info.Text = fwdText
		info.MsgImageId = fwdImageId
	} else {
		// Parent is a normal message: quote its content.
		info.Text = text
		info.MsgImageId = imageId
	}
	return info, nil
}

// GetMessageById retrieves a message with its sender and reactions.
// Returns ErrMessageNotFound if the message does not exist in the specified chat.
func (db *appdbimpl) GetMessageById(chatId, messageId string) (Message, error) {

	// Query the message and its sender from the database
	query := `SELECT m.id, m.chat_id, m.send_time, m.is_deleted, m.is_init_message,
		m.is_join_message, m.is_leave_message,
		m.text, m.image_id, m.reply_to_msg_id,
		m.is_forward_message, m.forwarded_from_chat_id, m.forwarded_from_msg_id,
		u.id, u.name
	FROM messages m
	JOIN users u ON m.sender_id = u.id
	WHERE m.id = ? AND m.chat_id = ?`

	// Save the results into variables
	var (
		id, targetChatId, sendTimeStr, senderId, senderName string
		isDeleted, isInit, isForward, isJoin, isLeave       bool
		dbText, dbImageId, dbReplyTo                        *string
		dbFwdChatId, dbFwdMsgId                             *string
	)
	err := db.c.QueryRow(query, messageId, chatId).Scan(
		&id, &targetChatId, &sendTimeStr, &isDeleted, &isInit,
		&isJoin, &isLeave,
		&dbText, &dbImageId, &dbReplyTo,
		&isForward, &dbFwdChatId, &dbFwdMsgId,
		&senderId, &senderName,
	)

	// Handle the case where the message is not found by returning the specific custom ErrMessageNotFound error
	if errors.Is(err, sql.ErrNoRows) {
		return Message{}, ErrMessageNotFound
	}

	// Handle any other errors that occurred during the query
	if err != nil {
		return Message{}, fmt.Errorf("querying message by id: %w", err)
	}

	// Parse the send_time string into a time.Time object
	sendTime, err := time.Parse(time.RFC3339, sendTimeStr)
	if err != nil {
		return Message{}, fmt.Errorf("parsing message send_time: %w", err)
	}

	// Construct the Message object
	msg := Message{
		ID:             id,
		ChatID:         targetChatId,
		SendTime:       sendTime,
		Sender:         User{Id: senderId, Name: senderName},
		IsDeleted:      isDeleted,
		IsInitMessage:  isInit,
		IsJoinMessage:  isJoin,
		IsLeaveMessage: isLeave,
	}

	// Only include content if the message is not deleted and not a system message (init/join/leave)
	if !isDeleted && !isInit && !isJoin && !isLeave {
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

		// Resolve the original message's content for display
		var fwdText, fwdImageId *string
		err := db.c.QueryRow(
			"SELECT text, image_id FROM messages WHERE id = ? AND is_deleted = 0",
			*dbFwdMsgId,
		).Scan(&fwdText, &fwdImageId)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return Message{}, fmt.Errorf("resolving forwarded content: %w", err)
		}
		if err == nil {
			msg.ForwardedFrom.Text = fwdText
			msg.ForwardedFrom.MsgImageId = fwdImageId
		}
	}

	// Include reply-to message ID if present
	if dbReplyTo != nil {
		msg.ReplyTo = dbReplyTo

		// Resolve the replied-to message's sender and content for display
		msg.RepliedTo, err = db.resolveRepliedTo(*dbReplyTo)
		if err != nil {
			return Message{}, err
		}
	}

	// Get and include reactions for the message (deleted messages report an empty list)
	if msg.IsDeleted {
		msg.ReactionsList = []EmojiReaction{}
	} else {
		reactions, err := db.getReactionsForMessage(messageId)
		if err != nil {
			return Message{}, fmt.Errorf("querying reactions: %w", err)
		}
		msg.ReactionsList = reactions
	}

	// Get message statuts (delivered, received, read) based on receiver_statuses
	status, err := db.ComputeMessageStatus(messageId)
	if err != nil {
		return Message{}, fmt.Errorf("computing message status: %w", err)
	}
	msg.Status = status

	return msg, nil
}

// SetMessageAsDeleted marks a message as deleted, clears its content, and removes any associated image file and record.
func (db *appdbimpl) SetMessageAsDeleted(chatId, messageId string) (Message, error) {

	// Get image_id before clearing it
	var imageId *string
	err := db.c.QueryRow(
		"SELECT image_id FROM messages WHERE id = ? AND chat_id = ?",
		messageId, chatId,
	).Scan(&imageId)

	// Handle errors, but ignore the case where no rows are found (message without an image are valid)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return Message{}, fmt.Errorf("error looking up message image before deletion: %w", err)
	}

	// Mark the message as deleted in the database.
	_, err = db.c.Exec(
		`UPDATE messages SET
			is_deleted = 1,
			text = NULL,
			image_id = NULL,
			reply_to_msg_id = NULL,
			is_forward_message = 0,
			forwarded_from_chat_id = NULL,
			forwarded_from_msg_id = NULL,
			is_init_message = 0
		WHERE id = ? AND chat_id = ?`,
		messageId, chatId,
	)
	if err != nil {
		return Message{}, fmt.Errorf("updating message as deleted: %w", err)
	}

	// Delete image file and record if exists
	if imageId != nil && *imageId != "" {

		// Get image file path from db
		var imagePath string
		err := db.c.QueryRow(
			"SELECT path FROM images WHERE id = ?", *imageId,
		).Scan(&imagePath)

		// If the image path was found, delete the file and the record from the database
		if err == nil {
			_ = os.Remove(imagePath)
		}
		_, _ = db.c.Exec("DELETE FROM images WHERE id = ?", *imageId)

		// Note: in future would make sense to handle the error cases here (TODO)
	}

	// Fetch the updated message to return
	msg, err := db.GetMessageById(chatId, messageId)
	if err != nil {
		return Message{}, fmt.Errorf("fetching updated message: %w", err)
	}

	return msg, nil
}

// GetImagePathByImageId returns the file path for an image by its ID.
// Returns ("", nil) if the image is not found.
func (db *appdbimpl) GetImagePathByImageId(imageId string) (string, error) {
	var path string
	err := db.c.QueryRow(
		"SELECT path FROM images WHERE id = ?", imageId,
	).Scan(&path)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("querying image path: %w", err)
	}
	return path, nil
}

// IsImageInChat checks if an image is associated with any message in a chat.
func (db *appdbimpl) IsImageInChat(chatId, imageId string) (bool, error) {
	var exists bool
	err := db.c.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM messages WHERE chat_id = ? AND image_id = ? LIMIT 1)",
		chatId, imageId,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking image in chat: %w", err)
	}
	return exists, nil
}

// GetChatMessages retrieves the last 50 messages for a chat, ordered by send_time descending.
func (db *appdbimpl) GetChatMessages(chatId string) ([]Message, error) {
	rows, err := db.c.Query(
		`SELECT m.id, m.chat_id, m.send_time, m.is_deleted, m.is_init_message,
			m.is_join_message, m.is_leave_message,
			m.text, m.image_id, m.reply_to_msg_id,
			m.is_forward_message, m.forwarded_from_chat_id, m.forwarded_from_msg_id,
			u.id, u.name
		FROM messages m
		JOIN users u ON m.sender_id = u.id
		WHERE m.chat_id = ?

		-- send_time has only second precision, so we use rowid as tiebreaker
		ORDER BY m.send_time DESC, m.rowid DESC

		-- Limit to the last 50 messages (in future, we may want to implement pagination or infinite scrolling)
		LIMIT 50`,
		chatId,
	)
	if err != nil {
		return nil, fmt.Errorf("querying chat messages: %w", err)
	}
	defer rows.Close()

	// Extract messages list form the query results
	var messages []Message
	for rows.Next() {
		var (
			id, targetChatId, sendTimeStr, senderId, senderName string
			isDeleted, isInit, isForward, isJoin, isLeave       bool
			dbText, dbImageId, dbReplyTo                        *string
			dbFwdChatId, dbFwdMsgId                             *string
		)
		err := rows.Scan(
			&id, &targetChatId, &sendTimeStr, &isDeleted, &isInit,
			&isJoin, &isLeave,
			&dbText, &dbImageId, &dbReplyTo,
			&isForward, &dbFwdChatId, &dbFwdMsgId,
			&senderId, &senderName,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning message row: %w", err)
		}

		// Parse the send_time string into a time.Time object
		sendTime, err := time.Parse(time.RFC3339, sendTimeStr)
		if err != nil {
			return nil, fmt.Errorf("parsing message send_time: %w", err)
		}

		// Initialize the Message struct with basic info
		msg := Message{
			ID:             id,
			ChatID:         targetChatId,
			SendTime:       sendTime,
			Sender:         User{Id: senderId, Name: senderName},
			IsDeleted:      isDeleted,
			IsInitMessage:  isInit,
			IsJoinMessage:  isJoin,
			IsLeaveMessage: isLeave,
		}

		// Only include content if the message is not deleted and not a system message (init/join/leave)
		if !isDeleted && !isInit && !isJoin && !isLeave {
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
		if isForward && dbFwdChatId != nil && dbFwdMsgId != nil {
			msg.ForwardedFrom = &ForwardedFromInfo{
				ChatID:    *dbFwdChatId,
				MessageID: *dbFwdMsgId,
			}

			// Query the original message's content
			var fwdText, fwdImageId *string
			_ = db.c.QueryRow(
				"SELECT text, image_id FROM messages WHERE id = ? AND is_deleted = 0",
				*dbFwdMsgId,
			).Scan(&fwdText, &fwdImageId)
			msg.ForwardedFrom.Text = fwdText
			msg.ForwardedFrom.MsgImageId = fwdImageId
		}

		// Include reply to message ID if present
		if dbReplyTo != nil {
			msg.ReplyTo = dbReplyTo

			// Resolve the replied-to message's sender and content for display
			msg.RepliedTo, _ = db.resolveRepliedTo(*dbReplyTo)
		}

		// Get and include reactions for the message (deleted messages report an empty list)
		if isDeleted {
			msg.ReactionsList = []EmojiReaction{}
		} else {
			reactions, err := db.getReactionsForMessage(id)
			if err != nil {
				return nil, fmt.Errorf("querying reactions for %s: %w", id, err)
			}
			msg.ReactionsList = reactions
		}

		// Get message statuts (delivered, received, read) based on receiver_statuses
		status, err := db.ComputeMessageStatus(id)
		if err != nil {
			return nil, fmt.Errorf("computing message status for %s: %w", id, err)
		}
		msg.Status = status

		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating message rows: %w", err)
	}

	// Return an empty list instead of nil if there are no messages
	if messages == nil {
		messages = make([]Message, 0)
	}

	return messages, nil
}

// InsertImageVisibility records that an image is visible in a specific chat
func (db *appdbimpl) InsertImageVisibility(imageId, chatId string) error {
	_, err := db.c.Exec(
		"INSERT OR IGNORE INTO image_chat_visibility (image_id, chat_id) VALUES (?, ?)",
		imageId, chatId,
	)
	if err != nil {
		return fmt.Errorf("inserting image visibility: %w", err)
	}
	return nil
}

// IsImageVisibleInChat checks if an image is marked as visible in a specific chat
func (db *appdbimpl) IsImageVisibleInChat(chatId, imageId string) (bool, error) {
	var exists bool
	err := db.c.QueryRow(
		"SELECT EXISTS(SELECT 1 FROM image_chat_visibility WHERE chat_id = ? AND image_id = ? LIMIT 1)",
		chatId, imageId,
	).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("checking image visibility: %w", err)
	}
	return exists, nil
}
