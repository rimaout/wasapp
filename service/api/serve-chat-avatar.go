package api

import (
	// We use the "_" to silence the "unused import" warning.
	// The embed package is required to enable the //go:embed compiler directive,
	// but because directives look like standard comments, the Go import checker
	// doesn't recognize it as a "real" usage.
	_ "embed"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/rimaout/wasapp/service/api/reqcontext"
	"github.com/rimaout/wasapp/service/database"
)

// Embed the default avatar image into the binary, in the defaultGroupAvatar variable.
// This allows us to serve a default image without needing to read from disk.
//
//go:embed assets/default-group-avatar.png
var defaultGroupAvatar []byte

func (rt *_router) getChatAvatar(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Extract target chat id form url (http params)
	chatId := ps.ByName("chatId")

	// Check if chat exist (404 error)
	_, err := rt.db.GetChatById(chatId)
	if err == database.ErrChatNotFound {
		rt.respondWithError(w, http.StatusNotFound, "404", "chat not found")
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting chat by id")
		rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
		return
	}

	// Check if chat is a group chat (409 error)
	isGroup, err := rt.db.IsGroupChat(chatId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking if chat is a group chat")
		rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
		return
	}

	// Check if logged user is member of the chat (403 error)
	isMember, err := rt.db.IsActiveChatMember(chatId, ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking if user is member of chat")
		rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
		return
	}
	if !isMember {
		rt.respondWithError(w, http.StatusForbidden, "403", "user is not a member of the chat")
		return
	}

	// Initialize imagePath variable to hold the path of the avatar imagePath
	imagePath := ""

	if isGroup {
		// If the chat is a group chat, we need to get the avatar image of the group

		// Get the avatar image path for the group
		imagePath, err := rt.db.GetGroupAvatarPath(chatId)
		if err != nil {
			ctx.Logger.WithError(err).Error("error getting avatar image path")
			rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
			return
		}

		// If the group doesn't have an avatar, serve the default avatar.
		if imagePath == "" {
			w.Header().Set("Content-Type", "image/png")
			_, err = w.Write(defaultGroupAvatar)
			if err != nil {
				ctx.Logger.WithError(err).Error("error writing default avatar to response")
			}
			return
		}
	} else {
		// If the chat is not a group chat, we need to get the avatar image of the other user in the chat

		// Get the other user id in the chat
		otherUserId, err := rt.db.GetOtherMemberId(chatId, ctx.UserID)
		if err != nil {
			ctx.Logger.WithError(err).Error("error getting other member id")
			rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
			return
		}

		// Get the avatar image path for the other user
		imagePath, err = rt.db.GetUserAvatarPath(otherUserId)
		if err != nil {
			ctx.Logger.WithError(err).Error("error getting avatar image path")
			rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
			return
		}

		// If the user doesn't have an avatar, serve the default user avatar.
		if imagePath == "" {
			w.Header().Set("Content-Type", "image/png")
			_, err = w.Write(defaultGroupAvatar)
			if err != nil {
				ctx.Logger.WithError(err).Error("error writing default avatar to response")
			}
			return
		}
	}

	http.ServeFile(w, r, imagePath)
}
