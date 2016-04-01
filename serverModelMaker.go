package main

import (
	"fmt"
	"path"
	"strings"
)

type ServerModelSettings struct {
	idCol              string
	indexFunc          string
	newFunc            string
	loadFunc           string
	saveFunc           string
	deleteFunc         string
	controllerVar      string
	controllerName     string
	controllerFileName string
	modelName          string
	modelFileName      string
	indexRoute         string
	newRoute           string
	getRoute           string
	saveRoute          string
	deleteRoute        string
}

func (m *Model) GetServerSettings() *ServerModelSettings {
	return &ServerModelSettings{
		idCol:              fmt.Sprintf("%sId", m.Name),
		indexFunc:          fmt.Sprintf("Get%sList", strings.Title(m.Name)),
		newFunc:            fmt.Sprintf("New%s", strings.Title(m.Name)),
		loadFunc:           fmt.Sprintf("Load%s", strings.Title(m.Name)),
		saveFunc:           fmt.Sprintf("Save%s", strings.Title(m.Name)),
		deleteFunc:         fmt.Sprintf("Delete%s", strings.Title(m.Name)),
		controllerVar:      fmt.Sprintf("%sController", m.Name),
		controllerName:     fmt.Sprintf("%sController", strings.Title(m.Name)),
		controllerFileName: fmt.Sprintf("%sController.go", m.Name),
		modelName:          fmt.Sprintf("%sModel", strings.Title(m.Name)),
		modelFileName:      fmt.Sprintf("%sModel.go", m.Name),
		indexRoute:         fmt.Sprintf("/%s/list", m.Name),
		newRoute:           fmt.Sprintf("/%s/new", m.Name),
		getRoute:           fmt.Sprintf("/%s/:%s", m.Name, fmt.Sprintf("%sId", m.Name)),
		saveRoute:          fmt.Sprintf("/%s", m.Name),
		deleteRoute:        fmt.Sprintf("/%s/delete/:%s", m.Name, fmt.Sprintf("%sId", m.Name)),
	}
}

func (m *Model) GetServerModel(a *ServerModelSettings) (fileName, goCode string) {
	fileName = path.Join(m.appRef.GetServerSettings().directories["server"], a.modelFileName)

	fieldDefs := ""
	fieldChecks := ""
	for _, fld := range m.Fields {
		if fld.Type == Boolean {
			fieldDefs += fmt.Sprintf(`%s bool`, fld.Name)
		} else if fld.Type == Date {
			fieldDefs += fmt.Sprintf(`%s Time`, fld.Name)
		} else if fld.Type == Integer {
			fieldDefs += fmt.Sprintf(`%s uint64`, fld.Name)
		} else if fld.Type == Float {
			fieldDefs += fmt.Sprintf(`%s float64`, fld.Name)
		} else if fld.Type == String {
			fieldDefs += fmt.Sprintf(`%s string`, fld.Name)
		} else {
			fieldDefs += fmt.Sprintf(`%s string`, fld.Name)
		}

		fieldDefs += ` //`
		if fld.Validator.MinLen > 0 {
			fieldDefs += fmt.Sprintf(` minlength="%d",`, fld.Validator.MinLen)
		}
		if fld.Validator.MaxLen > 0 {
			fieldDefs += fmt.Sprintf(` maxlength="%d",`, fld.Validator.MaxLen)
		}
		if fld.Validator.MinValue > 0 {
			fieldDefs += fmt.Sprintf(` min="%d",`, fld.Validator.MinValue)
		}
		if fld.Validator.MaxValue > 0 {
			fieldDefs += fmt.Sprintf(` max="%d",`, fld.Validator.MaxValue)
		}
		if fld.Validator.Email {
			fieldDefs += ` valid email,`
		}
		if fld.Validator.Url {
			fieldDefs += ` valid url,`
		}
		if fld.Validator.IsAlpha {
			fieldDefs += ` isAlpha,`
		}
		if fld.Validator.IsAlphaNumeric {
			fieldDefs += ` isAlphaNumeric,`
		}
		if fld.Validator.Required {
			fieldDefs += ` required,`
		}
		fieldDefs += fmt.Sprintln()

		fieldChecks += fld.GetServerValidation(a) + fmt.Sprintln()
	}

	goCode = fmt.Sprintf(`package main

			type %s struct{
				%s
			}

			func (m *%s) Validate() (ok bool, modelErrors []string){
				%s
				ok = (len(modelErrors) >0)
				if ok {
					modelErrors=nil
				}
				return
			}`, a.modelName, fieldDefs, a.modelName, fieldChecks)

	return
}

