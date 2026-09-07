package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rimaout/wasapp/service/api/reqcontext"
)

type loginRequest struct {
	UserName string `json:"userName"`
}

type loginResponse struct {
	UserId string `json:"userId"`
	Token  string `json:"token"`
}

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Decode the request body into a loginRequest struct
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rt.respondWithError(w, http.StatusBadRequest, ErrCodeInvalidInput, "invalid request body") // 400
		return
	}

	// Validate the userName using the same rules as PATCH /profile/name
	if !isValidBaseName(req.UserName) {
		rt.respondWithError(w, http.StatusBadRequest, ErrCodeInvalidUserName, "userName must be 3-24 characters, alphanumeric + spaces/underscores/hyphens, at least one non-space") // 400
		return
	}

	// Check if the user already exists in the database
	userId, err := rt.db.GetUserIdByName(req.UserName)
	if err != nil {
		ctx.Logger.WithError(err).Error("error looking up user")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") // 500
		return
	}
	status := http.StatusOK // Default status is 200 OK

	// If the user does not exist, create a new user
	if userId == "" {
		userId, err = rt.db.CreateUser(req.UserName)
		if err != nil {
			ctx.Logger.WithError(err).Error("error creating user")
			rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") // 500
			return
		}
		status = http.StatusCreated // Set status to 201 Created if a new user was created
	}

	// Generate Session (Bearer) Token
	token, err := rt.db.GenerateUserSessionToken(userId)
	if err != nil {
		ctx.Logger.WithError(err).Error("error generating session token")
		rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") // 500
		return
	}

	// Respond with the user identifier and token in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(loginResponse{
		UserId: userId,
		Token:  "Bearer " + token,
	})
}
