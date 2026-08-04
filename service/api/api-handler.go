package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/liveness", rt.liveness)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// Login
	rt.router.POST("/session", rt.wrap(rt.doLogin))

	// User search (no auth required)
	rt.router.GET("/users", rt.wrap(rt.userSearch))

	// User avatar serving (no auth required)
	rt.router.GET("/users/:userId/avatar", rt.wrap(rt.getUserAvatar))

	// --- Authenticated routes (require Bearer token)

	// User
	rt.router.PATCH("/me/name", rt.wrapAuthenticated(rt.setMyUserName))
	rt.router.PUT("/me/avatar", rt.wrapAuthenticated(rt.setMyUserAvatar))

	// Chat creation
	rt.router.POST("/user/:userId/chats", rt.wrapAuthenticated(rt.createDirectChat))
	rt.router.POST("/chats", rt.wrapAuthenticated(rt.createGroupChat))

	// Chat edits
	rt.router.PATCH("/chats/:chatId/name", rt.wrapAuthenticated(rt.setGroupChatName))
	rt.router.PUT("/chats/:chatId/avatar", rt.wrapAuthenticated(rt.setGroupChatAvatar))

	// Chat Image Serving
	rt.router.GET("/chats/:chatId/avatar", rt.wrapAuthenticated(rt.getChatAvatar))

	// Members
	rt.router.GET("/chats/:chatId/members", rt.wrapAuthenticated(rt.getChatMembers))
	rt.router.POST("/chats/:chatId/members", rt.wrapAuthenticated(rt.addMemberToGroup))
	rt.router.DELETE("/chats/:chatId/members/me", rt.wrapAuthenticated(rt.leaveGroup))

	// Messages
	rt.router.POST("/chats/:chatId/messages", rt.wrapAuthenticated(rt.sendMessage))

	return rt.router
}
