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

	// Chats Reading
	rt.router.GET("/chats", rt.wrapAuthenticated(rt.getMyChats))

	// Chat Creation
	rt.router.POST("/user/:userId/chats", rt.wrapAuthenticated(rt.createDirectChat))
	rt.router.POST("/chats", rt.wrapAuthenticated(rt.createGroupChat))

	// Chat Edits
	rt.router.PATCH("/chats/:chatId/name", rt.wrapAuthenticated(rt.setGroupChatName))
	rt.router.PUT("/chats/:chatId/avatar", rt.wrapAuthenticated(rt.setGroupChatAvatar))

	// Members
	rt.router.GET("/chats/:chatId/members", rt.wrapAuthenticated(rt.getChatMembers))
	rt.router.POST("/chats/:chatId/members", rt.wrapAuthenticated(rt.addMemberToGroup))
	rt.router.DELETE("/chats/:chatId/members/me", rt.wrapAuthenticated(rt.leaveGroup))

	// Messages
	rt.router.GET("/chats/:chatId/messages", rt.wrapAuthenticated(rt.getChatMessages))
	rt.router.POST("/chats/:chatId/messages", rt.wrapAuthenticated(rt.sendMessage))
	rt.router.DELETE(`/chats/:chatId/messages/:messageId`, rt.wrapAuthenticated(rt.deleteMessage))
	rt.router.POST("/chats/:chatId/messages/:messageId/reply", rt.wrapAuthenticated(rt.replyMessage))

	// Message Forward
	rt.router.POST("/chats/:chatId/messages/:messageId/forwards", rt.wrapAuthenticated(rt.forwardMessage))

	// Mark Chat as Read
	rt.router.POST("/chats/:chatId/read", rt.wrapAuthenticated(rt.markChatRead))

	// Messages Reactions
	rt.router.POST("/chats/:chatId/messages/:messageId/reactions", rt.wrapAuthenticated(rt.addReactionToMessage))
	rt.router.DELETE("/chats/:chatId/messages/:messageId/reactions", rt.wrapAuthenticated(rt.removeReactionFromMessage))

	// Image Serving
	rt.router.PUT("/me/avatar", rt.wrapAuthenticated(rt.setMyUserAvatar))
	rt.router.GET("/chats/:chatId/avatar", rt.wrapAuthenticated(rt.getChatAvatar))
	rt.router.GET("/chats/:chatId/images/:imageId", rt.wrapAuthenticated(rt.getMessageImage))

	return rt.router
}
