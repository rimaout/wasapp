package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/database"
	"github.com/rimaout/wasapp/service/api/reqcontext"
)

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Extract target chat id from url (http params)
	chatId := ps.ByName("chatId")

	// Check if chat exists (404)
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

	// Check if chat is groupchat (409)
	isGroup, err := rt.db.IsGroupChat(chatId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking if chat is a group chat")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}
	if !isGroup {
		rt.respondWithError(w, http.StatusConflict, ErrCodeNotAGroupChat, "chat is not a group chat, only members of group chats can leave the group") //409
		return
	}

	// Check if logged user is group member (403)
	isMember, err := rt.db.IsActiveChatMember(chatId, ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking if user is a member of the chat")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}
	if !isMember {
		rt.respondWithError(w, http.StatusForbidden, ErrCodeForbiddenNotMember, "user is not a member of the group chat") //403
		return
	}

	// Set leave time
	err = rt.db.SetLeaveTime(chatId, ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error setting leave time for user in chat")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Return 204 No Content on success
	w.WriteHeader(http.StatusNoContent)
}
