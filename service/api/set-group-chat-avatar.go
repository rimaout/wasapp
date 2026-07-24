package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/api/reqcontext"
)

const groupChatAvatarUploadDir = "uploads/users/avatars/"

func (rt *_router) setGroupChatAvatar(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	//TODO: Check if chat exist (404 error)

	//TODO: Check if chat is group chat (409 error)

	//TODO: Check if logged user is member (403 error)

	// Get the current avatar image path
	oldImagePath, err := rt.db.GetGroupAvatarPath(ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error fetching old group avatar path")
		rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
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
	if err := rt.db.SetGroupAvatarPath(ctx.UserID, newImagePath); err != nil {
		ctx.Logger.WithError(err).Error("error updating group avatar image path in database")
		rt.deleteImageFile(ctx, newImagePath)
		rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
		return
	}

	// Now it's safe to delete the old image from disk
	rt.deleteImageFile(ctx, oldImagePath)

	// Return Response
	w.WriteHeader(http.StatusNoContent)
}
