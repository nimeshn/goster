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
	logoFile             string
	bootstrapSocialFile  string
	baseTemplateName     string
	indexTemplateName    string
	clientHTMLFile       string
	homeHTMLFile         string
	homeControllerFile   string
	loginHTMLFile        string
	loginControllerFile  string
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
		logoFile:             "logo.png",
		bootstrapSocialFile:  "bootstrap-social.css",
		baseTemplateName:     "base.tmpl",
		indexTemplateName:    "index.tmpl",
		clientHTMLFile:       "index.html",
		homeHTMLFile:         "home.view.htm",
		homeControllerFile:   "home.controller.js",
		loginHTMLFile:        "login.view.htm",
		loginControllerFile:  "login.controller.js",
		errorHandler:         "errorhandler.view.html",
		appTemplateSrcPath:   "appTemplate",
		jsBeautifierCmd:      "jsbeautifier-go",
		directories: map[string]string{
			"client":            path.Join(a.AppDir, "client"),
			"app":               path.Join(a.AppDir, "client/app"),
			"home":              path.Join(a.AppDir, "client/app/home"),
			"login":             path.Join(a.AppDir, "client/app/login"),
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
