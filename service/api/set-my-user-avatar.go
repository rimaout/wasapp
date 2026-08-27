package api

import (
	"net/http"
	"path/filepath"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/api/reqcontext"
)

const userAvatarUploadDir = "users/avatars"

func (rt *_router) setMyUserAvatar(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Get the current avatar image path
	oldImagePath, err := rt.db.GetUserAvatarPath(ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error fetching old avatar path")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Extract and save the new image.
	// (saveUploadedImage writes 400, 413, 415, or 500 error responses directly on failure)
	newImagePath, ok := rt.saveUploadedImage(w, r, ctx, "binaryImage", filepath.Join(rt.uploadsDir, userAvatarUploadDir))
	if !ok {
		return
	}

	// Update the database with the new file path
	if err := rt.db.SetUserAvatarPath(ctx.UserID, newImagePath); err != nil {
		ctx.Logger.WithError(err).Error("error updating avatar image path in database")
		rt.deleteImageFile(ctx, newImagePath)
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Now it's safe to delete the old image from disk
	rt.deleteImageFile(ctx, oldImagePath)

	// Return Response
	w.WriteHeader(http.StatusNoContent)
}
