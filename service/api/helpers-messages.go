package api

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"

	"github.com/rimaout/wasapp/service/api/reqcontext"
)

const messageImageUploadDir = "messages/images"
const maxMessagePayloadSize = maxImageSize + (1 * 1024 * 1024) // 5MB max image + 1MB buffer for text/headers

type NewMessageContent struct {
	Text      string
	ImagePath string
}

// extractSendMessageRequest parses the multipart form, processes the image if present,
func (rt *_router) extractSendMessageContent(
	w http.ResponseWriter,
	r *http.Request,
	ctx reqcontext.RequestContext,
) (NewMessageContent, bool) {
	// Limit the size of the incoming request body
	r.Body = http.MaxBytesReader(w, r.Body, maxMessagePayloadSize)

	// Parse the multipart form to determine if it is correctly formed and not too large.
	if err := r.ParseMultipartForm(maxMessagePayloadSize); err != nil {

		// Check if request body exceeds the maximum size limit
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			rt.respondWithError(w, http.StatusRequestEntityTooLarge, ErrCodePayloadTooLarge, "image exceeds maximum size of 5MB") // 413
			return NewMessageContent{}, false
		}

		// Check if the error is due to invalid multipart form format
		rt.respondWithError(w, http.StatusBadRequest, ErrCodeInvalidInput, "invalid multipart form") // 400
		return NewMessageContent{}, false
	}

	// Extract fields
	text := r.FormValue("text")
	file, header, err := r.FormFile("imageFile")

	// Check for errors in reading the image file, but allow for the case where no file was uploaded
	if err != nil && !errors.Is(err, http.ErrMissingFile) {
		ctx.Logger.WithError(err).Error("error reading message image file")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") // 500
		return NewMessageContent{}, false
	}

	// If an image was attached, process and save it
	var imagePath string
	if err == nil {
		defer file.Close()

		// Validate content type (JPEG, PNG, WebP)
		contentType := header.Header.Get("Content-Type")
		if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" {
			rt.respondWithError(w, http.StatusUnsupportedMediaType, ErrCodeUnsupportedMedia, "unsupported image format, only JPEG, PNG, WebP allowed") // 415
			return NewMessageContent{}, false
		}

		// Ensure target directory exists
		targetDir := filepath.Join(rt.uploadsDir, messageImageUploadDir)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			ctx.Logger.WithError(err).Error("error creating message upload directory")
			rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") // 500
			return NewMessageContent{}, false
		}

		// Determine file extension
		ext := filepath.Ext(header.Filename)
		if ext == "" {
			switch contentType {
			case "image/jpeg":
				ext = ".jpg"
			case "image/png":
				ext = ".png"
			case "image/webp":
				ext = ".webp"
			}
		}

		// Generate random UUID for filename
		imageId, err := uuid.NewV4()
		if err != nil {
			ctx.Logger.WithError(err).Error("error generating image UUID")
			rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") // 500
			return NewMessageContent{}, false
		}
		imagePath = filepath.Join(targetDir, imageId.String()+ext)

		// Create destination file
		dst, err := os.Create(imagePath)
		if err != nil {
			ctx.Logger.WithError(err).Error("error creating message image file")
			rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") // 500
			return NewMessageContent{}, false
		}
		defer dst.Close()

		// Write content to disk (note: io.Copy directly writes to the disk without buffering the entire file in memory)
		if _, err := io.Copy(dst, file); err != nil {
			ctx.Logger.WithError(err).Error("error saving message image file")
			rt.deleteImageFile(ctx, imagePath)                                                                    // Cleanup partial file
			rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") // 500
			return NewMessageContent{}, false
		}
	}

	// Validate text content if present
	if text != "" {
		if valid, errMsg := isValidMessageText(text); !valid {
			rt.respondWithError(w, http.StatusBadRequest, ErrCodeInvalidMessageText, errMsg) // 400
			return NewMessageContent{}, false
		}
	}

	// Check for "oneOf" between image and text (must contain text, image, or both)
	if text == "" && imagePath == "" {
		rt.respondWithError(w, http.StatusBadRequest, ErrCodeMessageWithNoContent, "message must contain text, an image, or both") // 400
		return NewMessageContent{}, false
	}

	return NewMessageContent{
		Text:      text,
		ImagePath: imagePath,
	}, true
}

// isValidMessageText validates message text, it must be 1-3000 characters and contain at least one non-whitespace character.
func isValidMessageText(text string) (bool, string) {
	if len(text) > 3000 {
		return false, "text exceeds maximum length of 3000 characters"
	}
	nonSpace := false
	for _, c := range text {
		if c != ' ' && c != '\t' && c != '\n' && c != '\r' {
			nonSpace = true
			break
		}
	}
	if !nonSpace {
		return false, "text cannot be only whitespace"
	}
	return true, ""
}
