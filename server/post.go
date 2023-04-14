package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) registerPostRoutes() {
	s.router.GET("/", s.handleListPosts)
}

func (s *Server) handleListPosts(c *gin.Context) {
	c.HTML(http.StatusOK, "post/list", s.defaultContextWith(gin.H{
		"title": "all posts",
	}))
}
