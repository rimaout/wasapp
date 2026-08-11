package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/api/reqcontext"
)

type patchGroupNameRequest struct {
	GroupName string `json:"groupName"`
}

type GroupNameResponse struct {
	GroupName string `json:"groupName"`
}

func (rt *_router) setGroupChatName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Extract target chat id form url (http params)
	chatId := ps.ByName("chatId")

	// Read request (to extract new group name)
	var req patchGroupNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rt.respondWithError(w, http.StatusBadRequest, ErrCodeInvalidInput, "invalid request body") //400
		return
	}
	newGroupName := req.GroupName

	// Check if new group name is valid (400 error)
	if !isValidBaseName(newGroupName) {
		rt.respondWithError(w, http.StatusBadRequest, ErrCodeInvalidGroupName, "groupName must be 3-24 characters, alphanumeric + spaces/underscores/hyphens, at least one non-space") //400
		return
	}

	if !rt.validateChatAccess(w, r, ctx, chatId) {
		// Checks errors: 404 (chat not found), 403 (user not a member), 500 (internal error)
		return
	}

	// Check if chat is a group chat (409 error)
	isGroup, err := rt.db.IsGroupChat(chatId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking if chat is a group chat")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}
	if !isGroup {
		rt.respondWithError(w, http.StatusConflict, ErrCodeNotAGroupChat, "chat is not a group chat, only the names of group chats can be set/changed") //409
		return
	}

	// Set new group name
	if err := rt.db.SetGroupName(chatId, newGroupName); err != nil {
		ctx.Logger.WithError(err).Error("error updating group name")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Respond with the new group name
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(GroupNameResponse{
		GroupName: req.GroupName,
	})
}

