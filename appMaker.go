package main

import (
	"encoding/json"
	"fmt"
	"github.com/yosssi/gohtml"
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
	fileMap := map[string]string{
		clientSettings.helperJSFileName:      clientSettings.directories["app"],
		clientSettings.moduleJSFileName:      clientSettings.directories["app"],
		clientSettings.logoFile:              clientSettings.directories["images"],
		clientSettings.bootstrapSocialFile:   clientSettings.directories["css"],
		clientSettings.baseTemplateName:      clientSettings.directories["layout templates"],
		clientSettings.indexTemplateName:     clientSettings.directories["include templates"],
		clientSettings.errorHandler:          clientSettings.directories["errorHandler"],
		clientSettings.homeHTMLFile:          clientSettings.directories["home"],
		clientSettings.homeControllerFile:    clientSettings.directories["home"],
		clientSettings.loginHTMLFile:         clientSettings.directories["login"],
		clientSettings.loginControllerFile:   clientSettings.directories["login"],
		clientSettings.profileHTMLFile:       clientSettings.directories["profile"],
		clientSettings.profileControllerFile: clientSettings.directories["profile"],
	}
	for filename, dir := range fileMap {
		//copy app helper js file
		CopyFile(path.Join(clientSettings.appTemplateSrcPath, filename),
			path.Join(dir, filename))
	}
	user := &Model{
		Name:        "user",
		DisplayName: "User",
		ViewType:    List,
		Fields: []*Field{
			&Field{
				Name:        "fn",
				DisplayName: "First Name",
				Type:        String,
				Validator:   &FieldValidation{Required: true, IsAlpha: true, MaxLen: 25},
			},
			&Field{
				Name:        "ln",
				DisplayName: "Last Name",
				Type:        String,
				Validator:   &FieldValidation{Required: true, IsAlpha: true, MaxLen: 25},
			},
			&Field{
				Name:        "emlid",
				DisplayName: "EmailId",
				Type:        String,
				Validator:   &FieldValidation{Required: true, Email: true, MaxLen: 255},
			},
			&Field{
				Name:        "sex",
				DisplayName: "Gender",
				Type:        String,
				Validator:   &FieldValidation{Required: true, IsAlpha: true, MaxLen: 1},
			},
			&Field{
				Name:       "fbuid",
				Type:       String,
				HideInEdit: true,
				Validator:  &FieldValidation{MaxLen: 500},
			},
			&Field{
				Name:       "gpuid",
				Type:       String,
				HideInEdit: true,
				Validator:  &FieldValidation{MaxLen: 500},
			},
			&Field{
				Name:        "active",
				DisplayName: "Active",
				Type:        Boolean,
				Validator:   &FieldValidation{},
			},
		},
	}
	app.AddModel(user)
	return
}

func LoadApp(appDir string) (app *App) {
	buffstr, err := ioutil.ReadFile(path.Join(app.AppDir, GetAppSettings()["configFile"]))
	Check(err)
	app = &App{}
	json.Unmarshal([]byte(buffstr), app)
	return
}
