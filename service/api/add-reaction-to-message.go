package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/database"
	"github.com/rimaout/wasapp/service/api/reqcontext"
)

type emojiReactionRequest struct {
	EmojiId int32 `json:"emojiId"`
}

func (rt *_router) addReactionToMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	chatId := ps.ByName("chatId")
	messageId := ps.ByName("messageId")

	if !rt.validateChatAccess(w, r, ctx, chatId) {
		// Checks errors: 404 (chat not found), 403 (user not a member), 500 (internal error)
		return
	}

	// Check if the message exists
	_, err := rt.db.GetMessageById(chatId, messageId)
	if err == database.ErrMessageNotFound {
		rt.respondWithError(w, http.StatusNotFound, ErrCodeMessageNotFound, "message not found") //404
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting message by id")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Decode the request body
	var req emojiReactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rt.respondWithError(w, http.StatusBadRequest, ErrCodeInvalidInput, "invalid request body") //400
		return
	}

	// Validate the emoji id
	if req.EmojiId < 0 || req.EmojiId > 9 {
		rt.respondWithError(w, http.StatusBadRequest, ErrCodeInvalidEmojiId, "invalid emoji id") //400
		return
	}

	// Add the reaction to the message
	err = rt.db.CreateReaction(messageId, ctx.UserID, req.EmojiId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error adding reaction to message")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Get updated message with reactions
	msg, err := rt.db.GetMessageById(chatId, messageId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error fetching updated message")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Respond with the updated message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)
}
