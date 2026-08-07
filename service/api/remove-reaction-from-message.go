package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/database"
	"github.com/rimaout/wasapp/service/api/reqcontext"
)

func (rt *_router) removeReactionFromMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	chatId := ps.ByName("chatId")
	messageId := ps.ByName("messageId")

	// Check if the chat exists
	_, err := rt.db.GetChatById(chatId)
	if err == database.ErrChatNotFound {
		rt.respondWithError(w, http.StatusNotFound, ErrCodeChatNotFound, "chat not found") //404
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting chat by id")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Check if logged user is a member of the chat
	isMember, err := rt.db.IsActiveChatMember(chatId, ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking if user is member of chat")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}
	if !isMember {
		rt.respondWithError(w, http.StatusForbidden, ErrCodeForbiddenNotMember, "user is not a member of the chat") //403
		return
	}

	// Check if the message exists
	_, err = rt.db.GetMessageById(chatId, messageId)
	if err == database.ErrMessageNotFound {
		rt.respondWithError(w, http.StatusNotFound, ErrCodeMessageNotFound, "message not found") //404
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting message by id")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Delete the reaction associated with the message for the logged user
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
