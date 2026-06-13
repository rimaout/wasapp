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
	Identifier string `json:"identifier"`
}

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	// Decode the request body into a loginRequest struct
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"code":"400","message":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	// Validate the userName length
	if len(req.UserName) < 3 || len(req.UserName) > 24 {
		http.Error(w, `{"code":"400","message":"userName must be 3-24 characters"}`, http.StatusBadRequest)
		return
	}

	// Check if the user already exists in the database
	userId, err := rt.db.GetUserByName(req.UserName)
	if err != nil {
		ctx.Logger.WithError(err).Error("error looking up user")
		http.Error(w, `{"code":"500","message":"internal server error"}`, http.StatusInternalServerError)
		return
	}
	status := http.StatusOK	// Default status is 200 OK

	// If the user does not exist, create a new user
	if userId == "" {
		userId, err = rt.db.CreateUser(req.UserName)
		if err != nil {
			ctx.Logger.WithError(err).Error("error creating user")
			http.Error(w, `{"code":"500","message":"internal server error"}`, http.StatusInternalServerError)
			return
		}
		status = http.StatusCreated // Set status to 201 Created if a new user was created
	}

	// Respond with the user identifier in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(loginResponse{Identifier: userId})
}
