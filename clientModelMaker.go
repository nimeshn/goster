package main

import (
	"fmt"
	"github.com/yosssi/gohtml"
	"path"
	"strings"
)

type ClientModelSettings struct {
	idCol                   string
	formName                string
	formData                string
	indexData               string
	indexFunc               string
	isNewFunc               string
	newFunc                 string
	loadFunc                string
	saveFunc                string
	deleteFunc              string
	controllerName          string
	indexControllerName     string
	indexViewFileName       string
	editViewFileName        string
	showViewFileName        string
	controllerFileName      string
	indexControllerFileName string
	indexRoute              string
	newRoute                string
	editRoute               string
	showRoute               string
	saveRoute               string
	deleteRoute             string
}

func (m *Model) GetClientSettings() *ClientModelSettings {
	return &ClientModelSettings{
		idCol:                   fmt.Sprintf("%sId", m.Name),
		formName:                fmt.Sprintf("%sForm", m.Name),
		formData:                fmt.Sprintf("%sData", m.Name),
		indexData:               fmt.Sprintf("%sList", m.Name),
		indexFunc:               fmt.Sprintf("Get%sList", strings.Title(m.Name)),
		isNewFunc:               fmt.Sprintf("IsNew%s", strings.Title(m.Name)),
		newFunc:                 fmt.Sprintf("New%s", strings.Title(m.Name)),
		loadFunc:                fmt.Sprintf("Load%s", strings.Title(m.Name)),
		saveFunc:                fmt.Sprintf("Save%s", strings.Title(m.Name)),
		deleteFunc:              fmt.Sprintf("Delete%s", strings.Title(m.Name)),
		controllerName:          fmt.Sprintf("%sController", m.Name),
		indexControllerName:     fmt.Sprintf("%sIndexController", m.Name),
		indexViewFileName:       fmt.Sprintf("%s.index.view.htm", m.Name),
		editViewFileName:        fmt.Sprintf("%s.edit.view.htm", m.Name),
		showViewFileName:        fmt.Sprintf("%s.show.view.htm", m.Name),
		controllerFileName:      fmt.Sprintf("%s.controller.js", m.Name),
		indexControllerFileName: fmt.Sprintf("%s.index.controller.js", m.Name),
		indexRoute:              fmt.Sprintf("/%s/list", m.Name),
		newRoute:                fmt.Sprintf("/%s/new", m.Name),
		editRoute:               fmt.Sprintf("/%s/edit/", m.Name),
		showRoute:               fmt.Sprintf("/%s/view/", m.Name),
		saveRoute:               fmt.Sprintf("/%s", m.Name),
		deleteRoute:             fmt.Sprintf("/%s/delete", m.Name),
	}
}

