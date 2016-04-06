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
	loadFunc                string
	validateFunc            string
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
		loadFunc:                fmt.Sprintf("Load%s", strings.Title(m.Name)),
		validateFunc:            fmt.Sprintf("Validate%s", strings.Title(m.Name)),
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
		deleteRoute:             fmt.Sprintf("/%s/", m.Name),
	}
}

func (m *Model) GetClientEditView(a *ClientModelSettings) (fileName, htmlCode string) {
	fileName = path.Join(m.appRef.GetClientSettings().directories["app"], m.Name, a.editViewFileName)

	htmlCode = `<div class="row" ng-include="'/app/error/errorhandler.view.html'"></div><div class="row text-center">`
	htmlCode += fmt.Sprintf(`<h3 ng-bind="(%s()?'New':'Edit') + ' %s'"></h3>`, a.isNewFunc, m.DisplayName)
	htmlCode += `<hr/></div>`

	htmlCode += fmt.Sprintf(`<form class="form-horizontal" role="form" name="%s" ng-submit="%s()"><div class="row">`, a.formName, a.saveFunc)
	for _, fld := range m.Fields {
		if fld.HideInEdit {
			continue
		}
		inputHtml := ""
		if fld.Type == Boolean {
			inputHtml = fmt.Sprintf(`<input id="%s" ng-model="%s.%s" title="%s" type="checkbox" `,
				fld.Name, a.formData, fld.Name, fld.DisplayName)
		} else {
			inputHtml = fmt.Sprintf(`<input id="%s" ng-model="%s.%s" class="form-control" placeholder="Enter %s" title="%s" `,
				fld.Name, a.formData, fld.Name, fld.DisplayName, fld.DisplayName)
			if fld.Type == Date {
				inputHtml += `type="date" `
			} else if fld.Type == Integer {
				inputHtml += `type="number" step="1" `
			} else if fld.Type == Float {
				inputHtml += `type="number" step=".01" `
			} else if fld.Type == String {
				if fld.Validator != nil && fld.Validator.Email {
					inputHtml += `type="email" `
				} else if fld.Validator != nil && fld.Validator.Url {
					inputHtml += `type="url" `
				} else {
					inputHtml += `type="text" `
				}
			} else {
				inputHtml += `type="text" `
			}
		}

		if fld.Validator != nil {
			if fld.Validator.MinLen > 0 {
				inputHtml += fmt.Sprintf(` minlength="%d"`, fld.Validator.MinLen)
			}
			if fld.Validator.MaxLen > 0 {
				inputHtml += fmt.Sprintf(` maxlength="%d"`, fld.Validator.MaxLen)
			}
			if fld.Validator.MinValue > 0 {
				inputHtml += fmt.Sprintf(` min="%d"`, fld.Validator.MinValue)
			}
			if fld.Validator.MaxValue > 0 {
				inputHtml += fmt.Sprintf(` max="%d"`, fld.Validator.MaxValue)
			}
			if fld.Validator.Email {
				inputHtml += ` pattern="[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,3}$"`
			}
			if fld.Validator.Url {
				inputHtml += ` pattern="https?://.+"`
			}
			if fld.Validator.IsAlpha {
				inputHtml += ` pattern="^[A-Za-z ]+$"`
			}
			if fld.Validator.IsAlphaNumeric {
				inputHtml += ` pattern="^[A-Za-z0-9 ]+$"`
			}
			if fld.Validator.Required {
				inputHtml += ` required`
			}
		}
		inputHtml += `/>`

		if fld.Type == Boolean {
			htmlCode += `<div class="form-group"><div class="col-sm-offset-4 col-sm-8">`
			inputHtml = fmt.Sprintf(`<label class="checkbox-inline">%s%s</label>`, inputHtml, fld.DisplayName)
		} else {
			htmlCode += fmt.Sprintf(`<div class="form-group">
			<label class="control-label col-sm-4" for="%s">%s:</label>
			<div class="col-sm-8">`, fld.Name, fld.DisplayName)
		}
		htmlCode += inputHtml + `</div></div>`
	}
	htmlCode += fmt.Sprintf(`<div class="form-group">
		<div class="col-sm-offset-2 col-sm-4">
			<a ng-href="{{'#%s'}}" alt="click to go back to %s List"><span class="glyphicon glyphicon-arrow-left"></span> Back</a>
		</div>
		<div class="col-sm-offset-2 col-sm-4">
			<button type="submit" class="btn btn-link"><span class="glyphicon glyphicon-save"></span> Save</button>
		</div>
	</div>
	</div></form>
	<div class="row"><hr/></div>`, a.indexRoute, m.DisplayName)

	htmlCode = gohtml.Format(htmlCode)
	return
}

