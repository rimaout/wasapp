package api

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/rimaout/wasapp/service/api/reqcontext"
	"github.com/rimaout/wasapp/service/database"
)

func (rt *_router) getMessageImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Extract target chat id and image id from url (chats/{chatId}/messages/images/{imageId})
	chatId := ps.ByName("chatId")
	imageId := ps.ByName("imageId")

	// Check if chat exists (404 error)
	_, err := rt.db.GetChatById(chatId)
	if errors.Is(err, database.ErrChatNotFound) {
		rt.respondWithError(w, http.StatusNotFound, ErrCodeChatNotFound, "chat not found") // 404
		return
	}
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting chat by id")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") // 500
		return
	}

	// Check if authenticated user is member of the chat (403 error)
	isMember, err := rt.db.IsActiveChatMember(chatId, ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking if user is member of chat")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") // 500
		return
	}
	if !isMember {
		rt.respondWithError(w, http.StatusForbidden, ErrCodeForbiddenNotMember, "user is not a member of the chat") // 403
		return
	}

	// Check if image is associated with a message in this chat (404 error)
	exists, err := rt.db.IsImageVisibleInChat(chatId, imageId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error checking if image is in chat")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") // 500
		return
	}
	if !exists {
		rt.respondWithError(w, http.StatusNotFound, ErrCodeMissingImageFile, "image not found") // 404
		return
	}

	// Get the image path
	imagePath, err := rt.db.GetImagePathByImageId(imageId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting image path")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") // 500
		return
	}
	if imagePath == "" {
		rt.respondWithError(w, http.StatusNotFound, ErrCodeMissingImageFile, "image not found") // 404
		return
	}

	http.ServeFile(w, r, imagePath)
}