func (m *Model) GetClientEditView(a *ClientModelSettings) (fileName, htmlCode string) {
	fileName = path.Join(m.appRef.GetClientSettings().directories["app"], m.Name, a.editViewFileName)

	htmlCode = `<div class="row" ng-include="'/app/error/errorhandler.view.html'"></div><div class="row text-center">`
	htmlCode += fmt.Sprintf(`<h3 ng-bind="(%s()?'New':'Edit') + ' %s'"></h3>`, a.isNewFunc, m.DisplayName)
	htmlCode += `<hr/></div>`

	htmlCode += fmt.Sprintf(`<form class="form-horizontal" role="form" name="%s" ng-submit="%s()"><div class="row">`, a.formName, a.saveFunc)
	for _, fld := range m.Fields {
		htmlCode += `<div class="form-group">`
		htmlCode += fmt.Sprintf(`<label class="control-label col-sm-4" for="%s">%s:</label>`, fld.Name, fld.DisplayName)
		htmlCode += `<div class="col-sm-8">`
		htmlCode += fmt.Sprintf(`<input type="text" class="form-control" id="%s" placeholder="Enter %s" ng-model="%s.%s" title="%s" `,
			fld.Name, fld.DisplayName, a.formData, fld.Name, fld.DisplayName)

		if fld.Type == Boolean {
			htmlCode += `type="checkbox" `
		} else if fld.Type == Date {
			htmlCode += `type="date" `
		} else if fld.Type == Integer {
			htmlCode += `type="number" step="1" `
		} else if fld.Type == Float {
			htmlCode += `type="number" step=".01" `
		} else if fld.Type == String {
			if fld.Validator.Email {
				htmlCode += `type="email" `
			} else if fld.Validator.Url {
				htmlCode += `type="url" `
			} else {
				htmlCode += `type="text" `
			}
		}

		if fld.Validator.MinLen > 0 {
			htmlCode += fmt.Sprintf(` minlength="%s"`, fld.Validator.MinLen)
		}
		if fld.Validator.MaxLen > 0 {
			htmlCode += fmt.Sprintf(` maxlength="%s"`, fld.Validator.MaxLen)
		}
		if fld.Validator.MinValue > 0 {
			htmlCode += fmt.Sprintf(` min="%s"`, fld.Validator.MinValue)
		}
		if fld.Validator.MaxValue > 0 {
			htmlCode += fmt.Sprintf(` max="%s"`, fld.Validator.MaxValue)
		}
		if fld.Validator.Email {
			htmlCode += ` pattern="[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,3}$"`
		}
		if fld.Validator.Url {
			htmlCode += ` pattern="https?://.+"`
		}
		if fld.Validator.IsAlpha {
			htmlCode += ` pattern="[A-Za-z]"`
		}
		if fld.Validator.IsAlphaNumeric {
			htmlCode += ` pattern="[A-Za-z0-9]"`
		}
		if fld.Validator.Required {
			htmlCode += ` required`
		}
		htmlCode += `/>`
		htmlCode += `</div>`
		htmlCode += `</div>`
	}
	htmlCode += `</div></form>`
	htmlCode += `<div class="row"><hr/></div>`
	htmlCode = gohtml.Format(htmlCode)
	return
}

func (m *Model) GetClientShowView(a *ClientModelSettings) (fileName, htmlCode string) {
	fileName = path.Join(m.appRef.GetClientSettings().directories["app"], m.Name, a.showViewFileName)

	htmlCode = fmt.Sprintf(`<div class="row" ng-include="'/app/error/errorhandler.view.html'"></div>`+
		`<div class="row text-center"><h3>%s Details</h3><hr/></div>`, m.DisplayName)

	htmlCode += `<div class="row"><div class="col-sm-12">`
	for _, fld := range m.Fields {
		htmlCode += fmt.Sprintf(`<div class="row"><div class="col-sm-12"><h3>%s</h3><p ng-bind="%s.%s"></p></div></div>`,
			fld.DisplayName, a.formData, fld.Name)
	}
	htmlCode += `</div></div>`
	htmlCode += `<div class="row"><hr/></div>`
	htmlCode = gohtml.Format(htmlCode)
	return
}

