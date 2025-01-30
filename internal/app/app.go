package app

import (
	"forum/conf"
	"html/template"
	"log"
)

type Application struct {
	ErrorLog      *log.Logger
	InfoLog       *log.Logger
	TemplateCache map[string]*template.Template
	GoogleConfig  conf.GoogleConfig
	GithubConfig  conf.GithubConfig
}

func New(infoLog, errorLog *log.Logger, templateCache map[string]*template.Template, googleConfig conf.GoogleConfig, githubConfig conf.GithubConfig) *Application {
	return &Application{
		ErrorLog:      errorLog,
		InfoLog:       infoLog,
		TemplateCache: templateCache,
		GoogleConfig:  googleConfig,
		GithubConfig:  githubConfig,
	}
}
