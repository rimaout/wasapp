package api

import (
	"net/http"
	"encoding/json"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/database"
	"github.com/rimaout/wasapp/service/api/reqcontext"
)


func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Extract target chat id from url
	chatId    := ps.ByName("chatId")
	messageId := ps.ByName("messageId")

	if !rt.validateChatAccess(w, r, ctx, chatId) {
		// Checks errors: 404 (chat not found), 403 (user not a member), 500 (internal error)
		return
	}

	// Check if message exists (404)
	msg, err := rt.db.GetMessageById(chatId, messageId)
	if err == database.ErrMessageNotFound {
		rt.respondWithError(w, http.StatusNotFound, ErrCodeMessageNotFound, "message not found")
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting message by id")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error")
		return
	}

	// Check if message is owned by the logged user (403)
	if msg.Sender.Id != ctx.UserID {
		rt.respondWithError(w, http.StatusForbidden, ErrCodeForbiddenNotSender, "user is not the owner of the message")
		return
	}

	// Check if message is already deleted (400)
	if msg.IsDeleted {
		rt.respondWithError(w, http.StatusBadRequest, ErrCodeAlreadyDeleted, "message already deleted")
		return
	}

	// Check if message is initMessage (400)
	if msg.IsInitMessage {
		rt.respondWithError(w, http.StatusBadRequest, ErrCodeItsInitMessage, "cannot delete init message")
		return
	}

	// Mark the db record as delete message from
	updatedMsg, err := rt.db.SetMessageAsDeleted(chatId, messageId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error marking message as deleted")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error")
		return
	}

	// Return the updated message as JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(updatedMsg)
	if err != nil {
		ctx.Logger.WithError(err).Error("error encoding response JSON")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error")
		return
	}
}
