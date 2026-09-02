package api

import (
	"encoding/json"
	"net/http"
)

// Machine-readable error codes returned in the JSON error payload.
// NOTE: To keep in sync with the Error schema enum in doc/api.yaml.
const (
	ErrCodeChatNotFound			  = "CHAT-NOT-FOUND"
	ErrCodeUserNotFound			  = "USER-NOT-FOUND"
	ErrCodeAlreadyInGroup		  = "ALREADY-IN-GROUP"
	ErrCodeNotAGroupChat		  = "NOT-A-GROUP-CHAT"
	ErrCodeNoChatAvatar			  = "NO-CHAT-AVATAR"
	ErrCodeNoUserAvatar			  = "NO-USER-AVATAR"
	ErrCodeForbiddenNotMember	  = "FORBIDDEN-NOT-MEMBER"
	ErrCodeInvalidInput			  = "INVALID-INPUT"
	ErrCodeInternalError		  = "INTERNAL-SERVER-ERROR"
	ErrCodeUnauthorized			  = "UNAUTHORIZED"
	ErrCodeUsernameTaken		  = "USERNAME-TAKEN"
	ErrCodeChatAlreadyExists	  = "CHAT-ALREADY-EXISTS"
	ErrCodePayloadTooLarge		  = "PAYLOAD-TOO-LARGE"
	ErrCodeUnsupportedMedia		  = "UNSUPPORTED-MEDIA-TYPE"
	ErrCodeInvalidUserName		  = "INVALID-USER-NAME"
	ErrCodeInvalidGroupName		  = "INVALID-GROUP-NAME"
	ErrCodeEmptyMemberList		  = "EMPTY-MEMBER-LIST"
	ErrCodeMissingImageFile		  = "MISSING-IMAGE-FILE"
	ErrCodeForbiddenNotSender	  = "FORBIDDEN-NOT-SENDER"
	ErrCodeAlreadyDeleted		  = "MESSAGE-ALREADY-DELETED"
	ErrCodeItsInitMessage		  = "ITS-INIT-MESSAGE"
	ErrCodeMessageNotFound        = "MESSAGE-NOT-FOUND"
	ErrCodeInvalidEmojiId         = "INVALID-EMOJI-ID"
	ErrCodeReactionNotFound       = "REACTION-NOT-FOUND"
	ErrDestinationChatNotFound    = "DESTINATION-CHAT-NOT-FOUND"
	ErrCodeOriginChatNotFound	  = "ORIGIN-CHAT-NOT-FOUND"
	ErrCodeInvalidMessageText     = "INVALID-MESSAGE-TEXT"
	ErrCodeMessageWithNoContent   = "MESSAGE-WITH-NO-CONTENT"
	ErrCodeCannotForwardForwarded = "CANNOT-FORWARD-FORWARDED-MESSAGE"

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
