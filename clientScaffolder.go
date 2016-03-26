package main

import (
	"fmt"
	"strings"
)

type ModelTokens struct {
	idCol                   string
	modelForm               string
	modelFormData           string
	modelIndexData          string
	modelIndexFunc          string
	modelIsNewFunc          string
	modelLoadFunc           string
	modelSaveFunc           string
	modelDeleteFunc         string
	indexViewFileName       string
	editViewFileName        string
	showViewFileName        string
	controllerFileName      string
	indexControllerFileName string
	indexRoute              string
	newRoute                string
	editRoute               string
	showRoute               string
	deleteRoute             string
}

func (m *Model) GetTokens() *ModelTokens {
	return &ModelTokens{
		idCol:                   fmt.Sprintf("%sId", m.Name),
		modelForm:               fmt.Sprintf("%sForm", m.Name),
		modelFormData:           fmt.Sprintf("%sData", m.Name),
		modelIndexData:          fmt.Sprintf("%sList", m.Name),
		modelIndexFunc:          fmt.Sprintf("Get%sList", strings.Title(m.Name)),
		modelIsNewFunc:          fmt.Sprintf("IsNew%s", strings.Title(m.Name)),
		modelLoadFunc:           fmt.Sprintf("Load%s", strings.Title(m.Name)),
		modelSaveFunc:           fmt.Sprintf("Save%s", strings.Title(m.Name)),
		modelDeleteFunc:         fmt.Sprintf("Delete%s", strings.Title(m.Name)),
		indexViewFileName:       fmt.Sprintf("%sIndex.view.htm", m.Name),
		editViewFileName:        fmt.Sprintf("%sEdit.view.htm", m.Name),
		showViewFileName:        fmt.Sprintf("%sShow.view.htm", m.Name),
		controllerFileName:      fmt.Sprintf("%s.controller.js", m.Name),
		indexControllerFileName: fmt.Sprintf("%s.index.controller.js", m.Name),
		indexRoute:              fmt.Sprintf("/%s/list", m.Name),
		newRoute:                fmt.Sprintf("/%s/new", m.Name),
		editRoute:               fmt.Sprintf("/%s/edit/", m.Name),
		showRoute:               fmt.Sprintf("/%s/view/", m.Name),
		deleteRoute:             fmt.Sprintf("/%s/delete", m.Name),
	}
}

