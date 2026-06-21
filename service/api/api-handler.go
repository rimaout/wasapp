package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// Login
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	// Authenticated routes (require Bearer token)
	rt.router.GET("/profile", rt.wrapAuthenticated(rt.getMyProfile))
	rt.router.PATCH("/profile/name", rt.wrapAuthenticated(rt.setMyUserName))
	rt.router.PUT("/profile/image", rt.wrapAuthenticated(rt.setMyPhoto))

	// Public image serving (no auth needed)
	rt.router.GET("/users/:userId/profile-image", rt.wrap(rt.getUserProfileImage))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
