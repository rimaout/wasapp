package api

import (
	"net/http"
	"encoding/json"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/database"
	"github.com/rimaout/wasapp/service/api/reqcontext"
)

func (rt *_router) getChatMembers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Extract target chat id form url (http params)
	chatId := ps.ByName("chatId")

	// Check if chat exist (404 error)
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

	// Check if logged user is a member of the chat (403 error)
	isMember, err := rt.db.IsActiveChatMember(chatId, ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking if user is a member of the chat")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}
	if !isMember {
		rt.respondWithError(w, http.StatusForbidden, ErrCodeForbiddenNotMember, "user is not a member of the chat") //403
		return
	}

	// Get members of the chat
	members, err := rt.db.GetChatMembers(chatId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting chat members")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Respond with members list (200 OK)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(members)
}
