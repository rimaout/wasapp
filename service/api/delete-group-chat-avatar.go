package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/rimaout/wasapp/service/api/reqcontext"
)

func (rt *_router) deleteGroupChatAvatar(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Extract target chat id from url (http params)
	chatId := ps.ByName("chatId")

	// Check chat exists and membership (404/403/500 errors)
	if !rt.validateChatAccess(w, r, ctx, chatId) {
		return
	}

	// Check if chat is a group chat (409 error)
	isGroup, err := rt.db.IsGroupChat(chatId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking if chat is a group chat")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}
	if !isGroup {
		rt.respondWithError(w, http.StatusConflict, ErrCodeNotAGroupChat, "chat is not a group chat, only group avatars can be removed") //409
		return
	}

	// Get the current group avatar image path to delete later
	oldImagePath, err := rt.db.GetGroupAvatarPath(chatId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error fetching old group avatar path")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Clear the avatar path in the database
	if err := rt.db.SetGroupAvatarPath(chatId, ""); err != nil {
		ctx.Logger.WithError(err).Error("error clearing group avatar image path in database")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Now it's safe to delete the old image from disk
	rt.deleteImageFile(ctx, oldImagePath)

	// Return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
