package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func (a *App) GetServerRoutes(t *ServerAppSettings) (fileName, GoCode string) {
	fileName = path.Join(t.directories["server"], t.actionRoutesFileName)

	routing := ""
	for _, mods := range a.Models {
		routing += mods.GetServerRoutes(t) + fmt.Sprintln()
	}

	GoCode = fmt.Sprintf(
		`package main

		import(
			"net/http"
		)

		func MakeActionRoutes() {
			%s
		}`, routing)
	return
}

func (app *App) MakeServer() {
	t := app.GetServerSettings()
	//create actionRoutes.go
	fileName, content := app.GetServerRoutes(t)
	CreateFile(fileName, content)
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

		a := mods.GetServerSettings()
		fileName, content = mods.GetServerModel(a)
		CreateFile(fileName, content)
		//
		fileName, content = mods.GetServerController(a)
		CreateFile(fileName, content)
	}
	//SaveAppSettings(app)
	//
	filepath.Walk(app.AppDir, app.fileServerWalker)
}

func (app *App) fileServerWalker(path string, f os.FileInfo, err error) error {
	cmdTxt := app.GetServerSettings().goBeautifierCmd
	//we are going to search for js files
	if !f.IsDir() {
		if filepath.Ext(path) == ".go" { //got a .go file
			//fmt.Println(cmdTxt, filepath.Dir(path), filepath.Base(path))
			cmd := exec.Command(cmdTxt, "-w", filepath.Base(path))
			cmd.Dir = filepath.Dir(path)
			Check(cmd.Run())
			fmt.Println("Formatted", path)
		}
	}
	return nil
}