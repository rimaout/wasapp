package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/database"
	"github.com/rimaout/wasapp/service/api/reqcontext"
)

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Extract target chat id from url (http params)
	targetChatId := ps.ByName("chatId")

	// Extract user id from context (logged user)
	userId := ctx.UserID

	// Check if chat exists (404)
	_, err := rt.db.GetChatById(targetChatId)
	if err == database.ErrChatNotFound {
		rt.respondWithError(w, http.StatusNotFound, "404", "chat not found")
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting chat by id")
		rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
		return
	}

	// Check if chat is groupchat (409)
	isGroup, err := rt.db.IsGroupChat(targetChatId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking if chat is a group chat")
		rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
		return
	}
	if !isGroup {
		rt.respondWithError(w, http.StatusConflict, "409", "chat is not a group chat, only members of group chats can leave the group")
		return
	}

	// Check if logged user is group member (403)
	isMember, err := rt.db.IsActiveChatMember(userId, targetChatId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking if user is a member of the chat")
		rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
		return
	}
	if !isMember {
		rt.respondWithError(w, http.StatusForbidden, "403", "user is not a member of the group chat")
		return
	}

	// Set leave time
	err = rt.db.SetLeaveTime(targetChatId, userId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error setting leave time for user in chat")
		rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
		return
	}

	// Return 204 No Content on success
	w.WriteHeader(http.StatusNoContent)
}
