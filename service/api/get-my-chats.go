package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/api/reqcontext"
	"github.com/rimaout/wasapp/service/database"
)

type chatsPreviewListResponse struct {
	ChatsPreviewList []chatPreviewResponse `json:"chatsPreviewList"`
}

type chatPreviewResponse struct {
	ID          string                 `json:"id"`
	DisplayName string                 `json:"displayName"`
	IsGroupChat bool                   `json:"isGroupChat"`
	LastMessage chatLastMessagePreview `json:"lastMessage"`
	UnreadCount int                    `json:"unreadCount"`
}

type chatLastMessagePreview struct {
	Status         string                   `json:"status"`
	SendTime       string                   `json:"sendTime"`
	SenderName     string                   `json:"senderName"`
	IsDeleted      bool                     `json:"isDeleted"`
	IsInitMessage  bool                     `json:"isInitMessage"`
	IsJoinMessage  bool                     `json:"isJoinMessage"`
	IsLeaveMessage bool                     `json:"isLeaveMessage"`
	IsForward      bool                     `json:"isForward"`
	Content        *database.MessageContent `json:"content,omitempty"`
}

func (rt *_router) getMyChats(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Get the list of chats for the authenticated user
	chats, err := rt.db.GetMyChats(ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting chats")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
		return
	}

	list := make([]chatPreviewResponse, 0, len(chats))
	for _, cp := range chats {

		preview := chatPreviewResponse{
			ID:          cp.ChatID,
			DisplayName: cp.DisplayName,
			IsGroupChat: cp.IsGroupChat,
			UnreadCount: cp.UnreadCount,
			LastMessage: chatLastMessagePreview{
				Status:         cp.Status,
				SendTime:       cp.SendTime.UTC().Format(time.RFC3339),
				SenderName:     cp.SenderName,
				IsDeleted:      cp.IsDeleted,
				IsInitMessage:  cp.IsInitMsg,
				IsJoinMessage:  cp.IsJoinMsg,
				IsLeaveMessage: cp.IsLeaveMsg,
				IsForward:      cp.IsForward,
			},
		}

		// Only include the content if the message is not deleted, not an init message, and not a system message
		if !cp.IsDeleted && !cp.IsInitMsg && !cp.IsJoinMsg && !cp.IsLeaveMsg {
			content := &database.MessageContent{}
			hasContent := false
			if cp.Text != nil {
				content.Text = cp.Text
				hasContent = true
			}
			if cp.ImageID != nil {
				content.MsgImageId = cp.ImageID
				hasContent = true
			}
			if hasContent {
				preview.LastMessage.Content = content
			}
		}

		list = append(list, preview)
	}

	resp := chatsPreviewListResponse{
		ChatsPreviewList: list,
	}

	// Respond with the list of chats
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
