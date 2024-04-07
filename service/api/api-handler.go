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
	rt.router.POST("/users)", rt.wrap(HandleAddUser))
	rt.router.GET("/users/:username", rt.wrap(HandleGetUserProfile))
	rt.router.POST("/users/:username", rt.wrap(HandleSetUsername))
	rt.router.POST("/photos/:photoId/likes", rt.wrap(HandleLikePhoto))
	rt.router.DELETE("/photos/:photoId/likes", rt.wrap(HandleUnlikePhoto))
	rt.router.POST("/users/:username/follows", rt.wrap(HandleFollowUser))
	rt.router.DELETE("/users/:username/follows", rt.wrap(HandleUnfollowUser))

	return rt.router
}
