package api

import (
	"encoding/json"
	"net/http"

	"github.com/rimaout/wasapp/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

type profileResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (rt *_router) getMyProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userName, err := rt.db.GetUserNameById(ctx.UserID)
	if err != nil {
		ctx.Logger.WithError(err).Error("error getting user name")
		rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(profileResponse{
		Id:   ctx.UserID,
		Name: userName,
	})
}
