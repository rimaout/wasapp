package api

import (
	"encoding/json"
	"net/http"

	"github.com/rimaout/wasapp/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
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
		rt.respondWithError(w, http.StatusBadRequest, "400", "invalid request body")
		return
	}

	// Validate the userName using the same rules as PATCH /profile/name
	if !isValidBaseName(req.UserName) {
		rt.respondWithError(w, http.StatusBadRequest, "400", "userName must be 3-24 characters, alphanumeric + spaces/underscores/hyphens, at least one non-space")
		return
	}

	// Check if the user already exists in the database
	userId, err := rt.db.GetUserIdByName(req.UserName)
	if err != nil {
		ctx.Logger.WithError(err).Error("error looking up user")
		rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
		return
	}
	status := http.StatusOK	// Default status is 200 OK

	// If the user does not exist, create a new user
	if userId == "" {
		userId, err = rt.db.CreateUser(req.UserName)
		if err != nil {
			ctx.Logger.WithError(err).Error("error creating user")
			rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
			return
		}
		status = http.StatusCreated // Set status to 201 Created if a new user was created
	}

	// Generate Session (Bearer) Token
	token, err := rt.db.GenerateUserSessionToken(userId)
    if err != nil {
        ctx.Logger.WithError(err).Error("error generating session token")
		rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
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
