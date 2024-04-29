package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)
	rt.router.GET("/bans", rt.wrap(handleGetBannedUsers))
	rt.router.GET("/users/:username", rt.wrap(HandleGetUserProfile))
	rt.router.GET("/photos", rt.wrap(handleGetPhotos))
	rt.router.GET("/users", rt.wrap(HandleGetAllUsers))
	rt.router.GET("/photos/:photoId/comments", rt.wrap(handleGetComments))
	rt.router.GET("/stream", rt.wrap(handleGetMyStream))
	rt.router.GET("/users/:username/followers", rt.wrap(handleGetFollowers))
	rt.router.POST("/users", rt.wrap(HandleAddUser))
	rt.router.POST("/photos/:photoId/comments", rt.wrap(handleCommentPhoto))
	rt.router.POST("/photos/:photoId/likes", rt.wrap(HandleLikePhoto))
	rt.router.POST("/users/:username/bans", rt.wrap(handleBanUser))
	rt.router.POST("/users/:username/follows", rt.wrap(HandleFollowUser))
	rt.router.POST("/session", rt.wrap(doLogin))
	rt.router.POST("/photos", rt.wrap(handleUploadPhoto))
	rt.router.PATCH("/users/:username", rt.wrap(HandleSetUsername))
	rt.router.DELETE("/photos/:photoId", rt.wrap(handleDeletePhoto))
	rt.router.DELETE("/photos/:photoId/likes", rt.wrap(HandleUnlikePhoto))
	rt.router.DELETE("/users/:username/bans", rt.wrap(handleUnbanUser))
	rt.router.DELETE("/users/:username/follows", rt.wrap(HandleUnfollowUser))
	rt.router.DELETE("/comments/:commentId", rt.wrap(handleUncommentPhoto))

	return rt.router
}