func (m *Model) GetClientShowView(a *ClientModelSettings) (fileName, htmlCode string) {
	fileName = path.Join(m.appRef.GetClientSettings().directories["app"], m.Name, a.showViewFileName)

	htmlCode = fmt.Sprintf(`<div class="row" ng-include="'/app/error/errorhandler.view.html'"></div>`+
		`<div class="row text-center"><h3>%s Details</h3><hr/></div>`, m.DisplayName)

	htmlCode += `<div class="row"><div class="col-sm-12">`
	for _, fld := range m.Fields {
		htmlCode += fmt.Sprintf(`<div class="row"><div class="col-sm-12"><h4>%s</h4><p ng-bind="%s.%s"></p></div></div>`,
			fld.DisplayName, a.formData, fld.Name)
	}
	htmlCode += fmt.Sprintf(`<div class="row">
			<div class="col-sm-12">
				<a ng-href="{{'#%s'}}" alt="click to go back to %s List"><span class="glyphicon glyphicon-arrow-left"></span> Back</a>
			</div>
		</div>`, a.indexRoute, m.DisplayName)
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
		a.indexFunc, a.newRoute, m.DisplayName)

	htmlCode += fmt.Sprintf(`<div ng-if="%s.length==0" class="row"><div class="col-sm-12 text-center"><h3>0 Records Found.</h3></div></div>`,
		a.indexData)

	htmlCode += fmt.Sprintf(`<div ng-if="%s.length>0" class="row"><div class="col-sm-12">`, a.indexData)
	if m.ViewType == List {
		htmlCode += fmt.Sprintf(`<div class="row" ng-repeat="x in %s | orderBy:modifiedAt:reverse">`, a.indexData)
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
				htmlCode += fmt.Sprintf(`<div class="col-sm-1"><span ng-bind="x.%s"></span></div>`,
					fld.Name)
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

	//modelLoadFunc
	loadFunc := fmt.Sprintf(
		`//function to load model entity
		$scope.%s =function(){
		$http.get(apiPath + "/%s/" + $scope.%s)
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
	};`, a.loadFunc, m.Name, a.idCol, a.formData)

	validateFunc := ""
	for _, fld := range m.Fields {
		if fld.HideInEdit {
			continue
		}
		validateFunc += fld.GetClientValidation(a)
	}
	validateFunc = fmt.Sprintf(
		`$scope.%s=function(){
			$scope.errors = []
			%s
			return ($scope.errors.length==0);
		};`, a.validateFunc, validateFunc)
	//modelSaveFunc
	saveFunc := fmt.Sprintf(
		`//function to save model entity
		$scope.%s =function(){
			if (!$scope.%s()){
				handleAPIError($scope, {status:404,data:{errors:$scope.errors}});
				return
			}
			$http({
					method: $scope.%s()?'POST':'PUT',
					url: apiPath + "%s",
					data: $scope.%s
				}).then(
				function(response) {
					if (response.status == 200){
						clearAPIError($scope);
						$location.path("%s");
					} 
					else {
					  $scope.message = data.message;
					}
				},
				function(response){
					handleAPIError($scope, response);
			  });
		};`, a.saveFunc, a.validateFunc, a.isNewFunc, a.saveRoute, a.formData, a.indexRoute)

	JSCode = isNewFunc + fmt.Sprintln() + loadFunc +
		fmt.Sprintln() + validateFunc + fmt.Sprintln() + saveFunc

	JSCode = fmt.Sprintf(`app.controller('%s', 
		['$scope', '$http', '$location', '$routeParams', 'apiPath', 'appVars',
			function($scope, $http, $location, $routeParams, apiPath, appVars) {
		%s
		//check if the user has access to this page	
		checkPageAccess($location, appVars.user);	
		$scope.%s = $routeParams.%s;
		if (!$scope.%s()){//New 
			$scope.%s();
		}
	}]);`, a.controllerName, JSCode, a.idCol, a.idCol, a.isNewFunc, a.loadFunc)

	return
}

func (m *Model) GetClientIndexController(a *ClientModelSettings) (fileName, JSCode string) {
	fileName = path.Join(m.appRef.GetClientSettings().directories["app"], m.Name, a.indexControllerFileName)

	//modelIndexFunc
	listFunc := fmt.Sprintf(
		`$scope.%s = function(){
			$http.get(apiPath + "/%s")
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
		};`, a.indexFunc, m.Name, a.indexData)

	//modelDeleteFunc
	deleteFunc := fmt.Sprintf(`$scope.%s = function(%s){
			$http.delete(apiPath + "%s" + %s)
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
