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

	// User search (no auth required)
	rt.router.GET("/users", rt.wrap(rt.userSearch))

	// Authenticated routes (require Bearer token)

	// User
	rt.router.PATCH("/me/name", rt.wrapAuthenticated(rt.setMyUserName))
	rt.router.PUT("/me/avatar", rt.wrapAuthenticated(rt.setMyUserAvatar))

	// Chat creation
	rt.router.POST("/user/:userId/chats", rt.wrapAuthenticated(rt.createDirectChat))
	rt.router.POST("/chats/groups", rt.wrapAuthenticated(rt.createGroupChat))

	// Chat edits
	rt.router.PUT("/chats/:chatId/name:", rt.wrapAuthenticated(rt.setGroupChatName))
	rt.router.PUT("/chats/:chatId/avatar:", rt.wrapAuthenticated(rt.setGroupChatAvatar))

	// Public image serving (no auth needed)
	rt.router.GET("/users/:userId/avatar", rt.wrap(rt.getUserAvatar))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
