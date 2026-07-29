package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/database"
	"github.com/rimaout/wasapp/service/api/reqcontext"
)

const groupChatAvatarUploadDir = "uploads/users/avatars/"

func (rt *_router) setGroupChatAvatar(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Extract target chat id form url (http params)
	targetChatId := ps.ByName("chatId")

	// Check if chat exist (404 error)
	_, err := rt.db.GetChatById(targetChatId)
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
	isGroup, err := rt.db.IsGroupChat(targetChatId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking if chat is a group chat")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}
	if !isGroup {
		rt.respondWithError(w, http.StatusConflict, ErrCodeNotAGroupChat, "chat is not a group chat, only the names of group chats can be set/changed") //409
		return
	}

	// Check if logged user ia a member of the group chat (403 error)
	isMember, err := rt.db.IsActiveChatMember(ctx.UserID, targetChatId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking if user is a member of the chat")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}
	if !isMember {
		rt.respondWithError(w, http.StatusForbidden, ErrCodeForbiddenNotMember, "you are not a member of this group chat") //403
		return
	}

	// Get the current avatar image path
	oldImagePath, err := rt.db.GetGroupAvatarPath(targetChatId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error fetching old group avatar path")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Extract and save the new image
	newImagePath, ok := rt.saveUploadedImage(w, r, ctx, "binaryImage", groupChatAvatarUploadDir)
	if !ok {
		// TODO:
		//	 - 400 (invalid form / missing file),
		//	 - 413 (too large), 415 (bad format), 500 (disk error).
		return
	}

	// Update the database with the new file path
	if err := rt.db.SetGroupAvatarPath(targetChatId, newImagePath); err != nil {
		ctx.Logger.WithError(err).Error("error updating group avatar image path in database")
		rt.deleteImageFile(ctx, newImagePath)
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Now it's safe to delete the old image from disk
	rt.deleteImageFile(ctx, oldImagePath)

	// Return Response
	w.WriteHeader(http.StatusNoContent)
}
