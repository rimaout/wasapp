package api

import (
	"net/http"
	"strings"

	"github.com/rimaout/wasapp/service/api/reqcontext"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// wrapAuthenticated creates a request context and validates the Bearer token
// before calling the handler. Returns 401 if the token is missing, expired, or invalid.
func (rt *_router) wrapAuthenticated(fn httpRouterHandler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		reqUUID, err := uuid.NewV4()
		if err != nil {
			rt.baseLogger.WithError(err).Error("can't generate a request UUID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx := reqcontext.RequestContext{
			ReqUUID: reqUUID,
		}

		ctx.Logger = rt.baseLogger.WithFields(logrus.Fields{
			"reqid":     ctx.ReqUUID.String(),
			"remote-ip": r.RemoteAddr,
		})

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.Logger.Warning("missing or invalid Authorization header")
			rt.respondWithError(w, http.StatusUnauthorized, "401", "authentication required")
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			ctx.Logger.Warning("empty token in Authorization header")
			rt.respondWithError(w, http.StatusUnauthorized, "401", "authentication required")
			return
		}

		userId, err := rt.db.GetUserIDByToken(token)
		if err != nil {
			ctx.Logger.WithError(err).Error("error validating token")
			rt.respondWithError(w, http.StatusInternalServerError, "500", "internal server error")
			return
		}

		if userId == "" {
			ctx.Logger.Warning("invalid or expired token")
			rt.respondWithError(w, http.StatusUnauthorized, "401", "authentication failed")
			return
		}

		ctx.UserID = userId

		fn(w, r, ps, ctx)
	}
}
