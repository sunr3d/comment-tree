package httphandlers

import (
	"github.com/wb-go/wbf/ginext"

	"github.com/sunr3d/comment-tree/internal/interfaces/services"
)

type Handler struct {
	svc services.CommentTree
}

func New(svc services.CommentTree) *Handler {
	return &Handler{
		svc: svc,
	}
}

func (h *Handler) RegisterHandlers() *ginext.Engine {
	router := ginext.New()
	router.Use(ginext.Logger(), ginext.Recovery())

	// Web-интерфейс
	router.Static("/static", "./static")
	router.GET("/", func(c *ginext.Context) {
		c.File("./static/index.html")
	})

	// API
	router.POST("/comments", h.writeComment)
	router.GET("/comments", h.getComments)
	router.DELETE("/comments/:id", h.deleteComment)

	return router
}
