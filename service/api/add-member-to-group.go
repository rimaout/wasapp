package api

import (
	"net/http"
	"encoding/json"
	"github.com/julienschmidt/httprouter"

	"github.com/rimaout/wasapp/service/database"
	"github.com/rimaout/wasapp/service/api/reqcontext"
)

type AddUserToGroupRequest struct {
	UserID string `json:"userId"`
}

func (rt *_router) addMemberToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Extract target chat id from url (http params)
	targetChatId := ps.ByName("chatId")

	// Extract user id from context (logged user)
	loggedUserId := ctx.UserID

	// Extract new member id from request body (JSON)
	var req AddUserToGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rt.respondWithError(w, http.StatusBadRequest, ErrCodeInvalidInput, "invalid request body") //400
		return
	}
	newMemberId := req.UserID

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
		rt.respondWithError(w, http.StatusConflict, ErrCodeNotAGroupChat, "chat is not a group chat") //409
		return
	}

	// Check if logged user ia a member of the group chat (403 error)
	isMember, err := rt.db.IsActiveChatMember(loggedUserId, targetChatId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking if user is a member of the chat")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}
	if !isMember {
		rt.respondWithError(w, http.StatusForbidden, ErrCodeForbiddenNotMember, "you are not a member of this group chat") //403
		return
	}

	// Check if new member exists (404)
	newMemberName, err := rt.db.GetUserNameById(newMemberId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting user by id")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}
	if newMemberName == "" {
		rt.respondWithError(w, http.StatusNotFound, ErrCodeUserNotFound, "user not found") //404
		return
	}

	// Check if new member is already a member of the group (409)
	isAlreadyMember, err := rt.db.IsActiveChatMember(newMemberId, targetChatId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking if user is a member of the chat")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}
	if isAlreadyMember {
		rt.respondWithError(w, http.StatusConflict, ErrCodeAlreadyInGroup, "user is already a member of the group chat") //409
		return
	}

	// Add new member to group
	if err := rt.db.AddChatMember(targetChatId, newMemberId); err != nil {
		ctx.Logger.WithError(err).Error("error adding new user to group")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Get updated list of members
	updatedMembers, err := rt.db.GetChatMembers(targetChatId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting updated list of members")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Respond with success (201) and the updated list of members
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(updatedMembers)
}

