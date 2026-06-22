package api

import (
	"encoding/json"
	"net/http"

	"github.com/rimaout/wasapp/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type patchUserNameRequest struct {
	UserName string `json:"userName"`
}

type userResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Read request
	var req patchUserNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rt.respondWithError(w, http.StatusBadRequest, "400", "invalid request body")
		return
	}

	// Check if new username is valied (respects the basename rules)
	if !isValidBaseName(req.UserName) {
		rt.respondWithError(w, http.StatusBadRequest, "400", "userName must be 3-24 characters, alphanumeric + spaces/underscores/hyphens, at least one non-space")
		return
	}

	// Check if new username is not already used
	existingId, err := rt.db.GetUserIdByName(req.UserName)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking user name availability")
		rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
		return
	}
	if existingId != "" && existingId != ctx.UserID {
		rt.respondWithError(w, http.StatusConflict, "409", "username already in use")
		return
	}

	// Set new username
	if err := rt.db.SetUserName(ctx.UserID, req.UserName); err != nil {
		ctx.Logger.WithError(err).Error("error updating user name")
		rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
		return
	}

	// Read updated username
	newName, err := rt.db.GetUserNameById(ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting updated user name")
		rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
		return
	}

	// Return userId and userName (TODO: add retuen profile image)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userResponse{
		Id:   ctx.UserID,
		Name: newName,
	})
}
