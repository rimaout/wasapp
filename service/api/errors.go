package api

import (
	"encoding/json"
	"net/http"
)

// Machine-readable error codes returned in the JSON error payload.
// NOTE: Keep in sync with the Error schema enum in doc/api.yaml.
const (
	ErrCodeChatNotFound        = "CHAT_NOT_FOUND"
	ErrCodeUserNotFound        = "USER_NOT_FOUND"
	ErrCodeAlreadyInGroup      = "ALREADY_IN_GROUP"
	ErrCodeNotAGroupChat       = "NOT_A_GROUP_CHAT"
	ErrCodeNoChatAvatar        = "NO_CHAT_AVATAR"
	ErrCodeNoUserAvatar        = "NO_USER_AVATAR"
	ErrCodeForbiddenNotMember  = "FORBIDDEN_NOT_MEMBER"
	ErrCodeInvalidInput        = "INVALID_INPUT"
	ErrCodeInternalError       = "INTERNAL_SERVER_ERROR"
	ErrCodeUnauthorized        = "UNAUTHORIZED"
	ErrCodeUsernameTaken       = "USERNAME_TAKEN"
	ErrCodeChatAlreadyExists   = "CHAT_ALREADY_EXISTS"
	ErrCodePayloadTooLarge     = "PAYLOAD_TOO_LARGE"
	ErrCodeUnsupportedMedia    = "UNSUPPORTED_MEDIA_TYPE"
	ErrCodeInvalidUserName     = "INVALID_USER_NAME"
	ErrCodeInvalidGroupName    = "INVALID_GROUP_NAME"
	ErrCodeEmptyMemberList	   = "EMPTY_MEMBER_LIST"
	ErrCodeMissingImageFile    = "MISSING_IMAGE_FILE"
	ErrCodeForbiddenNotSender  = "FORBIDDEN_NOT_SENDER"
	ErrCodeAlreadyDeleted      = "MESSAGE_ALREADY_DELETED"
	ErrCodeItsInitMessage      = "ITS_INIT_MESSAGE"
	ErrCodeMessageNotFound     = "MESSAGE_NOT_FOUND"
	ErrCodeInvalidEmojiId      = "INVALID_EMOJI_ID"
	ErrCodeReactionNotFound    = "REACTION_NOT_FOUND"
	ErrDestinationChatNotFound = "DESTINATION_CHAT_NOT_FOUND"
	ErrCodeOriginChatNotFound  = "ORIGIN_CHAT_NOT_FOUND"
	ErrCodeInvalidMessageText  = "INVALID_MESSAGE_TEXT"
	ErrCodeMessageWithNoContent= "MESSAGE_WITH_NO_CONTENT"
	ErrCodeCannotForwardForwarded = "CANNOT_FORWARD_FORWARDED_MESSAGE"

)

// errorResponse defines a structural type for consistent error responses
type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// respondWithError sets the JSON headers, writes the HTTP status,
// and marshals the custom string code and error message.
func (rt *_router) respondWithError(w http.ResponseWriter, status int, code string, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// We ignore the error here because encoding a basic hardcoded struct
	// to a client response stream will practically never fail.
	_ = json.NewEncoder(w).Encode(errorResponse{
		Code:    code,
		Message: message,
	})
}

// respondWithChatAlreadyExists responds with a 409 Conflict, including the
// ID of the existing chat so the client can open it directly.
func (rt *_router) respondWithChatAlreadyExists(w http.ResponseWriter, chatId string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusConflict)

	// We ignore the error here because encoding a basic hardcoded struct
	// to a client response stream will practically never fail.
	_ = json.NewEncoder(w).Encode(struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		ChatId  string `json:"chatId"`
	}{
		Code:    ErrCodeChatAlreadyExists,
		Message: "you already have a direct chat with this user",
		ChatId:  chatId,
	})
}