func (m *Model) GetClientIndexView(a *ClientModelSettings) (fileName, htmlCode string) {
	fileName = path.Join(m.appRef.GetClientSettings().directories["app"], m.Name, a.indexViewFileName)

	htmlCode = fmt.Sprintf(`<div class="row" ng-include="'/app/error/errorhandler.view.html'"></div>`+
		`<div class="row text-center"><h3>%s List</h3><hr/></div>`, m.DisplayName)

	htmlCode += fmt.Sprintf(`<div class="row text-center"><a href="" ng-click="%s()"><span class="glyphicon glyphicon-refresh"/> Refresh</a>`+
		`<a href="#%s" class="col-sm-offset-1"><span class="glyphicon glyphicon-plus"/> New %s</a></div><br/>`,
		a.loadFunc, a.newRoute, m.DisplayName)

	htmlCode += fmt.Sprintf(`<div ng-if="%s.length==0" class="row"><div class="col-sm-12 text-center"><h3>0 Records Found.</h3></div></div>`,
		a.formData)

	htmlCode += fmt.Sprintf(`<div ng-if="%s.length>0" class="row"><div class="col-sm-12">`, a.formData)
	if m.ViewType == List {
		htmlCode += fmt.Sprintf(`<div class="row" ng-repeat="x in %s | orderBy:createdOn:reverse">`, a.formData)
		htmlCode += fmt.Sprintf(`<div class="col-sm-1"><a href="#%s{{x.id}}" alt="View %s" title="View %s">`+
			`<span class="glyphicon glyphicon-folder-open"></span></a></div>`,
			a.showRoute, m.DisplayName, m.DisplayName)
		htmlCode += fmt.Sprintf(`<div class="col-sm-1"><a href="#%s{{x.id}}" alt="Edit %s" title="Edit %s">`+
			`<span class="glyphicon glyphicon-edit"></span></a></div>`,
			a.editRoute, m.DisplayName, m.DisplayName)
		htmlCode += fmt.Sprintf(`<div class="col-sm-1"><a href="" alt="Delete %s" title="Delete %s"><span class="glyphicon glyphicon-remove" `+
			`ng-click="%s(x.id);"></span></a></div>`,
			m.DisplayName, m.DisplayName, a.deleteFunc)
		for _, fld := range m.Fields {
			if !fld.HideInIndex {
				htmlCode += fmt.Sprintf(`<div class="col-sm-1"><span ng-bind="%s.%s"></span></div>`,
					a.formData, fld.Name)
			}
		}
		htmlCode += `</div>`
	} else if m.ViewType == Table {

	}
	htmlCode += `</div></div>`
	htmlCode += `<div class="row"><hr/></div>`
	htmlCode = gohtml.Format(htmlCode)
	return
}

func (m *Model) GetClientController(a *ClientModelSettings) (fileName, JSCode string) {
	fileName = path.Join(m.appRef.GetClientSettings().directories["app"], m.Name, a.controllerFileName)

	//modelIsNewFunc
	isNewFunc := fmt.Sprintf(
		`//function to check if the view is for new model entity
		$scope.%s = function(){
			return (!$scope.%s || $scope.%s == "" || $scope.%s == null);
		};`, a.isNewFunc, a.idCol, a.idCol, a.idCol)

	//modelNewFunc
	newFunc := fmt.Sprintf(
		`//function to Get New model entity
		$scope.%s =function(){
		$http.get(apiPath + "%s")
			.then(function(response) {
				if (response.status == 200){
					$scope.%s = response.data.%s;
					clearAPIError($scope);
				}
			},
			function(response) {
				handleAPIError($scope, response);
			}
		);
	};`, a.newFunc, a.newRoute, a.formData, m.Name)

	//modelLoadFunc
	loadFunc := fmt.Sprintf(
		`//function to load model entity
		$scope.%s =function(){
		$http.get(apiPath + "/%s/" + $scope.%s)
			.then(function(response) {
				if (response.status == 200){
					$scope.%s = response.data.%s;
					clearAPIError($scope);
				}
			},
			function(response) {
				handleAPIError($scope, response);
			}
		);
	};`, a.loadFunc, m.Name, a.idCol, a.formData, m.Name)

	//modelSaveFunc
	saveFunc := fmt.Sprintf(
		`//function to save model entity
		$scope.%s =function(){
			$http({
					method: $scope.%s()?'POST':'PUT',
					url: apiPath + "%s",
					data: $scope.%s
				}).then(
				function(response) {
					if (response.status == 200){
						clearAPIError($scope);
						$location.path("#%s");
					} 
					else {
					  $scope.message = data.message;
					}
				},
				function(response){
					handleAPIError($scope, response);
			  });
		};`, a.saveFunc, a.isNewFunc, a.saveRoute, a.formData, a.indexRoute)

	JSCode = isNewFunc + fmt.Sprintln() + newFunc + fmt.Sprintln() + loadFunc + fmt.Sprintln() + saveFunc

	JSCode = fmt.Sprintf(`app.controller('%s', 
		['$scope', '$http', '$location', '$routeParams', 'apiPath', 'appVars',
			function($scope, $http, $location, $routeParams, apiPath, appVars) {
		%s
		//check if the user has access to this page	
		checkPageAccess($location, appVars.user);	
		$scope.%s = $routeParams.%s;
		if ($scope.%s()){//New 
			$scope.%s();
		}else{
			$scope.%s();
		}
	}]);`, a.controllerName, JSCode, a.idCol, a.idCol, a.isNewFunc, a.newFunc, a.loadFunc)

	return
}