func (m *Model) GetClientEditView(a *ModelTokens) (fileName, htmlCode string) {
	fileName = a.editViewFileName

	htmlCode = `<div class="row" ng-include="'/app/error/errorhandler.view.html'"></div><div class="row text-center">`
	htmlCode += fmt.Sprintf(`<h3 ng-bind="(%s()?'New':'Edit') + ' %s'"></h3>`, a.modelIsNewFunc, m.DisplayName)
	htmlCode += `<hr/></div>`

	htmlCode += fmt.Sprintf(`<form class="form-horizontal" role="form" name="%s" ng-submit="%s()"><div class="row">`, a.modelForm, a.modelSaveFunc)
	for _, fld := range m.Fields {
		htmlCode += `<div class="form-group">`
		htmlCode += fmt.Sprintf(`<label class="control-label col-sm-4" for="%s">%s:</label>`, fld.Name, fld.DisplayName)
		htmlCode += `<div class="col-sm-8">`
		htmlCode += fmt.Sprintf(`<input type="text" class="form-control" id="%s" placeholder="Enter %s" ng-model="%s.%s" title="%s" `,
			fld.Name, fld.DisplayName, a.modelFormData, fld.Name, fld.DisplayName)

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
	htmlCode += `<div class="row"><hr></div>`
	return
}

func (m *Model) GetClientShowView(a *ModelTokens) (fileName, htmlCode string) {
	fileName = a.showViewFileName

	htmlCode = fmt.Sprintf(`<div class="row" ng-include="'/app/error/errorhandler.view.html'"></div>`+
		`<div class="row text-center"><h3>%s Details</h3><hr></div>`, m.DisplayName)

	htmlCode += `<div class="row"><div class="col-sm-12">`
	for _, fld := range m.Fields {
		htmlCode += fmt.Sprintf(`<div class="row"><div class="col-sm-12"><h3 ng-bind="%s"></h3><p ng-bind="%s.%s"></p></div></div>`,
			fld.DisplayName, a.modelFormData, fld.Name)
	}
	htmlCode += `</div></div>`
	htmlCode += `<div class="row"><hr></div>`
	return
}

func (m *Model) GetClientIndexView(a *ModelTokens) (fileName, htmlCode string) {
	fileName = a.indexViewFileName

	htmlCode = fmt.Sprintf(`<div class="row" ng-include="'/app/error/errorhandler.view.html'"></div>`+
		`<div class="row text-center"><h3>%s List</h3><hr></div>`, m.DisplayName)

	htmlCode += fmt.Sprintf(`<div class="row text-center"><a href="" ng-click="%s()"><span class="glyphicon glyphicon-refresh"/> Refresh</a>`+
		`<a href="%s" class="col-sm-offset-1"><span class="glyphicon glyphicon-plus"/> New %s</a></div><br/>`,
		a.modelLoadFunc, a.newRoute, m.DisplayName)

	htmlCode += fmt.Sprintf(`<div ng-if="%s.length==0" class="row"><div class="col-sm-12 text-center"><h3>0 Records Found.</h3></div></div>`,
		a.modelFormData)

	htmlCode += fmt.Sprintf(`<div ng-if="%s.length>0" class="row"><div class="col-sm-12">`, a.modelFormData)
	if m.ViewType == List {
		htmlCode += fmt.Sprintf(`<div class="row" ng-repeat="x in %s | orderBy:createdOn:reverse">`, a.modelFormData)
		htmlCode += fmt.Sprintf(`<div class="col-sm-1"><a href="%s{{x.id}}" alt="View %s" title="View %s">`+
			`<span class="glyphicon glyphicon-folder-open"></span></a></div>`,
			a.showRoute, m.DisplayName, m.DisplayName)
		htmlCode += fmt.Sprintf(`<div class="col-sm-1"><a href="%s{{x.id}}" alt="Edit %s" title="Edit %s">`+
			`<span class="glyphicon glyphicon-edit"></span></a></div>`,
			a.editRoute, m.DisplayName, m.DisplayName)
		htmlCode += fmt.Sprintf(`<div class="col-sm-1"><a href="" alt="Delete %s" title="Delete %s"><span class="glyphicon glyphicon-remove" `+
			`ng-click="%s(x.id);"></span></a></div>`,
			m.DisplayName, m.DisplayName, a.modelDeleteFunc)
		for _, fld := range m.Fields {
			if fld.ShowInIndex {
				htmlCode += fmt.Sprintf(`<div class="col-sm-1"><span ng-bind="%s.%s"></span></div>`,
					a.modelFormData, fld.Name)
			}
		}
		htmlCode += `</div>`
	} else if m.ViewType == Table {

	}
	htmlCode += `</div></div>`
	htmlCode += `<div class="row"><hr></div>`
	return
}

func (m *Model) GetClientController(a *ModelTokens) (fileName, JSCode string) {
	fileName = a.controllerFileName

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
		};`, a.modelIndexFunc, a.indexRoute, a.modelIndexData)

	//modelIsNewFunc
	isNewFunc := fmt.Sprintf(
		`$scope.%s = function(){
		return (!$scope.%s || $scope.%s == "" || $scope.%s == null);
	}`, a.modelIsNewFunc, a.idCol)

	//modelLoadFunc
	loadFunc := fmt.Sprintf(`$scope.%s =function(){
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
	}`, a.modelLoadFunc, m.Name, a.idCol, a.modelFormData, m.Name)

	//modelSaveFunc
	saveFunc := fmt.Sprintf(
		`$scope.%s =function(){
			$http({
					method: %s()?'POST:'PUT'',
					url: apiPath + "/%s",
					data: $scope.%s
				}).then(
				function(response) {
					if (response.status == 200){
						clearAPIError($scope);
						$location.path("/%s");
					} 
					else {
					  $scope.message = data.message;
					}
				},
				function(response){
					handleAPIError($scope, response);
			  });
		};`, a.modelSaveFunc, a.modelIsNewFunc, m.Name, a.modelFormData, a.indexRoute)

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
		};`, a.modelDeleteFunc, a.idCol, a.deleteRoute, a.idCol, a.modelIndexData, a.idCol)

	JSCode = listFunc + isNewFunc + loadFunc + saveFunc + deleteFunc
	return
}

func (m *Model) GetClientIndexController(a *ModelTokens) (fileName, JSCode string) {
	fileName = a.indexControllerFileName
	return
}
