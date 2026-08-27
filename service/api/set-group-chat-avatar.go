package api

import (
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/database"
	"github.com/rimaout/wasapp/service/api/reqcontext"
)

const groupChatAvatarUploadDir = "groups/avatars"

func (rt *_router) setGroupChatAvatar(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Extract target chat id from url (http params)
	chatId := ps.ByName("chatId")

	// Check if chat exists (404 error)
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

	// Check if chat is a group chat (409 error)
	isGroup, err := rt.db.IsGroupChat(chatId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking if chat is a group chat")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}
	if !isGroup {
		rt.respondWithError(w, http.StatusConflict, ErrCodeNotAGroupChat, "chat is not a group chat, only group avatars can be updated") // 409
		return
	}

	// Check if logged user is a member of the group chat (403 error)
	isMember, err := rt.db.IsActiveChatMember(chatId, ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking if user is a member of the chat")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}
	if !isMember {
		rt.respondWithError(w, http.StatusForbidden, ErrCodeForbiddenNotMember, "you are not a member of this group chat") //403
		return
	}

	// Get the current group avatar image path to delete later
	oldImagePath, err := rt.db.GetGroupAvatarPath(chatId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error fetching old group avatar path")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Extract and save the new image.
	// (saveUploadedImage writes 400, 413, 415, or 500 error responses directly on failure)
	newImagePath, ok := rt.saveUploadedImage(w, r, ctx, "binaryImage", filepath.Join(rt.uploadsDir, groupChatAvatarUploadDir))
	if !ok {
		return
	}

	// Update the database with the new file path
	if err := rt.db.SetGroupAvatarPath(chatId, newImagePath); err != nil {
		ctx.Logger.WithError(err).Error("error updating group avatar image path in database")
		rt.deleteImageFile(ctx, newImagePath)
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Now it's safe to delete the old image from disk
	rt.deleteImageFile(ctx, oldImagePath)

	// Return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
