package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/api/reqcontext"
)

type groupChatRequest struct {
	GroupName   string   `json:"groupName"`
	MembersList []string `json:"membersList"`
}

func (rt *_router) createGroupChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var req groupChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rt.respondWithError(w, http.StatusBadRequest, "400", "invalid request body")
		return
	}

	// Check group name
	if !isValidBaseName(req.GroupName) {
		rt.respondWithError(w, http.StatusBadRequest, "400", "groupName must be 3-24 characters, alphanumeric + spaces/underscores/hyphens, at least one non-space")
		return
	}

	// Check members list: at least 1 member
	if len(req.MembersList) == 0 {
		rt.respondWithError(w, http.StatusBadRequest, "400", "membersList must contain at least one user")
		return
	}

	// Check members list: no duplicates
	seen := make(map[string]bool)
	for _, id := range req.MembersList {
		if seen[id] {
			rt.respondWithError(w, http.StatusConflict, "409", "duplicate user in membersList")
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
			rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
			return
		}

		// If name is empty, user does not exist in the database
		if name == "" {
			rt.respondWithError(w, http.StatusBadRequest, "400", "user not found: "+id)
			return
		}
	}

	// Create chat
	chatId, err := rt.db.CreateChat(true, req.GroupName, "")
	if err != nil {
		ctx.Logger.WithError(err).Error("error creating group chat")
		rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
		return
	}

	// Add creator (logged user) as member
	if err := rt.db.AddChatMember(chatId, ctx.UserID); err != nil {
		ctx.Logger.WithError(err).Error("error adding creator to group")
		rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
		return
	}

	// add all members from meber list
	for _, id := range req.MembersList {
		if id == ctx.UserID {
			continue // already added as creator
		}
		if err := rt.db.AddChatMember(chatId, id); err != nil {
			ctx.Logger.WithError(err).Error("error adding member to group")
			rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
			return
		}
	}

	//TODO: create init message for group chat

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(chatPreviewResponse{
		Id:          chatId,
		DisplayName: req.GroupName,
		IsGroupChat: true,
	})
}
