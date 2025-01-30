package main

import (
	"crypto/tls"
	"forum/internal/app"
	"forum/internal/handlers"
	"forum/internal/repo"
	"forum/internal/service"
	"forum/conf"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dsn := "./forum.db"
	configPath := "./config.json"

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := repo.NewDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	
	config, err := conf.Load(configPath)
	if err != nil {
		log.Fatal(err)
	}
	templateCache, err := app.NewTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := app.New(infoLog, errorLog, templateCache, config.GoogleConfig, config.GithubConfig)
	service := service.NewService(db)
	handlers := handlers.New(app, service)

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences: []tls.CurveID{
			tls.CurveP256,
			tls.CurveP384,
		},
	}

	srv := &http.Server{
		Addr:         ":8080",
		ErrorLog:     errorLog,
		Handler:      handlers.RateLimiter(handlers.Routes()),
		TLSConfig:    tlsConfig,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	infoLog.Print("Starting server on https://localhost:8080")
	err = srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem")
	errorLog.Fatal(err)
}
