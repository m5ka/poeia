package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/marzusia/poeia"
	"github.com/marzusia/poeia/server"
	log "github.com/sirupsen/logrus"
)

func RunServer() error {
	log.Info("🌼 poeia ─ no-js microblogging")

	if err := godotenv.Load(); err != nil {
		return errors.New("could not load environment variables from .env")
	}

	db, err := poeia.Open(os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_USERNAME"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_DATABASE"))
	if err != nil {
		return err
	}

	settings := server.Settings{
		Addr:         os.Getenv("POEIA_ADDR"),
		TemplateDir:  os.Getenv("POEIA_TEMPLATE_DIR"),
		SiteName:     os.Getenv("POEIA_SITE_NAME"),
		SiteLanguage: os.Getenv("POEIA_SITE_LANGUAGE"),
		StaticPath:   os.Getenv("POEIA_STATIC_PATH"),
		StaticDir:    os.Getenv("POEIA_STATIC_DIR"),
		HostStatic:   os.Getenv("POEIA_HOST_STATIC") == "true",
	}
	server := server.NewServer(settings).WithDB(db)
	if err = server.Open(); err != nil {
		return err
	}

	log.Infof("🌍 listening on %s", server.Addr())
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	<-ch

	log.Info("👋 shutting down...")
	server.Close()
	return nil
}

func main() {
	if err := RunServer(); err != nil {
		log.Fatal(fmt.Sprintf("💀 %v", err))
	}
}
