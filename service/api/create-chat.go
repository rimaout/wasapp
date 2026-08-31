package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/api/reqcontext"
)

type groupChatRequest struct {
	GroupName   string   `json:"groupName"`
	MembersList []string `json:"membersList"`
}

type ChatPreviewResponse struct {
	Id          string `json:"id"`
	DisplayName string `json:"displayName"`
	IsGroupChat bool   `json:"isGroupChat"`
}

func (rt *_router) createGroupChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Decode the request body
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
	chatId, err := rt.db.CreateChat(true, req.GroupName)
	if err != nil {
		ctx.Logger.WithError(err).Error("error creating group chat")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Add creator (authenticated user) as member
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
	messageId, err := rt.db.CreateMessage(
		chatId,
		ctx.UserID,
		time.Now().UTC().Format(time.RFC3339),
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

	// Create the receiver statuses for the message (one for each member of the chat)
	// This means that the message is considered "sent" to all members, but not yet "delivered" or "read"
	_ = rt.db.InsertReceiverStatuses(messageId, chatId, ctx.UserID)

	// Respond with the chat preview for the newly created group chat
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ChatPreviewResponse{
		Id:          chatId,
		DisplayName: req.GroupName,
		IsGroupChat: true,
	})
}

func (rt *_router) createDirectChat(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get userId from endpoint url (/user/{userId}/chats)
	targetUserId := ps.ByName("userId")

	// Check user exists
	targetName, err := rt.db.GetUserNameById(targetUserId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error looking up target user")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}
	if targetName == "" {
		rt.respondWithError(w, http.StatusNotFound, ErrCodeUserNotFound, "user not found") //404
		return
	}

	// Cannot create a chat with yourself
	if targetUserId == ctx.UserID {
		rt.respondWithError(w, http.StatusBadRequest, ErrCodeInvalidInput, "cannot create a chat with yourself") //400
		return
	}

	// Check if private chat already exists between these two users
	existingChatId, err := rt.db.FindPrivateChatBetween(ctx.UserID, targetUserId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking existing private chat")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}
	if existingChatId != "" {
		rt.respondWithChatAlreadyExists(w, existingChatId)
		return
	}

	// Create chat
	chatId, err := rt.db.CreateChat(false, "")
	if err != nil {
		ctx.Logger.WithError(err).Error("error creating chat")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Add both members
	if err := rt.db.AddChatMember(chatId, ctx.UserID); err != nil {
		ctx.Logger.WithError(err).Error("error adding creator to chat")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}
	if err := rt.db.AddChatMember(chatId, targetUserId); err != nil {
		ctx.Logger.WithError(err).Error("error adding target user to chat")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Create init message (sender=creator, is_init_message=1, no content)
	messageId, err := rt.db.CreateMessage(
		chatId,
		ctx.UserID,
		time.Now().UTC().Format(time.RFC3339),
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

	// Create the receiver statuses for the message (one for each member of the chat)
	// This means that the message is considered "sent" to all members, but not yet "delivered" or "read"
	_ = rt.db.InsertReceiverStatuses(messageId, chatId, ctx.UserID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(ChatPreviewResponse{
		Id:          chatId,
		DisplayName: targetName,
		IsGroupChat: false,
	})
}
