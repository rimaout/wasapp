package api

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/rimaout/wasapp/service/api/reqcontext"
	"github.com/rimaout/wasapp/service/database"
)

func (rt *_router) getChatAvatar(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Extract target chat id form url (http params)
	chatId := ps.ByName("chatId")

	// Check if chat exist (404 error)
	_, err := rt.db.GetChatById(chatId)
	if errors.Is(err, database.ErrChatNotFound) {
		rt.respondWithError(w, http.StatusNotFound, ErrCodeChatNotFound, "chat not found") // 404
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting chat by id")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") // 500
		return
	}

	// Check if chat is a group chat (409 error)
	isGroup, err := rt.db.IsGroupChat(chatId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking if chat is a group chat")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") // 500
		return
	}

	// Check if authenticated user is member of the chat (403 error)
	isMember, err := rt.db.IsActiveChatMember(chatId, ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking if user is member of chat")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") // 500
		return
	}
	if !isMember {
		rt.respondWithError(w, http.StatusForbidden, ErrCodeForbiddenNotMember, "user is not a member of the chat") // 403
		return
	}

	// Initialize imagePath variable to hold the path of the avatar imagePath
	imagePath := ""

	if isGroup {
		// If the chat is a group chat, we need to get the avatar image of the group

		// Get the avatar image path for the group
		imagePath, err = rt.db.GetGroupAvatarPath(chatId)
		if err != nil {
			ctx.Logger.WithError(err).Error("error getting avatar image path")
			rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") // 500
			return
		}

		if imagePath == "" {
			rt.respondWithError(w, http.StatusNotFound, ErrCodeNoChatAvatar, "no avatar set for this chat") // 404
			return
		}
	} else {
		// If the chat is not a group chat, we need to get the avatar image of the other user in the chat

		// Get the other user id in the chat
		otherUserId, err := rt.db.GetOtherMemberId(chatId, ctx.UserID)
		if err != nil {
			ctx.Logger.WithError(err).Error("error getting other member id")
			rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") // 500
			return
		}

		// Get the avatar image path for the other user
		imagePath, err = rt.db.GetUserAvatarPath(otherUserId)
		if err != nil {
			ctx.Logger.WithError(err).Error("error getting avatar image path")
			rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") // 500
			return
		}

		if imagePath == "" {
			rt.respondWithError(w, http.StatusNotFound, ErrCodeNoChatAvatar, "no avatar set for this chat") // 404
			return
		}
	}

	http.ServeFile(w, r, imagePath)
}
