package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/api/reqcontext"
	"github.com/rimaout/wasapp/service/database"
)

func (rt *_router) replyMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Extract target chat id and parent message id from url (chats/{chatId}/messages/{messageId}/reply)
	chatId := ps.ByName("chatId")
	messageId := ps.ByName("messageId")

	// Validate that the authenticated user has access to the chat (404, 403, 500 errors)
	if !rt.validateChatAccess(w, r, ctx, chatId) {
		// Checks errors: 404 (chat not found), 403 (user not a member), 500 (internal error)
		return
	}

	// Validate the parent message exists in the chat (404 error)
	_, err := rt.db.GetMessageById(chatId, messageId)
	if err == database.ErrMessageNotFound {
		rt.respondWithError(w, http.StatusNotFound, ErrCodeMessageNotFound, "message not found") //404
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error("error fetching parent message")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Call the common message creation logic for replying to a message
	rt.createMessageLogic(w, r, ctx, chatId, messageId)
	// Note: it also handles errors: 400 (invalid request), 413 (image too large), 415 (unsupported media type), 500 (internal error)
}
