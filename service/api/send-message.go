package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/database"
	"github.com/rimaout/wasapp/service/api/reqcontext"
)


func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Extract target chat id from url
	chatId := ps.ByName("chatId")

	// Check if chat exists (404)
	_, err := rt.db.GetChatById(chatId)
	if err == database.ErrChatNotFound {
		rt.respondWithError(w, http.StatusNotFound, ErrCodeChatNotFound, "chat not found")
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting chat by id")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error")
		return
	}

	// Check if logged user is a member of the chat (403)
	isMember, err := rt.db.IsActiveChatMember(chatId, ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking if user is a member of the chat")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error")
		return
	}
	if !isMember {
		rt.respondWithError(w, http.StatusForbidden, ErrCodeForbiddenNotMember, "user is not a member of the chat")
		return
	}

	// Extract and validate message request content (text, image)
	// Note: if an image was sent, it will be saved to disk and its path will be returned in reqContent.ImagePath
	reqContent, ok := rt.extractSendMessageContent(w, r, ctx)
	if !ok {
		return // Helper already sent the 400, 413, 415, or 500 response
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

			rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error")
			return
		}
	}

	// Save message to DB
	sendTime := time.Now().UTC().Format(time.RFC3339)
	messageId, err := rt.db.CreateMessage(
		chatId,
		ctx.UserID,
		sendTime,
		reqContent.Text,
		imageId,
		"",
		false,
		false,
		"",
		"",
	)
	if err != nil {
		ctx.Logger.WithError(err).Error("error saving message to database")

		// Clean up uploaded image if database insert fails
		if reqContent.ImagePath != "" {
			_ = rt.db.DeleteMessageImagePath(imageId)          // Delete image record from database
			rt.deleteImageFile(ctx, reqContent.ImagePath)  // Delete image file from disk
		}

		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error")
		return
	}

	// Fetch the created message to return as response
	msg, err := rt.db.GetMessageById(chatId, messageId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error fetching created message")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error")
		return
	}

	// Send 201 Created response with the full message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)
}
