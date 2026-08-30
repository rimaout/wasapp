package api

import (
	"github.com/rimaout/wasapp/service/api/reqcontext"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

// httpRouterHandler is the signature for functions that accepts a reqcontext.RequestContext in addition to those
// required by the httprouter package.
type httpRouterHandler func(http.ResponseWriter, *http.Request, httprouter.Params, reqcontext.RequestContext)

// wrap parses the request and adds a reqcontext.RequestContext instance related to the request.
func (rt *_router) wrap(fn httpRouterHandler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		// Generate a new request UUID for logging (used to trace logs for this request)
		reqUUID, err := uuid.NewV4()
		if err != nil {
			rt.baseLogger.WithError(err).Error("can't generate a request UUID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Create a new request context with the generated UUID
		var ctx = reqcontext.RequestContext{
			ReqUUID: reqUUID,
		}

		// Add the request UUID and remote IP (ip of the client) to the logger for this request
		ctx.Logger = rt.baseLogger.WithFields(logrus.Fields{
			"reqid":     ctx.ReqUUID.String(),
			"remote-ip": r.RemoteAddr,
		})

		// Call the next handler in chain (usually, the handler function for the path)
		fn(w, r, ps, ctx)
	}
}

// wrapAuthenticated creates a request context and validates the Bearer token
// before calling the handler. Returns 401 if the token is missing, expired, or invalid.
func (rt *_router) wrapAuthenticated(fn httpRouterHandler) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		// Generate a new request UUID for logging (used to trace logs for this request)
		reqUUID, err := uuid.NewV4()
		if err != nil {
			rt.baseLogger.WithError(err).Error("can't generate a request UUID")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Create a new request context with the generated UUID
		ctx := reqcontext.RequestContext{
			ReqUUID: reqUUID,
		}

		// Add the request UUID and remote IP (ip of the client) to the logger for this request
		ctx.Logger = rt.baseLogger.WithFields(logrus.Fields{
			"reqid":     ctx.ReqUUID.String(),
			"remote-ip": r.RemoteAddr,
		})

		// Extract the Bearer token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.Logger.Warning("missing or invalid Authorization header")
			rt.respondWithError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "authentication required") //401
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ") // Extract the token value after "Bearer "
		if token == "" {
			ctx.Logger.Warning("empty token in Authorization header")
			rt.respondWithError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "authentication required") //401
			return
		}

		// Validate the token and retrieve the associated user ID
		userId, err := rt.db.GetUserIDByToken(token)
		if err != nil {
			ctx.Logger.WithError(err).Error("error validating token")
			rt.respondWithError(w, http.StatusInternalServerError, ErrCodeInternalError, "internal server error") //500
			return
		}

		if userId == "" {
			ctx.Logger.Warning("invalid or expired token")
			rt.respondWithError(w, http.StatusUnauthorized, ErrCodeUnauthorized, "authentication failed") //401
			return
		}

		// Set the authenticated user ID in the request context
		ctx.UserID = userId

		// Call the next handler in chain (usually, the handler function for the path)
		fn(w, r, ps, ctx)
	}
}
