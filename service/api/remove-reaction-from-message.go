package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/database"
	"github.com/rimaout/wasapp/service/api/reqcontext"
)

func (rt *_router) removeReactionFromMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Extract target chat id and message id from url (chats/{chatId}/messages/{messageId}/reaction)
	chatId := ps.ByName("chatId")
	messageId := ps.ByName("messageId")

	// Validate that the authenticated user has access to the chat (404, 403, 500 errors)
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

	// Delete the reaction associated with the message for the authenticated user
	err = rt.db.DeleteReaction(messageId, ctx.UserID)
	if err == database.ErrReactionNotFound {
		rt.respondWithError(w, http.StatusNotFound, ErrCodeReactionNotFound, "reaction not found") //404
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error("error deleting reaction")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Fetch the updated message after removing the reaction
	msg, err := rt.db.GetMessageById(chatId, messageId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error fetching updated message")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Respond with the updated message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(msg)
}