func (m *Model) GetServerController(a *ServerModelSettings) (fileName, goCode string) {
	fileName = path.Join(m.appRef.GetServerSettings().directories["server"], a.controllerFileName)

	//modelListFunc
	indexFunc := fmt.Sprintf(
		`//function to Get List of model
		func (c *%s) %s() (%sList []*%s){
			return
		}`, a.controllerName, a.indexFunc, m.Name, a.modelName)

	//modelNewFunc
	newFunc := fmt.Sprintf(
		`//function to Get New model entity
		func (c *%s) %s() (%s *%s){
			return
		}`, a.controllerName, a.newFunc, m.Name, a.modelName)

	//modelLoadFunc
	loadFunc := fmt.Sprintf(
		`//function to Get model entity by id
		func (c *%s) %s(%s uint64) (%s *%s){
			return
		}`, a.controllerName, a.loadFunc, a.idCol, m.Name, a.modelName)

	//modelSaveFunc
	saveFunc := fmt.Sprintf(
		`//function to Get model entity by id
		func (c *%s) %s(%s *%s) (ok bool, modelErrors []string){
			ok, modelErrors = %s.Validate()
			if !ok{
				return ok, modelErrors
			}			
			return
		}`, a.controllerName, a.saveFunc, m.Name, a.modelName, m.Name)

	//modelSaveFunc
	deleteFunc := fmt.Sprintf(
		`//function to Get model entity by id
		func (c *%s) %s(%s uint64) (ok bool, err error){
						
			return true, nil
		}`, a.controllerName, a.deleteFunc, a.idCol)

	goCode = indexFunc + fmt.Sprintln() + newFunc + fmt.Sprintln() + loadFunc +
		fmt.Sprintln() + saveFunc + fmt.Sprintln() + deleteFunc

	goCode = fmt.Sprintf(`package main
		
			import(
				"net/http"
			)

			type %s struct{
				Name string
			}

			var %s *%s = &%s{
				Name:"%s",
			}

			func (c *%s) HandleAction(rw http.ResponseWriter, req *http.Request) {
			}

			%s`, a.controllerName, a.controllerVar, a.controllerName, a.controllerName, a.controllerName, a.controllerName, goCode)

	return
}

func (m *Model) GetServerRoutes(s *ServerAppSettings) (routes string) {
	a := m.GetServerSettings()

	indexRoute := fmt.Sprintf(
		`http.HandleFunc("%s", %s.HandleAction)`,
		path.Join(s.apiPath, a.indexRoute), a.controllerVar)

	newRoute := fmt.Sprintf(
		`http.HandleFunc("%s", %s.HandleAction)`,
		path.Join(s.apiPath, a.newRoute), a.controllerVar)

	getRoute := fmt.Sprintf(
		`http.HandleFunc("%s", %s.HandleAction)`,
		path.Join(s.apiPath, a.getRoute), a.controllerVar)

	saveRoute := fmt.Sprintf(
		`http.HandleFunc("%s", %s.HandleAction)`,
		path.Join(s.apiPath, a.saveRoute), a.controllerVar)

	deleteRoute := fmt.Sprintf(
		`http.HandleFunc("%s", %s.HandleAction)`,
		path.Join(s.apiPath, a.deleteRoute), a.controllerVar)

	routes = indexRoute + fmt.Sprintln() + newRoute + fmt.Sprintln() + getRoute +
		fmt.Sprintln() + saveRoute + fmt.Sprintln() + deleteRoute
	return
}
