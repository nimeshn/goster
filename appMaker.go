package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

var app *App

func GetApp2() *App {
	return app
}

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func CopyFile(srcPath, destPath string) {
	buffStr, err := ioutil.ReadFile(srcPath)
	Check(err)
	Check(ioutil.WriteFile(destPath, buffStr, 0644))
}

func CreateFile(fileName, content string) {
	Check(os.MkdirAll(filepath.Dir(fileName), os.ModeDir))
	//
	Check(ioutil.WriteFile(fileName, []byte(content), 0644))
}

func SaveAppSettings(app *App) {
	settings := GetAppSettings()
	//create goster config file
	appJson, err := json.Marshal(app)
	Check(err)
	Check(ioutil.WriteFile(path.Join(app.AppDir, settings["configFile"]), []byte(appJson), 0644))
}

func CreateNewApp(name, displayName, companyName, versionNo, appDir string) (app *App) {
	app = &App{
		Name:        name,
		DisplayName: displayName,
		CompanyName: companyName,
		VersionNo:   versionNo,
		AppDir:      appDir,
		Models:      []*Model{},
	}
	clientSettings := app.GetClientSettings()
	//create app directories
	fmt.Println("Creating App folder for ", app.Name)
	for name, dir := range clientSettings.directories {
		Check(os.MkdirAll(dir, os.ModeDir))
		fmt.Println("Created ", name, " folder:", dir)
	}
	SaveAppSettings(app)
	//copy app helper js file
	CopyFile(path.Join(clientSettings.appTemplateSrcPath, clientSettings.helperJSFileName),
		path.Join(clientSettings.directories["app"], clientSettings.helperJSFileName))
	//copy app module js file
	CopyFile(path.Join(clientSettings.appTemplateSrcPath, clientSettings.moduleJSFileName),
		path.Join(clientSettings.directories["app"], clientSettings.moduleJSFileName))
	//copy app index template & base template
	CopyFile(path.Join(clientSettings.appTemplateSrcPath, clientSettings.baseTemplateName),
		path.Join(clientSettings.directories["layout templates"], clientSettings.baseTemplateName))
	CopyFile(path.Join(clientSettings.appTemplateSrcPath, clientSettings.indexTemplateName),
		path.Join(clientSettings.directories["include templates"], clientSettings.indexTemplateName))
	//prep the base.tmpl and index.tmpl to prepare index.html
	InitTemplates(clientSettings.directories["templates"])
	return
}

func LoadApp(appDir string) (app *App) {
	buffstr, err := ioutil.ReadFile(path.Join(app.AppDir, GetAppSettings()["configFile"]))
	Check(err)
	app = &App{}
	json.Unmarshal([]byte(buffstr), app)
	return
}

func (app *App) MakeClient() {
	t := app.GetClientSettings()
	//create app.var.routes.js
	fileName, content := app.GetClientVarsRoutes(t)
	CreateFile(fileName, content)
	//create index.tmpl
	fileName, content = app.GetClientNavScriptLinks(t)
	CreateFile(fileName, content)
	//create index.html from the templates
	RenderTemplateToFile(path.Join(t.directories["app"], t.clientHTMLFile), "base.tmpl", nil)
	//
	for _, mods := range app.Models {
		a := mods.GetClientSettings()
		fileName, content = mods.GetClientController(a)
		CreateFile(fileName, content)
		//
		fileName, content = mods.GetClientIndexController(a)
		CreateFile(fileName, content)
		//
		fileName, content = mods.GetClientShowView(a)
		CreateFile(fileName, content)
		//
		fileName, content = mods.GetClientEditView(a)
		CreateFile(fileName, content)
		//
		fileName, content = mods.GetClientIndexView(a)
		CreateFile(fileName, content)
	}
	SaveAppSettings(app)
}
