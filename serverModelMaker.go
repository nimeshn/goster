package main

import (
	"fmt"
	"path"
	"strings"
)

type ServerModelSettings struct {
	idCol              string
	indexFunc          string
	getFunc            string
	createFunc         string
	updateFunc         string
	deleteFunc         string
	getFuncFormat      string
	deleteFuncFormat   string
	indexRoute         string
	getRoute           string
	createRoute        string
	updateRoute        string
	deleteRoute        string
	controllerVar      string
	controllerName     string
	controllerFileName string
	modelName          string
	modelFileName      string
}

func (m *Model) GetServerSettings() *ServerModelSettings {
	return &ServerModelSettings{
		idCol:              fmt.Sprintf("%sId", m.Name),
		indexFunc:          "GetAll",
		getFunc:            "GetById",
		createFunc:         "Create",
		updateFunc:         "Update",
		deleteFunc:         "DeleteById",
		getFuncFormat:      "GetBy%s",
		deleteFuncFormat:   "DeleteBy%s",
		indexRoute:         fmt.Sprintf("/%s", m.Name),
		getRoute:           fmt.Sprintf("/%s/:%s", m.Name, fmt.Sprintf("%sId", m.Name)),
		createRoute:        fmt.Sprintf("/%s", m.Name),
		updateRoute:        fmt.Sprintf("/%s", m.Name),
		deleteRoute:        fmt.Sprintf("/%s/:%s", m.Name, fmt.Sprintf("%sId", m.Name)),
		controllerVar:      fmt.Sprintf("%sController", m.Name),
		controllerName:     fmt.Sprintf("%sController", strings.Title(m.Name)),
		controllerFileName: fmt.Sprintf("%sController.go", m.Name),
		modelName:          fmt.Sprintf("%sModel", strings.Title(m.Name)),
		modelFileName:      fmt.Sprintf("%sModel.go", m.Name),
	}
}

func (m *Model) GetServerModel(sm *ServerModelSettings) (fileName, goCode string) {
	fileName = path.Join(m.appRef.GetServerSettings().directories["server"], sm.modelFileName)

	fieldDefs := ""
	fieldChecks := ""
	imports := ""
	for _, fld := range m.Fields {
		field := strings.Title(fld.Name)

		fieldDefs += fmt.Sprintln()
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

		if fld.Validator != nil {
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
		}

		fieldChecks += fld.GetServerValidation(sm) + fmt.Sprintln()
	}
	fieldDefs += fmt.Sprintln()

	goCode = fmt.Sprintf(`package main

			%s

			type %s struct{%s}

			//Validate data for this model
			func (m *%s) Validate() (ok bool, modelErrors []string){
				%s
				ok = (len(modelErrors) >0)
				if ok {
					modelErrors=nil
				}
				return
			}`, imports, sm.modelName, fieldDefs, sm.modelName, fieldChecks)

	return
}

