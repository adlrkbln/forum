package main

import (
	"forum/internal/app"
	"forum/internal/handlers"
	"forum/internal/repo"
	"forum/internal/service"
	"log"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dsn := "./forum.db"

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := repo.NewDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	templateCache, err := app.NewTemplateCache()
	if err != nil {
		errorLog.Fatal(err)
	}

	app := app.New(infoLog, errorLog, templateCache)
	service := service.NewService(db)
	handlers := handlers.New(app, service)

	srv := &http.Server{
		Addr:     ":8080",
		ErrorLog: errorLog,
		Handler:  handlers.Routes(),
	}

	infoLog.Print("Starting server on http://localhost:8080")
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
