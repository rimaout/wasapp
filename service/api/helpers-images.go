package api

import (
	"io"
	"os"
	"errors"
	"net/http"
	"path/filepath"
	
	"github.com/gofrs/uuid"
	
	"github.com/rimaout/wasapp/service/api/reqcontext"
)

const maxImageSize = 5 * 1024 * 1024

// deleteImageFile deleat a image file from disk and logs warnings if it fails.
func (rt *_router) deleteImageFile(ctx reqcontext.RequestContext, filePath string) {
	if filePath == "" {
		return
	}
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		ctx.Logger.WithError(err).Warnf("failed to delete file from disk: %s", filePath)
	}
}

// saveUploadedImage extracts, validates, and saves an input image.
func (rt *_router) saveUploadedImage(
	w http.ResponseWriter,
	r *http.Request,
	ctx reqcontext.RequestContext,
	formField, targetDir string,
) (string, bool) {
	// Limit the size of the incoming request body
	r.Body = http.MaxBytesReader(w, r.Body, maxImageSize)

	// Parse the multipart form to determine if is correctly formed and not too large.
	if err := r.ParseMultipartForm(maxImageSize); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			rt.respondWithError(w, http.StatusRequestEntityTooLarge, "413", "image exceeds maximum size of 5MB")
			return "", false
		}

		rt.respondWithError(w, http.StatusBadRequest, "400", "invalid multipart form")
		return "", false
	}

	// Extract the file from the form data
	file, header, err := r.FormFile(formField)
	if err != nil {
		rt.respondWithError(w, http.StatusBadRequest, "400", "missing image file")
		return "", false
	}
	defer file.Close()

	// Validate content type (only JPEG, PNG, WEBP are allowed)
	contentType := header.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" {
		rt.respondWithError(w, http.StatusUnsupportedMediaType, "415", "unsupported image format, only JPEG, PNG, WebP allowed")
		return "", false
	}

	// Ensure the target directory exists
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		ctx.Logger.WithError(err).Error("error creating upload directory")
		rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
		return "", false
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

	// Generate a random image name (path)
	imageId, err := uuid.NewV4()
	if err != nil {
		ctx.Logger.WithError(err).Error("error generating image UUID")
		rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
		return "", false
	}
	imagePath := filepath.Join(targetDir, imageId.String()+ext)

	// Create the destination file
	dst, err := os.Create(imagePath)
	if err != nil {
		ctx.Logger.WithError(err).Error("error creating image file")
		rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
		return "", false
	}
	defer dst.Close()

	// Copy the uploaded file to the destination file
	if _, err := io.Copy(dst, file); err != nil {
		ctx.Logger.WithError(err).Error("error saving image file")
		rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
		return "", false
	}

	return imagePath, true
}
