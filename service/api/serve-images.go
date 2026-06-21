package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/api/reqcontext"
)

func (rt *_router) getUserProfileImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userId := ps.ByName("userId")

	userName, err := rt.db.GetUserNameById(userId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking if user exists")
		rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
		return
	}
	if userName == "" {
		rt.respondWithError(w, http.StatusNotFound, "404", "user not found")
		return
	}

	imagePath, err := rt.db.GetUserProfileImagePath(userId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting profile image path")
		rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
		return
	}
	if imagePath == "" {
		rt.respondWithError(w, http.StatusNotFound, "404", "profile image not found")
		// TODO: this should return a default image instead of 404
		return
	}

	http.ServeFile(w, r, imagePath)
}
