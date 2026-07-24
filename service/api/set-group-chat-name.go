package api

import (
	"encoding/json"
	"net/http"

	"github.com/rimaout/wasapp/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type patchGroupNameRequest struct {
	GroupName string `json:"userName"`
}

type GroupNameResponse struct {
	GroupName string `json:"userName"`
}

func (rt *_router) setGroupNameName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Extract target chat id form url (http params)
	targetChatId := ps.ByName("chatId")

	// Read request (to extract new group name)
	var req patchGroupNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rt.respondWithError(w, http.StatusBadRequest, "400", "invalid request body")
		return
	}
	newGroupName := req.GroupName

	// Check if new group name is valid (respects the basename rules)
	if !isValidBaseName(newGroupName) {
		rt.respondWithError(w, http.StatusBadRequest, "400", "groupName must be 3-24 characters, alphanumeric + spaces/underscores/hyphens, at least one non-space")
		return
	}

	// TODO: Check if chat exist (404 error)

	// TODO: Check if chat is a group chat (409 error)

	// TODO: Check if logged user ia a member of the group chat (403 error)

	// Set new group name
	if err := rt.db.SetGroupName(targetChatId, newGroupName); err != nil {
		ctx.Logger.WithError(err).Error("error updating group name")
		rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
		return
	}

	// Respond with the new group name
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(GroupNameResponse{
		GroupName: req.GroupName,
	})
}

