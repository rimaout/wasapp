package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/api/reqcontext"
)

const userAvatarUploadDir = "uploads/users/avatars/"

func (rt *_router) setMyUserAvatar(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Get the current avatar image path
	oldImagePath, err := rt.db.GetUserAvatarPath(ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error fetching old avatar path")
		rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
		return
	}

	// Extract and save the new image
	newImagePath, ok := rt.saveUploadedImage(w, r, ctx, "binaryImage", userAvatarUploadDir)
	if !ok {
		// TODO:
		//	 - 400 (invalid form / missing file),
		//	 - 413 (too large), 415 (bad format), 500 (disk error).
		return
	}

	// Update the database with the new file path
	if err := rt.db.SetUserAvatarPath(ctx.UserID, newImagePath); err != nil {
		ctx.Logger.WithError(err).Error("error updating avatar image path in database")
		rt.deleteImageFile(ctx, newImagePath)
		rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
		return
	}

	// Now it's safe to delete the old image from disk
	rt.deleteImageFile(ctx, oldImagePath)

	// Return Response
	userName, _ := rt.db.GetUserNameById(ctx.UserID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userResponse{Id: ctx.UserID, Name: userName})
}
