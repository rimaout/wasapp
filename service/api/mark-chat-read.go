package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/api/reqcontext"
)

func (rt *_router) markChatRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	chatId := ps.ByName("chatId")

	if !rt.newMessageRequestValidation(w, r, ctx, chatId) {
		// Checks errors: 404 (chat not found), 403 (user not a member), 500 (internal error)
		return
	}

	_ = rt.db.MarkMessagesReadByUser(chatId, ctx.UserID)

	w.WriteHeader(http.StatusNoContent)
}