func (m *Model) GetClientIndexController(a *ClientModelSettings) (fileName, JSCode string) {
	fileName = path.Join(m.appRef.GetClientSettings().directories["app"], m.Name, a.indexControllerFileName)

	//modelIndexFunc
	listFunc := fmt.Sprintf(
		`$scope.%s = function(){
			$http.get(apiPath + "%s")
				.then(function(response) {
					if (response.status == 200){
						$scope.%s = response.data;
						clearAPIError($scope);
					}
				},
				function(response) {
					handleAPIError($scope, response);
				}
			);
		};`, a.indexFunc, a.indexRoute, a.indexData)

	//modelDeleteFunc
	deleteFunc := fmt.Sprintf(`$scope.%s = function(%s){
			$http.put(apiPath + "%s", {id: %s})
					.then(function(response) {
						if (response.status == 200){
							$scope.%s.removeByKey("id", %s);
							clearAPIError($scope);
						}
					},
					function(response) {
						handleAPIError($scope, response);
					}
				);
		};`, a.deleteFunc, a.idCol, a.deleteRoute, a.idCol, a.indexData, a.idCol)

	JSCode = listFunc + fmt.Sprintln() + deleteFunc

	JSCode = fmt.Sprintf(`app.controller('%s', 
		['$scope', '$http', '$location', 'apiPath', 'appVars',
			function($scope, $http, $location, apiPath, appVars) {
		%s
		//check if the user has access to this page	
		checkPageAccess($location, appVars.user);	

		$scope.%s();
	}]);`, a.indexControllerName, JSCode, a.indexFunc)

	return
}

func (m *Model) GetRoutes(a *ClientModelSettings) (routes string) {
	indexRoute := fmt.Sprintf(
		`.when('%s', {
			templateUrl: 'app/%s/%s',
			controller: '%s',
			title : '%s List'
		})`, a.indexRoute, m.Name, a.indexViewFileName, a.indexControllerName, strings.Title(m.DisplayName))

	newRoute := fmt.Sprintf(
		`.when('%s', {
			templateUrl: 'app/%s/%s',
			controller: '%s',
			title : 'New %s'
		})`, a.newRoute, m.Name, a.editViewFileName, a.controllerName, strings.Title(m.DisplayName))

	editRoute := fmt.Sprintf(
		`.when('%s:%s', {
			templateUrl: 'app/%s/%s',
			controller: '%s',
			title : 'Edit %s'
		})`, a.editRoute, a.idCol, m.Name, a.editViewFileName, a.controllerName, strings.Title(m.DisplayName))

	showRoute := fmt.Sprintf(
		`.when('%s:%s', {
			templateUrl: 'app/%s/%s',
			controller: '%s',
			title : '%s Details'
		})`, a.showRoute, a.idCol, m.Name, a.showViewFileName, a.controllerName, strings.Title(m.DisplayName))

	routes = indexRoute + fmt.Sprintln() + newRoute + fmt.Sprintln() + editRoute + fmt.Sprintln() + showRoute
	return
}
