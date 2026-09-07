package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/rimaout/wasapp/service/api/reqcontext"
)

func (rt *_router) deleteMyUserAvatar(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Get the current avatar image path to delete it later
	oldImagePath, err := rt.db.GetUserAvatarPath(ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error fetching old avatar path")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") // 500
		return
	}

	// Clear the avatar path in the database
	if err := rt.db.SetUserAvatarPath(ctx.UserID, ""); err != nil {
		ctx.Logger.WithError(err).Error("error clearing avatar image path in database")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") // 500
		return
	}

	// Now it's safe to delete the old image from disk
	rt.deleteImageFile(ctx, oldImagePath)

	// Return 204 No Content
	w.WriteHeader(http.StatusNoContent)
}
