package server

import (
	"context"
	"errors"
	"net"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/marzusia/poeia"
	log "github.com/sirupsen/logrus"
)

type Settings struct {
	Addr         string
	TemplateDir  string
	SiteName     string
	SiteLanguage string
	StaticPath   string
	StaticDir    string
	HostStatic   bool
}

type Server struct {
	server   *http.Server
	listener net.Listener
	router   *gin.Engine
	db       *poeia.DB

	settings Settings
}

func NewServer(settings Settings) *Server {
	// Create server object
	server := &Server{
		server:   &http.Server{},
		router:   gin.Default(),
		settings: settings,
	}

	// Initialise router
	server.router.Use(server.authenticate())
	server.registerPostRoutes()

	// Host static files if needed
	if server.settings.HostStatic && server.settings.StaticPath != "" && server.settings.StaticDir != "" {
		server.router.Static(server.settings.StaticPath, server.settings.StaticDir)
	}

	// Set handler and return server
	server.server.Handler = server.router
	return server
}

func (s *Server) WithDB(db *poeia.DB) *Server {
	s.db = db
	return s
}

func (s *Server) Open() (err error) {
	// Ensure database has been set
	if s.db == nil {
		return errors.New("database has not been set")
	}

	// Ensure address is set
	if s.settings.Addr == "" {
		log.Info("⚠️ no address was set, defaulting to :8000")
		s.settings.Addr = ":8000"
	}

	// Load templates
	if s.settings.TemplateDir == "" {
		return errors.New("template directory has not been set")
	}
	s.router.LoadHTMLGlob(filepath.Join(s.settings.TemplateDir, "**/*"))

	// Listen at given address in new goroutine
	if s.listener, err = net.Listen("tcp", s.settings.Addr); err != nil {
		return err
	}
	go s.server.Serve(s.listener)
	return nil
}

func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func (s *Server) Addr() string {
	return s.settings.Addr
}

func (s *Server) authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: check authentication
		c.Next()
	}
}

func (s *Server) defaultContextWith(h gin.H) gin.H {
	defaultContext := gin.H{
		"siteName":     s.settings.SiteName,
		"siteLanguage": s.settings.SiteLanguage,
		"staticPath":   s.settings.StaticPath,
	}
	for k, v := range h {
		defaultContext[k] = v
	}
	return defaultContext
}
