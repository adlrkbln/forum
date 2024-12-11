package app

import (
	"html/template"
	"log"
)

type Application struct {
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	TemplateCache map[string]*template.Template
}

func New(infoLog, errorLog *log.Logger, templateCache map[string]*template.Template) *Application {
	return &Application{
		ErrorLog:      errorLog,
		InfoLog:       infoLog,
		TemplateCache: templateCache,
	}
}
