package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/api/reqcontext"
)

func (rt *_router) leaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Extract target chat id from url (chats/{chatId}/leave)
	chatId := ps.ByName("chatId")

	// Validate that the authenticated user has access to the chat (404, 403, 500 errors)
	if !rt.validateChatAccess(w, r, ctx, chatId) {
		// Checks errors: 404 (chat not found), 403 (user not a member), 500 (internal error)
		return
	}

	// Check if chat is groupchat (409)
	isGroup, err := rt.db.IsGroupChat(chatId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking if chat is a group chat")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") // 500
		return
	}
	if !isGroup {
		rt.respondWithError(w, http.StatusConflict, ErrCodeNotAGroupChat, "chat is not a group chat, only members of group chats can leave the group") // 409
		return
	}

	// Set leave time
	err = rt.db.SetLeaveTime(chatId, ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error setting leave time for user in chat")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") // 500
		return
	}

	// Create a leave system message (sender is the leaving user)
	leaveMsgId, err := rt.db.CreateSystemMessage(chatId, ctx.UserID, false)
	if err != nil {
		ctx.Logger.WithError(err).Error("error creating leave message")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") // 500
		return
	}

	// Create the receiver statuses for the leave message (the leaver is excluded,
	// as their group_leave_time is already set)
	_ = rt.db.InsertReceiverStatuses(leaveMsgId, chatId, ctx.UserID)

	// Return 204 No Content on success
	w.WriteHeader(http.StatusNoContent)
}
