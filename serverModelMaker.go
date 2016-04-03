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
	getRoute           string
	newRoute           string
	saveRoute          string
	deleteRoute        string
}

func (m *Model) GetServerSettings() *ServerModelSettings {
	return &ServerModelSettings{
		idCol:              fmt.Sprintf("%sId", m.Name),
		indexFunc:          "GetAll",
		newFunc:            "Create",
		loadFunc:           "GetById",
		saveFunc:           "Save",
		deleteFunc:         "DeleteById",
		controllerVar:      fmt.Sprintf("%sController", m.Name),
		controllerName:     fmt.Sprintf("%sController", strings.Title(m.Name)),
		controllerFileName: fmt.Sprintf("%sController.go", m.Name),
		modelName:          fmt.Sprintf("%sModel", strings.Title(m.Name)),
		modelFileName:      fmt.Sprintf("%sModel.go", m.Name),
		indexRoute:         fmt.Sprintf("/%s", m.Name),
		getRoute:           fmt.Sprintf("/%s/:%s", m.Name, fmt.Sprintf("%sId", m.Name)),
		newRoute:           fmt.Sprintf("/%s", m.Name),
		saveRoute:          fmt.Sprintf("/%s", m.Name),
		deleteRoute:        fmt.Sprintf("/%s/:%s", m.Name, fmt.Sprintf("%sId", m.Name)),
	}
}

