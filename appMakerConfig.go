package main

import (
	"path"
)

var clientSettings map[string]string = map[string]string{
	"configFile":   "goster.json",
	"templatesDir": "templates",
}

func GetAppSettings() map[string]string {
	return clientSettings
}

type ClientAppSettings struct {
	helperJSFileName     string
	moduleJSFileName     string
	varsRoutesJSFileName string
	baseTemplateName     string
	indexTemplateName    string
	clientHTMLFile       string
	errorHandler         string
	appTemplateSrcPath   string
	jsBeautifierCmd      string
	directories          map[string]string
}

func (a *App) GetClientSettings() *ClientAppSettings {
	return &ClientAppSettings{
		helperJSFileName:     "app.helper.js",
		moduleJSFileName:     "app.module.js",
		varsRoutesJSFileName: "app.vars.routes.js",
		baseTemplateName:     "base.tmpl",
		indexTemplateName:    "index.tmpl",
		clientHTMLFile:       "index.html",
		errorHandler:         "errorhandler.view.html",
		appTemplateSrcPath:   "appTemplate",
		jsBeautifierCmd:      "jsbeautifier-go",
		directories: map[string]string{
			"client":            path.Join(a.AppDir, "client"),
			"app":               path.Join(a.AppDir, "client/app"),
			"errorHandler":      path.Join(a.AppDir, "client/app/error"),
			"css":               path.Join(a.AppDir, "client/assets/css"),
			"images":            path.Join(a.AppDir, "client/assets/images"),
			"scripts":           path.Join(a.AppDir, "client/scripts"),
			"templates":         path.Join(a.AppDir, "client/templates"),
			"layout templates":  path.Join(a.AppDir, "client/templates/layouts"),
			"include templates": path.Join(a.AppDir, "client/templates/includes"),
		},
	}
}
