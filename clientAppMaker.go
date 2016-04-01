package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

func (a *App) GetClientVarsRoutes(t *ClientAppSettings) (fileName, JSCode string) {
	fileName = path.Join(t.directories["app"], t.varsRoutesJSFileName)
	s := a.GetServerSettings()
	routing := ""
	for _, mods := range a.Models {
		routing += mods.GetRoutes(mods.GetClientSettings())
	}

	JSCode = fmt.Sprintf(
		`app.constant('appName','%s');
		app.constant('appVersion','%s');
		app.constant('compName','%s');
		app.constant('apiPath','..%s');
		app.value('appVars', {
			user : {
					sessionId:'',
					memberId:'',
					userName:'',
					fbToken:'',
					gpToken:'',
					address:'',
					isNewSignUp : false
				},
			fbAppId : _fbAppId,
			gpClientId : _gpClientId,
			fbSDKLoaded: false,
			gpSDKLoaded: false,
			fbSDKLoadedHndlr: null,
			gpSDKLoadedHndlr: null
		});

		//Routes for Application
		app.config(['$routeProvider', function($routeProvider) {	
			$routeProvider
			.when('/home', {
			   templateUrl: 'app/home/home.view.htm',
			   controller: 'homeController',
			   title : 'Welcome'
			})
			.when('/login', {
			   templateUrl: 'app/login/login.view.htm',
			   controller: 'loginController',
			   title : 'Login'
			})
			%s
			.otherwise({
			   redirectTo: '/home'
			});
		}]);`, a.Name, a.VersionNo, a.CompanyName, s.apiPath, routing)
	return
}

func (a *App) GetClientNavScriptLinks(t *ClientAppSettings) (fileName, content string) {
	fileName = path.Join(t.directories["include templates"], t.indexTemplateName)

	navLinks, scriptLinks := fmt.Sprintln(), fmt.Sprintln()
	for _, mod := range a.Models {
		mt := mod.GetClientSettings()
		navLinks += fmt.Sprintf(`<li><a href="#%s">%s</a></li>`, mt.indexRoute, mod.DisplayName) + fmt.Sprintln()
		scriptLinks += fmt.Sprintf(`<script src = "app/%s/%s"></script>`, mod.Name, mt.indexControllerFileName) +
			fmt.Sprintln() +
			fmt.Sprintf(`<script src = "app/%s/%s"></script>`, mod.Name, mt.controllerFileName) +
			fmt.Sprintln()
	}
	navLinks = fmt.Sprintf(`{{ define "navLinks" }}%s{{ end }}`, navLinks)
	scriptLinks = fmt.Sprintf(`{{ define "scriptLinks" }}%s{{ end }}`, scriptLinks)
	content = navLinks + fmt.Sprintln() + fmt.Sprintln() + scriptLinks
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
	filepath.Walk(app.AppDir, app.fileClientWalker)
}

func (app *App) fileClientWalker(path string, f os.FileInfo, err error) error {
	cmdTxt := app.GetClientSettings().jsBeautifierCmd
	//we are going to search for js files
	if !f.IsDir() {
		if filepath.Ext(path) == ".js" || filepath.Ext(path) == ".json" { //got a .js file
			cmd := exec.Command(cmdTxt, fmt.Sprintf("--outfile=%s", filepath.Base(path)), filepath.Base(path))
			cmd.Dir = filepath.Dir(path)
			Check(cmd.Run())
			fmt.Println("Beautified", path)
		}
	}
	//fmt.Println(path, f)
	return nil
}
