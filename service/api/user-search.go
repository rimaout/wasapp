package api

import (
	"encoding/json"
	"net/http"

	"github.com/rimaout/wasapp/service/api/reqcontext"
	"github.com/rimaout/wasapp/service/database"
	"github.com/julienschmidt/httprouter"
)

type userSearchResponse struct {
	UsersList []database.User `json:"usersList"`
}

func (rt *_router) userSearch(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	query := r.URL.Query().Get("username")

	users, err := rt.db.SearchUsers(query)
	if err != nil {
		ctx.Logger.WithError(err).Error("error searching users")
		rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userSearchResponse{UsersList: users})
}
