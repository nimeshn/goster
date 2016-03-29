package main

import (
	"encoding/json"
	"fmt"
	"github.com/yosssi/gohtml"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
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

func FormatHTMLFile(srcPath string) {
	buffStr, err := ioutil.ReadFile(srcPath)
	Check(err)
	Check(ioutil.WriteFile(srcPath, []byte(gohtml.Format(string(buffStr))), 0644))
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
	fmt.Println("Creating App directories for", app.Name)
	for name, dir := range clientSettings.directories {
		Check(os.MkdirAll(dir, os.ModeDir))
		fmt.Println(name, "folder:", dir)
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
	//copy Error
	CopyFile(path.Join(clientSettings.appTemplateSrcPath, clientSettings.errorHandler),
		path.Join(clientSettings.directories["errorHandler"], clientSettings.errorHandler))
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
	//prep the base.tmpl and index.tmpl to prepare index.html
	InitTemplates(t.directories["templates"])
	//create index.html from the templates
	Check(RenderTemplateToFile(path.Join(t.directories["client"], t.clientHTMLFile), "base.tmpl", map[string]interface{}{"dummy": "dummy Data"}))
	FormatHTMLFile(path.Join(t.directories["client"], t.clientHTMLFile))
	//
	for _, mods := range app.Models {
		if mods.DisplayName == "" {
			mods.DisplayName = mods.Name
		}
		mods.DisplayName = strings.Title(mods.DisplayName)
		for _, flds := range mods.Fields {
			if flds.DisplayName == "" {
				flds.DisplayName = flds.Name
			}
			flds.DisplayName = strings.Title(flds.DisplayName)
		}

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
	//
	filepath.Walk(app.AppDir, app.fileWalker)
}

func (app *App) fileWalker(path string, f os.FileInfo, err error) error {
	cmdTxt := app.GetClientSettings().jsBeautifierCmd
	//we are going to search for js files
	if !f.IsDir() {
		if filepath.Ext(path) == ".js" { //got a .js file
			cmd := exec.Command(cmdTxt, fmt.Sprintf("--outfile=%s", filepath.Base(path)), filepath.Base(path))
			cmd.Dir = filepath.Dir(path)
			Check(cmd.Run())
			fmt.Println("Beautified", path)
		}
	}
	//fmt.Println(path, f)
	return nil
}
