package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/database"
	"github.com/rimaout/wasapp/service/api/reqcontext"
)

type patchGroupNameRequest struct {
	GroupName string `json:"userName"`
}

type GroupNameResponse struct {
	GroupName string `json:"userName"`
}

func (rt *_router) setGroupChatName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Extract target chat id form url (http params)
	targetChatId := ps.ByName("chatId")

	// Read request (to extract new group name)
	var req patchGroupNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rt.respondWithError(w, http.StatusBadRequest, ErrCodeInvalidInput, "invalid request body") //400
		return
	}
	newGroupName := req.GroupName

	// Check if new group name is valid (400 error)
	if !isValidBaseName(newGroupName) {
		rt.respondWithError(w, http.StatusBadRequest, ErrCodeInvalidInput, "groupName must be 3-24 characters, alphanumeric + spaces/underscores/hyphens, at least one non-space") //400
		return
	}

	// Check if chat exist (404 error)
	_, err := rt.db.GetChatById(targetChatId)
	if err == database.ErrChatNotFound {
		rt.respondWithError(w, http.StatusNotFound, ErrCodeChatNotFound, "chat not found") //404
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting chat by id")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Check if chat is a group chat (409 error)
	isGroup, err := rt.db.IsGroupChat(targetChatId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking if chat is a group chat")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}
	if !isGroup {
		rt.respondWithError(w, http.StatusConflict, ErrCodeNotAGroupChat, "chat is not a group chat, only the names of group chats can be set/changed") //409
		return
	}

	// Check if logged user ia a member of the group chat (403 error)
	isMember, err := rt.db.IsActiveChatMember(ctx.UserID, targetChatId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking if user is a member of the chat")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}
	if !isMember {
		rt.respondWithError(w, http.StatusForbidden, ErrCodeForbiddenNotMember, "you are not a member of this group chat") //403
		return
	}

	// Set new group name
	if err := rt.db.SetGroupName(targetChatId, newGroupName); err != nil {
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

