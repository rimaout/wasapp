package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/structures"
	"github.com/rimaout/wasapp/service/api/reqcontext"
)

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
		rt.respondWithError(w, http.StatusConflict, ErrCodeChatAlreadyExists, "you already have a direct chat with this user") //409
		return
	}

	// Create chat
	chatId, err := rt.db.CreateChat(false, "", "")
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
	_, err = rt.db.CreateMessage(
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

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(structures.ChatPreviewResponse{
		Id:          chatId,
		DisplayName: targetName,
		IsGroupChat: false,
	})
}
