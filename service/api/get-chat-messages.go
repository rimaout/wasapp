package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/api/reqcontext"
	"github.com/rimaout/wasapp/service/database"
)

type chatMessagesResponse struct {
	Messages []database.Message `json:"messages"`
}

func (rt *_router) getChatMessages(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Extract target chat id form url (chats/{chatId}/messages)
	chatId := ps.ByName("chatId")

	// Validate that the authenticated user has access to the chat (404, 403, 500 errors)
	if !rt.validateChatAccess(w, r, ctx, chatId) {
		// Checks errors: 404 (chat not found), 403 (user not a member), 500 (internal error)
		return
	}

	// Mark all messages in this chat as received by the requesting user,
	// before computing their statuses so the response reflects it.
	_ = rt.db.MarkMessagesReceivedByUser(chatId, ctx.UserID)

	// Get chat messages slice from the database
	messages, err := rt.db.GetChatMessages(chatId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting chat messages")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Return the messages as JSON
	resp := chatMessagesResponse{Messages: messages}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
