package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/yosssi/gohtml"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	dirPerm, filePerm os.FileMode = 0755, 0644
	app               *App
	camelingRegex     = regexp.MustCompile("[0-9A-Za-z]+")
)

func CamelCase(src string) string {
	byteSrc := []byte(src)
	chunks := camelingRegex.FindAll(byteSrc, -1)
	for idx, val := range chunks {
		if idx > 0 {
			chunks[idx] = bytes.Title(val)
		} else {
			valStr := string(val)
			chunks[idx] = []byte(strings.ToLower(valStr[:1]) + string(valStr[1:]))
		}
	}
	return string(bytes.Join(chunks, nil))
}

func GetApp() *App {
	return app
}

func Check(err error) {
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func CopyFile(srcPath, destPath string) {
	buffStr, err := ioutil.ReadFile(srcPath)
	Check(err)
	Check(ioutil.WriteFile(destPath, buffStr, filePerm))
	fmt.Println("Copied", destPath)
}

func FormatHTMLFile(srcPath string) {
	buffStr, err := ioutil.ReadFile(srcPath)
	Check(err)
	Check(ioutil.WriteFile(srcPath, []byte(gohtml.Format(string(buffStr))), filePerm))
	fmt.Println("Formatted", srcPath)
}

func CreateFile(fileName, content string) {
	Check(os.MkdirAll(filepath.Dir(fileName), dirPerm))
	Check(ioutil.WriteFile(fileName, []byte(content), filePerm))
	fmt.Println("Created", fileName)
}

func SaveAppSettings(app *App) {
	settings := GetAppSettings()
	//create goster config file
	appJson, err := json.Marshal(app)
	Check(err)
	Check(ioutil.WriteFile(path.Join(app.AppDir, settings["configFile"]), []byte(appJson), 0644))
}

func CreateNewApp(name, displayName, companyName, versionNo string, portNumber int) (app *App) {
	workDir, err := os.Getwd()
	Check(err)
	fmt.Println("CamelCasing Appname:", CamelCase(name))
	app = &App{
		Name:        name,
		DisplayName: displayName,
		CompanyName: companyName,
		VersionNo:   versionNo,
		AppDir:      path.Join(path.Dir(workDir), CamelCase(name)),
		Models:      []*Model{},
		PortNumber:  portNumber,
	}
	clientSettings := app.GetClientSettings()
	serverSettings := app.GetServerSettings()
	//create app directories
	fmt.Println("Creating App directories for", app.Name)
	for _, dir := range clientSettings.directories {
		Check(os.MkdirAll(dir, dirPerm))
		fmt.Println("Created", dir)
	}
	fileMap := map[string]string{
		clientSettings.helperJSFileName:        clientSettings.directories["app"],
		clientSettings.moduleJSFileName:        clientSettings.directories["app"],
		clientSettings.logoFile:                clientSettings.directories["images"],
		clientSettings.bootstrapSocialFile:     clientSettings.directories["css"],
		clientSettings.baseTemplateName:        clientSettings.directories["layout templates"],
		clientSettings.indexTemplateName:       clientSettings.directories["include templates"],
		clientSettings.errorHandler:            clientSettings.directories["errorHandler"],
		clientSettings.homeHTMLFile:            clientSettings.directories["home"],
		clientSettings.homeControllerFile:      clientSettings.directories["home"],
		clientSettings.loginHTMLFile:           clientSettings.directories["login"],
		clientSettings.loginControllerFile:     clientSettings.directories["login"],
		clientSettings.profileHTMLFile:         clientSettings.directories["profile"],
		clientSettings.profileControllerFile:   clientSettings.directories["profile"],
		serverSettings.mainGoFileName:          serverSettings.directories["server"],
		serverSettings.helperGoFileName:        serverSettings.directories["server"],
		serverSettings.loginControllerFileName: serverSettings.directories["server"],
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
				Unique:      true,
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
			&Field{
				Name:        "createdAt",
				DisplayName: "Created At",
				Type:        Date,
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