func (m *Model) GetServerController(sm *ServerModelSettings) (fileName, goCode string) {
	fileName = path.Join(m.appRef.GetServerSettings().directories["server"], sm.controllerFileName)

	//modelListFunc
	indexFunc := fmt.Sprintf(
		`//function to Get List of model
		func (c *%s) %s() (%sList []*%s){
			fmt.Println("%s.%s executed")
			return
		}`, sm.controllerName, sm.indexFunc, m.Name, sm.modelName, sm.controllerName, sm.indexFunc)

	//modelLoadFunc
	getFunc := fmt.Sprintf(
		`//function to Get model entity by id
		func (c *%s) %s(%s uint64) (%s *%s){
			fmt.Println("%s.%s executed")
			return
		}`, sm.controllerName, sm.getFunc, sm.idCol, m.Name, sm.modelName, sm.controllerName, sm.getFunc)

	//modelNewFunc
	createFunc := fmt.Sprintf(
		`//function to Create New model entity
		func (c *%s) %s(%s *%s) (ok bool, modelErrors []string){
			fmt.Println("%s.%s executed")
			ok, modelErrors = %s.Validate()
			if !ok{
				return ok, modelErrors
			}			
			return
		}`, sm.controllerName, sm.createFunc, m.Name, sm.modelName, sm.controllerName, sm.createFunc, m.Name)

	//modelSaveFunc
	updateFunc := fmt.Sprintf(
		`//function to save model entity
		func (c *%s) %s(%s *%s) (ok bool, modelErrors []string){
			fmt.Println("%s.%s executed")
			ok, modelErrors = %s.Validate()
			if !ok{
				return ok, modelErrors
			}			
			return
		}`, sm.controllerName, sm.updateFunc, m.Name, sm.modelName, sm.controllerName, sm.updateFunc, m.Name)

	//modelDeleteFunc
	deleteFunc := fmt.Sprintf(
		`//function to delete model entity by id
		func (c *%s) %s(%s uint64) (ok bool, err error){
			fmt.Println("%s.%s executed")						
			return true, nil
		}`, sm.controllerName, sm.deleteFunc, sm.idCol, sm.controllerName, sm.deleteFunc)

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
				if fld.AutoGenerated {
					parse = fmt.Sprintf(`%s, _ := strconv.ParseUint(req.PostFormValue("%s"), 10, 64)`, fld.Name, fld.Name)
				} else {
					parse = fmt.Sprintf(`%s, _ := strconv.ParseInt(req.PostFormValue("%s"), 10, 64)`, fld.Name, fld.Name)
				}
			}
			formFld += fmt.Sprintf(
				`if req.PostFormValue("%s") != "" {
							%s
							model.%s = %s
				}`, fld.Name, parse, strings.Title(fld.Name), fld.Name) + fmt.Sprintln()
		}
	}

	uniqueGet, uniqueDel, uniqueFuncs := m.GetUniqueFieldAction(sm)

	serveHTTP := fmt.Sprintf(
		`//http Request handler for %s
		func (c *%s) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
			var model %s
			urlArr := strings.Split(req.URL.String(), "/")
			requestURL, requestField := "", ""
			if len(urlArr) > 1{
				requestField = urlArr[0]
				requestURL = urlArr[1]
			}else{
				requestURL = urlArr[0]
			}

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
				if req.Method == "POST" {
					c.%s(&model)
				} else {
					c.%s(&model)
				}
			case "GET":
				if requestField == "" {
					if requestURL == "" {
						c.%s()
					} else{
						id, _ := strconv.ParseUint(requestURL, 10, 64)
						c.%s(id)
					}
				}
				%s
			case "DELETE":
				if requestField == "" {
					if requestURL != "" {
						id, _ := strconv.ParseUint(requestURL, 10, 64)
						c.%s(id)
					}
				}
				%s
			case "PATCH":
			}
		}`, m.Name, sm.controllerName, sm.modelName, formFld, sm.createFunc, sm.updateFunc,
		sm.indexFunc, sm.getFunc, uniqueGet, sm.deleteFunc, uniqueDel)

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
			)`, timePack, sm.controllerName, sm.controllerVar, sm.controllerName, sm.controllerName,
		sm.controllerName)

	goCode += fmt.Sprintln() + serveHTTP + fmt.Sprintln() + indexFunc + fmt.Sprintln() + getFunc + fmt.Sprintln() +
		createFunc + fmt.Sprintln() + updateFunc + fmt.Sprintln() + deleteFunc + fmt.Sprintln() + uniqueFuncs

	return
}

func (m *Model) GetUniqueFieldAction(sm *ServerModelSettings) (get, del string, funcs string) {
	for _, fld := range m.Fields {
		if !fld.Unique {
			continue
		}
		parse := ""
		dataType := ""
		getFunc := fmt.Sprintf(sm.getFuncFormat, strings.Title(fld.Name))
		deleteFunc := fmt.Sprintf(sm.deleteFuncFormat, strings.Title(fld.Name))
		switch fld.Type {
		case String:
			parse = fmt.Sprintf(`%s := requestURL`, fld.Name)
			dataType = "string"
		case Date, Boolean, Float, Integer:
			switch fld.Type {
			case Date:
				parse = fmt.Sprintf(`%s, _ := time.Parse(longTimeForm, requestURL)`, fld.Name, fld.Name)
				dataType = "time.Time"
			case Boolean:
				parse = fmt.Sprintf(`%s, _ := strconv.ParseBool(requestURL)`, fld.Name, fld.Name)
				dataType = "bool"
			case Float:
				parse = fmt.Sprintf(`%s, _ := strconv.ParseFloat(requestURL, 64)`, fld.Name, fld.Name)
				dataType = "float64"
			case Integer:
				if fld.AutoGenerated {
					parse = fmt.Sprintf(`%s, _ := strconv.ParseUint(requestURL, 10, 64)`, fld.Name, fld.Name)
					dataType = "int64"
				} else {
					parse = fmt.Sprintf(`%s, _ := strconv.ParseInt(requestURL, 10, 64)`, fld.Name, fld.Name)
					dataType = "uint64"
				}
			}
		}
		get += fmt.Sprintf(
			`if requestField == "%s" {
				if requestURL != "" {
					%s
					c.%s(%s)
				}
			}`, fld.Name, parse, getFunc, fld.Name) + fmt.Sprintln()
		del += fmt.Sprintf(
			`if requestField == "%s" {
				if requestURL != "" {
					%s
					c.%s(%s)
				}
			}`, fld.Name, parse, deleteFunc, fld.Name) + fmt.Sprintln()

		funcs += fmt.Sprintf(
			`//function to Get model entity
			func (c *%s) %s(%s %s) (%s *%s){
				fmt.Println("%s.%s executed")
				return
			}`, sm.controllerName, getFunc, fld.Name, dataType, m.Name, sm.modelName,
			sm.controllerName, getFunc) + fmt.Sprintln() +

			fmt.Sprintf(
				`//function to delete model entity
			func (c *%s) %s(%s %s) (ok bool, err error){
				fmt.Println("%s.%s executed")						
				return true, nil
			}`, sm.controllerName, deleteFunc, fld.Name, dataType, sm.controllerName, deleteFunc)
	}
	return
}

func (m *Model) GetServerRoutes(sa *ServerAppSettings) (routes string) {
	ss := m.GetServerSettings()
	handlerPath := path.Join(sa.apiPath, m.Name)
	indexRoute := fmt.Sprintf(
		`//routes handler for %s
		http.Handle("%s/", http.StripPrefix("%s/", %s))
		http.Handle("%s", http.StripPrefix("%s", %s))`,
		m.Name, handlerPath, handlerPath, ss.controllerVar,
		handlerPath, handlerPath, ss.controllerVar)

	routes = indexRoute + fmt.Sprintln()
	return
}
