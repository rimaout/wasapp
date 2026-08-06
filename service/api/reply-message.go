package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/api/reqcontext"
	"github.com/rimaout/wasapp/service/database"
)

func (rt *_router) replyMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	chatId := ps.ByName("chatId")
	messageId := ps.ByName("messageId")

	if !rt.newMessageRequestValidation(w, r, ctx, chatId) {
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

	rt.createMessageLogic(w, r, ctx, chatId, messageId)
	// Note: it also handles errors: 400 (invalid request), 413 (image too large), 415 (unsupported media type), 500 (internal error)

}
