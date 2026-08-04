package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/structures"
	"github.com/rimaout/wasapp/service/api/reqcontext"
)

type groupChatRequest struct {
	GroupName   string   `json:"groupName"`
	MembersList []string `json:"membersList"`
}

func (rt *_router) createGroupChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var req groupChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rt.respondWithError(w, http.StatusBadRequest, ErrCodeInvalidInput, "invalid request body") //400
		return
	}

	// Check group name
	if !isValidBaseName(req.GroupName) {
		rt.respondWithError(w, http.StatusBadRequest, ErrCodeInvalidGroupName, "groupName must be 3-24 characters, alphanumeric + spaces/underscores/hyphens, at least one non-space") //400
		return
	}

	// Check members list: at least 1 member
	if len(req.MembersList) == 0 {
		rt.respondWithError(w, http.StatusBadRequest, ErrCodeEmptyMemberList, "membersList must contain at least one user") //400
		return
	}

	// Check members list: no duplicates
	seen := make(map[string]bool)
	for _, id := range req.MembersList {
		if seen[id] {
			rt.respondWithError(w, http.StatusConflict, ErrCodeAlreadyInGroup, "duplicate user in membersList") //409
			return
		}
		seen[id] = true
	}

	// Check members list: all users exist
	for _, id := range req.MembersList {
		// From user ID, get user name to check if user exists
		name, err := rt.db.GetUserNameById(id)
		if err != nil {
			ctx.Logger.WithError(err).Error("error looking up user in members list")
			rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
			return
		}

		// If name is empty, user does not exist in the database
		if name == "" {
			rt.respondWithError(w, http.StatusNotFound, ErrCodeUserNotFound, "user not found: "+id) //404
			return
		}
	}

	// Create chat
	chatId, err := rt.db.CreateChat(true, req.GroupName, "")
	if err != nil {
		ctx.Logger.WithError(err).Error("error creating group chat")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Add creator (logged user) as member
	if err := rt.db.AddChatMember(chatId, ctx.UserID); err != nil {
		ctx.Logger.WithError(err).Error("error adding creator to group")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// add all members from meber list
	for _, id := range req.MembersList {
		if id == ctx.UserID {
			continue // already added as creator
		}
		if err := rt.db.AddChatMember(chatId, id); err != nil {
			ctx.Logger.WithError(err).Error("error adding member to group")
			rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
			return
		}
	}

	// Create init message for group chat
	//	 The init message is a message with out content, just to have a message in the chat.
	//	 The frontend can use it to display "Group created by ..." or similar in the chat history.
	_, err = rt.db.CreateMessage(
		chatId,
		ctx.UserID,
		time.Now().UTC().Format(time.DateTime),
		"",
		"",
		"",
		true,
		false,
		"",
		"",
	)
	if err != nil {
		ctx.Logger.WithError(err).Error("error creating init message")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(structures.ChatPreviewResponse{
		Id:          chatId,
		DisplayName: req.GroupName,
		IsGroupChat: true,
	})
}
