package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/database"
	"github.com/rimaout/wasapp/service/api/reqcontext"
)

type forwardMessageRequest struct {
	ForwardTo string `json:"forwardTo"`
}

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Extract the origin chat ID and message ID from the URL parameters (chats/{chatId}/messages/{messageId}/forward)
	originChatId := ps.ByName("chatId")
	originMessageId := ps.ByName("messageId")

	// Check if the origin chat exists
	_, err := rt.db.GetChatById(originChatId)
	if err == database.ErrChatNotFound {
		rt.respondWithError(w, http.StatusNotFound, ErrCodeOriginChatNotFound, "origin chat not found") //404
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting origin chat by id")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Check if the user is a member of the origin chat
	isMember, err := rt.db.IsActiveChatMember(originChatId, ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking origin chat membership")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}
	if !isMember {
		rt.respondWithError(w, http.StatusForbidden, ErrCodeForbiddenNotMember, "user is not a member of the origin chat") //403
		return
	}

	// Check if the origin message exists
	originMsg, err := rt.db.GetMessageById(originChatId, originMessageId)
	if err == database.ErrMessageNotFound {
		rt.respondWithError(w, http.StatusNotFound, ErrCodeMessageNotFound, "message not found") //404
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting origin message")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Check if the origin message is already a forwarded message (we don't allow forwarding forwarded messages)
	if originMsg.ForwardedFrom != nil {
		rt.respondWithError(w, http.StatusBadRequest, ErrCodeCannotForwardForwarded, "cannot forward a forwarded message") //400
		return
	}

	// Decode the request body
	var req forwardMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rt.respondWithError(w, http.StatusBadRequest, ErrCodeInvalidInput, "invalid request body") //400
		return
	}
	if req.ForwardTo == "" {
		rt.respondWithError(w, http.StatusBadRequest, ErrCodeInvalidInput, "destination chat id is required") //400
		return
	}

	// Check if the destination chat is the same as the origin chat (we don't allow forwarding to the same chat)
	if req.ForwardTo == originChatId {
		rt.respondWithError(w, http.StatusBadRequest, ErrCodeInvalidInput, "cannot forward to the same chat") //400
		return
	}

	// Check if the destination chat exists
	_, err = rt.db.GetChatById(req.ForwardTo)
	if err == database.ErrChatNotFound {
		rt.respondWithError(w, http.StatusNotFound, ErrDestinationChatNotFound, "destination chat not found") //404
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting destination chat by id")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Check if the user is a member of the destination chat
	isDestMember, err := rt.db.IsActiveChatMember(req.ForwardTo, ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking destination chat membership")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}
	if !isDestMember {
		rt.respondWithError(w, http.StatusForbidden, ErrCodeForbiddenNotMember, "user is not a member of the destination chat") //403
		return
	}

	// If the original message has an image, insert a record in the image visibility table for the forwarded message
	if originMsg.Content != nil && originMsg.Content.MsgImageId != nil {
		_ = rt.db.InsertImageVisibility(*originMsg.Content.MsgImageId, req.ForwardTo)
	}

	// Create the forwarded message
	sendTime := time.Now().UTC().Format(time.RFC3339)
	messageId, err := rt.db.CreateMessage(
		req.ForwardTo,
		ctx.UserID,
		sendTime,
		"",
		"",
		"",
		false,
		true,
		originChatId,
		originMessageId,
	)
	if err != nil {
		ctx.Logger.WithError(err).Error("error creating forwarded message")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Fetch the created forward message to return in the response
	msg, err := rt.db.GetMessageById(req.ForwardTo, messageId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error fetching forwarded message")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	// Create the receiver statuses for the message (one for each member of the chat)
	// This means that the message is considered "sent" to all members, but not yet "delivered" or "read"
	_ = rt.db.InsertReceiverStatuses(messageId, req.ForwardTo, ctx.UserID)

	// Respond with the created forwarded message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)
}
