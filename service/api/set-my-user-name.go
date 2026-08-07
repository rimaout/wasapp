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
		rt.respondWithError(w, http.StatusBadRequest, ErrCodeInvalidInput, "invalid request body") //400
		return
	}

	// Check if new username is valied (respects the basename rules)
	if !isValidBaseName(req.UserName) {
		rt.respondWithError(w, http.StatusBadRequest, ErrCodeInvalidUserName, "userName must be 3-24 characters, alphanumeric + spaces/underscores/hyphens, at least one non-space") //400
		return
	}

	// Check if new username is not already used
	existingId, err := rt.db.GetUserIdByName(req.UserName)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking user name availability")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}
	if existingId != "" && existingId != ctx.UserID {
		rt.respondWithError(w, http.StatusConflict, ErrCodeUsernameTaken, "username already in use") //409
		return
	}

	// Set new username
	if err := rt.db.SetUserName(ctx.UserID, req.UserName); err != nil {
		ctx.Logger.WithError(err).Error("error updating user name")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Read updated username
	newName, err := rt.db.GetUserNameById(ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting updated user name")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Return userId and userName
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userResponse{
		Id:   ctx.UserID,
		Name: newName,
	})
}
