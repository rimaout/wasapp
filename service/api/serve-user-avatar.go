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
)

// Embed the default avatar image into the binary, in the defaultUserAvatar variable.
// This allows us to serve a default image without needing to read from disk.
//
//go:embed assets/default-user-avatar.png
var defaultUserAvatar []byte

func (rt *_router) getUserAvatar(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := ps.ByName("userId")

	// Check if the user exists.
	userName, err := rt.db.GetUserNameById(userId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking if user exists")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}
	if userName == "" {
		rt.respondWithError(w, http.StatusNotFound, ErrCodeUserNotFound, "user not found") //404
		return
	}

	// Get the avatar image path for the user.
	imagePath, err := rt.db.GetUserAvatarPath(userId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting avatar image path")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// If the user doesn't have an avatar, serve the default avatar.
	if imagePath == "" {
		w.Header().Set("Content-Type", "image/png")

		_, err = w.Write(defaultUserAvatar)
		if err != nil {
			ctx.Logger.WithError(err).Error("error writing default avatar to response")
		}
		return
	}

	http.ServeFile(w, r, imagePath)
}
