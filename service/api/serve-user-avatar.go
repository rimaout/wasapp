package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/api/reqcontext"
)

func (rt *_router) getUserAvatar(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Extract target user id from url (users/{userId}/avatar)
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

	// If the user doesn't have an avatar, return 404.
	if imagePath == "" {
		rt.respondWithError(w, http.StatusNotFound, ErrCodeNoUserAvatar, "user has no avatar") //404
		return
	}

	http.ServeFile(w, r, imagePath)
}