func (m *Model) GetServerModel(a *ServerModelSettings) (fileName, goCode string) {
	fileName = path.Join(m.appRef.GetServerSettings().directories["server"], a.modelFileName)

	fieldDefs := ""
	fieldChecks := ""
	imports := ""
	for _, fld := range m.Fields {
		field := strings.Title(fld.Name)
		if fld.Type == Boolean {
			fieldDefs += fmt.Sprintf(`%s bool`, field)
		} else if fld.Type == Date {
			fieldDefs += fmt.Sprintf(`%s time.Time`, field)
			imports = `import (
					"time"
					)`
		} else if fld.Type == Integer {
			fieldDefs += fmt.Sprintf(`%s uint64`, field)
		} else if fld.Type == Float {
			fieldDefs += fmt.Sprintf(`%s float64`, field)
		} else if fld.Type == String {
			fieldDefs += fmt.Sprintf(`%s string`, field)
		} else {
			fieldDefs += fmt.Sprintf(`%s string`, field)
		}
		fieldDefs += fmt.Sprintf(" `json:\"%s\"` ", fld.Name)

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

			%s

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
			}`, imports, a.modelName, fieldDefs, a.modelName, fieldChecks)

	return
}

func (m *Model) GetServerController(a *ServerModelSettings) (fileName, goCode string) {
	fileName = path.Join(m.appRef.GetServerSettings().directories["server"], a.controllerFileName)

	//modelListFunc
	indexFunc := fmt.Sprintf(
		`//function to Get List of model
		func (c *%s) %s() (%sList []*%s){
			fmt.Println("%s.%s executed")
			return
		}`, a.controllerName, a.indexFunc, m.Name, a.modelName, a.controllerName, a.indexFunc)

	//modelNewFunc
	newFunc := fmt.Sprintf(
		`//function to Create New model entity
		func (c *%s) %s(%s *%s) (ok bool, modelErrors []string){
			fmt.Println("%s.%s executed")
			ok, modelErrors = %s.Validate()
			if !ok{
				return ok, modelErrors
			}			
			return
		}`, a.controllerName, a.newFunc, m.Name, a.modelName, a.controllerName, a.newFunc, m.Name)

	//modelLoadFunc
	loadFunc := fmt.Sprintf(
		`//function to Get model entity by id
		func (c *%s) %s(%s uint64) (%s *%s){
			fmt.Println("%s.%s executed")
			return
		}`, a.controllerName, a.loadFunc, a.idCol, m.Name, a.modelName, a.controllerName, a.loadFunc)

	//modelSaveFunc
	saveFunc := fmt.Sprintf(
		`//function to save model entity
		func (c *%s) %s(%s *%s) (ok bool, modelErrors []string){
			fmt.Println("%s.%s executed")
			ok, modelErrors = %s.Validate()
			if !ok{
				return ok, modelErrors
			}			
			return
		}`, a.controllerName, a.saveFunc, m.Name, a.modelName, a.controllerName, a.saveFunc, m.Name)

	//modelSaveFunc
	deleteFunc := fmt.Sprintf(
		`//function to delete model entity by id
		func (c *%s) %s(%s uint64) (ok bool, err error){
			fmt.Println("%s.%s executed")						
			return true, nil
		}`, a.controllerName, a.deleteFunc, a.idCol, a.controllerName, a.deleteFunc)

	goCode = indexFunc + fmt.Sprintln() + newFunc + fmt.Sprintln() + loadFunc +
		fmt.Sprintln() + saveFunc + fmt.Sprintln() + deleteFunc

	formFld := ""
	timePack := ""
	for _, fld := range m.Fields {
		parse := ""
		switch fld.Type {
		case String:
			formFld += fmt.Sprintf(`model.%s = req.PostFormValue("%s")`, strings.Title(fld.Name), fld.Name) + fmt.Sprintln()
		case Date, Boolean, Float, Integer:
			switch fld.Type {
			case Date:
				parse = fmt.Sprintf(`%s, _ := time.Parse(longTimeForm, req.PostFormValue("%s"))`, fld.Name, fld.Name)
				timePack = `"time"`
			case Boolean:
				parse = fmt.Sprintf(`%s, _ := strconv.ParseBool(req.PostFormValue("%s"))`, fld.Name, fld.Name)
			case Float:
				parse = fmt.Sprintf(`%s, _ := strconv.ParseFloat(req.PostFormValue("%s"), 64)`, fld.Name, fld.Name)
			case Integer:
				parse = fmt.Sprintf(`%s, _ := strconv.ParseInt(req.PostFormValue("%s"), 10, 64)`, fld.Name, fld.Name)
			}
			formFld += fmt.Sprintf(
				`if req.PostFormValue("%s") != "" {
							%s
							model.%s = %s
				}`, fld.Name, parse, strings.Title(fld.Name), fld.Name) + fmt.Sprintln()
		}
	}

	goCode = fmt.Sprintf(`package main
		
			import(
				"encoding/json"
				"strings"
				"net/http"
				"fmt"
				"strconv"
				%s
			)

			type %s struct{
				Name string
			}

			var (
				%s *%s = &%s{
					Name:"%s",
				}
			)

			func (c *%s) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
				var model %s
				switch req.Method {
				case "POST", "PUT" :
					contentType := strings.ToLower(req.Header["Content-Type"][0])
					switch {
					case strings.Contains(contentType, "application/json"),
						strings.Contains(contentType, "text/plain"):
						json.NewDecoder(req.Body).Decode(&model)
					case strings.Contains(contentType, "application/x-www-form-urlencoded"):
						%s
					case strings.Contains(contentType, "multipart-form-data"):
					}
				case "GET":
					if req.URL.String() == "" {
						c.GetAll()
					} else{
						id, _ := strconv.ParseUint(req.PostFormValue("%s"), 10, 64)
						c.GetById(id)
					}
				case "DELETE":
				case "PATCH":
				}
			}
			%s`, timePack, a.controllerName, a.controllerVar, a.controllerName, a.controllerName,
		a.controllerName, a.controllerName, a.modelName, formFld, a.idCol, goCode)

	return
}

func (m *Model) GetServerRoutes(s *ServerAppSettings) (routes string) {
	a := m.GetServerSettings()
	handlerPath := path.Join(s.apiPath, m.Name)
	indexRoute := fmt.Sprintf(
		`http.Handle("%s/", http.StripPrefix("%s/", %s))
		http.Handle("%s", http.StripPrefix("%s", %s))`,
		handlerPath, handlerPath, a.controllerVar,
		handlerPath, handlerPath, a.controllerVar)

	routes = indexRoute + fmt.Sprintln()
	return
}
