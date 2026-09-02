package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/api/reqcontext"
	"github.com/rimaout/wasapp/service/database"
)

type AddUserToGroupRequest struct {
	UserID string `json:"userId"`
}

func (rt *_router) addMemberToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Extract target chat id from url (chats/{chatId}/members)
	chatId := ps.ByName("chatId")

	// Extract new member id from request body (JSON)
	var req AddUserToGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rt.respondWithError(w, http.StatusBadRequest, ErrCodeInvalidInput, "invalid request body") //400
		return
	}
	newMemberId := req.UserID

	// Validate that the authenticated user has access to the chat (404, 403, 500 errors)
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
		rt.respondWithError(w, http.StatusConflict, ErrCodeNotAGroupChat, "chat is not a group chat") //409
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
	isAlreadyMember, err := rt.db.IsActiveChatMember(chatId, newMemberId)
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
	if err := rt.db.AddChatMember(chatId, newMemberId); err != nil {
		ctx.Logger.WithError(err).Error("error adding new user to group")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Create a join system message (sender is the new member)
	joinMsgId, err := rt.db.CreateSystemMessage(chatId, newMemberId, true)
	if err != nil {
		ctx.Logger.WithError(err).Error("error creating join message")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Create the receiver statuses for the join message (one for each member of the chat)
	_ = rt.db.InsertReceiverStatuses(joinMsgId, chatId, newMemberId)

	// Get updated list of members
	updatedMembers, err := rt.db.GetChatMembers(chatId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting updated list of members")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Respond with success (201) and the updated list of members
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		MembersList []database.Member `json:"membersList"`
	}{MembersList: updatedMembers})
}

