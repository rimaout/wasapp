package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/api/reqcontext"
	"github.com/rimaout/wasapp/service/database"
)

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Extract target chat id from url (chats/{chatId}/messages)
	chatId := ps.ByName("chatId")

	if !rt.validateChatAccess(w, r, ctx, chatId) {
		// Checks errors: 404 (chat not found), 403 (user not a member), 500 (internal error)
		return
	}

	rt.createMessageLogic(w, r, ctx, chatId, "")
	// Note: it also handles errors: 400 (invalid request), 413 (image too large), 415 (unsupported media type), 500 (internal error)
}

func (rt *_router) validateChatAccess(w http.ResponseWriter, r *http.Request, ctx reqcontext.RequestContext, chatId string) bool {
	// Validate that the chat exists
	_, err := rt.db.GetChatById(chatId)
	if err == database.ErrChatNotFound {
		rt.respondWithError(w, http.StatusNotFound, ErrCodeChatNotFound, "chat not found") // 404
		return false
	}
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting chat by id")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") // 500
		return false
	}

	// Check if the authenticated user is a member of the chat
	isMember, err := rt.db.IsActiveChatMember(chatId, ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking if user is a member of the chat")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") // 500
		return false
	}
	if !isMember {
		rt.respondWithError(w, http.StatusForbidden, ErrCodeForbiddenNotMember, "user is not a member of the chat") // 403
		return false
	}

	return true
}

func (rt *_router) createMessageLogic(w http.ResponseWriter, r *http.Request, ctx reqcontext.RequestContext, chatId, replyTo string) {

	// Extract and validate message request content (text, image)
	// Note: if an image was sent, it will be saved to disk and its path will be returned in reqContent.ImagePath
	reqContent, ok := rt.extractSendMessageContent(w, r, ctx)
	if !ok {
		return // Helper already sent 400 (invalid request), 413 (image too large), 415 (unsupported media type), 500 (internal error)
	}

	// If an image was uploaded, save its path to the database (images table)
	var imageId string
	if reqContent.ImagePath != "" {
		var err error
		imageId, err = rt.db.SaveMessageImagePath(reqContent.ImagePath)
		if err != nil {
			ctx.Logger.WithError(err).Error("error saving image path to database")

			// Clean up file on disk since DB insert failed
			rt.deleteImageFile(ctx, reqContent.ImagePath)

			rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") // 500
			return
		}

		// Grant this chat visibility to the uploaded image
		_ = rt.db.InsertImageVisibility(imageId, chatId)
	}

	// Save message to DB
	sendTime := time.Now().UTC().Format(time.RFC3339)
	messageId, err := rt.db.CreateMessage(
		chatId,
		ctx.UserID,
		sendTime,
		reqContent.Text,
		imageId,
		replyTo,
		false,
		false,
		"",
		"",
	)
	if err != nil {
		ctx.Logger.WithError(err).Error("error saving message to database")

		// Clean up uploaded image if database insert fails
		if reqContent.ImagePath != "" {
			_ = rt.db.DeleteMessageImagePath(imageId)     // Delete image record from database
			rt.deleteImageFile(ctx, reqContent.ImagePath) // Delete image file from disk
		}

		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") // 500
		return
	}

	// Fetch the created message to return as response
	msg, err := rt.db.GetMessageById(chatId, messageId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error fetching created message")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") // 500
		return
	}

	// Create the receiver statuses for the message (one for each member of the chat)
	// This means that the message is considered "sent" to all members, but not yet "delivered" or "read"
	_ = rt.db.InsertReceiverStatuses(messageId, chatId, ctx.UserID)

	// Send 201 Created response with the full message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)
}
