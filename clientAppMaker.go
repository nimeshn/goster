package main

import (
	"fmt"
	"path"
)

func (a *App) GetClientVarsRoutes(t *ClientAppSettings) (fileName, JSCode string) {
	fileName = path.Join(t.directories["app"], t.varsRoutesJSFileName)

	routing := ""
	for _, mods := range a.Models {
		routing += mods.GetRoutes(mods.GetClientSettings())
	}

	JSCode = fmt.Sprintf(
		`app.constant('appName','%s');
		app.constant('appVersion','%s');
		app.constant('compName','%s');
		app.constant('apiPath','../server');
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
		}]);`, a.Name, a.VersionNo, a.CompanyName, routing)
	return
}

func (a *App) GetClientNavScriptLinks(t *ClientAppSettings) (fileName, content string) {
	fileName = path.Join(t.directories["include templates"], t.indexTemplateName)

	navLinks, scriptLinks := fmt.Sprintln(), fmt.Sprintln()
	for _, mod := range a.Models {
		mt := mod.GetClientSettings()
		navLinks += fmt.Sprintf(`<li><a href="%s">%s</a></li>`, mt.indexRoute, mod.DisplayName) + fmt.Sprintln()
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
