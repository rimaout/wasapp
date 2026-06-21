package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/api/reqcontext"
)

const profileUploadDir = "uploads/profile"

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the current profile image path
	// DB error → 500 (server can't proceed without knowing if old image exists)
	oldImagePath, err := rt.db.GetUserProfileImagePath(ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error fetching old profile path")
		rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
		return
	}

	// Extract and save the new image
	newImagePath, ok := rt.saveUploadedImage(w, r, ctx, "binaryImage", profileUploadDir)
	if !ok {
		// can emit: 400 (invalid form / missing file),
		// 413 (too large), 415 (bad format), 500 (disk error).
		return
	}

	// Update the database with the new file path
	if err := rt.db.SetUserProfileImage(ctx.UserID, newImagePath); err != nil {
		ctx.Logger.WithError(err).Error("error updating profile image path in database")
		rt.deleteImageFile(ctx, newImagePath)
		rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
		return
	}

	// Now it's safe to delete the old image from disk
	rt.deleteImageFile(ctx, oldImagePath)

	// Return Response
	userName, _ := rt.db.GetUserNameById(ctx.UserID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profileResponse{Id: ctx.UserID, Name: userName})
}
