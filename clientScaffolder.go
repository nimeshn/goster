package main

import (
	"fmt"
)

func (m *Model) GetEditViewClient() (fileName, htmlCode string) {
	fileName = m.Name + "Edit.view.htm"
	modelData := m.Name + "Data"
	modelForm := m.Name + "Form"
	modelNewFunc := fmt.Sprintf("isNew%s()", m.Name)

	htmlCode = `<div class="row" ng-include="'/app/error/errorhandler.view.html'"></div><div class="row text-center">`
	htmlCode += fmt.Sprintf(`<h3 ng-bind="(%s?'New':'Edit') + ' %s'"></h3>`, modelNewFunc, m.DisplayName)
	htmlCode += `<hr/></div>`

	htmlCode += fmt.Sprintf(`<form class="form-horizontal" role="form" name="%s" ng-submit="submit()"><div class="row">`, modelForm)
	for _, fld := range m.Fields {
		htmlCode += `<div class="form-group">`
		htmlCode += fmt.Sprintf(`<label class="control-label col-sm-4" for="%s">%s:</label>`, fld.Name, fld.DisplayName)
		htmlCode += `<div class="col-sm-8">`
		htmlCode += fmt.Sprintf(`<input type="text" class="form-control" id="%s" placeholder="Enter %s" ng-model="%s.%s" title="%s" `,
			fld.Name, fld.DisplayName, modelData, fld.Name, fld.DisplayName)

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

func (m *Model) GetShowViewClient() (fileName, htmlCode string) {
	fileName = m.Name + "Show.view.htm"
	modelData := m.Name + "Data"

	htmlCode = fmt.Sprintf(`<div class="row" ng-include="'/app/error/errorhandler.view.html'"></div>`+
		`<div class="row text-center"><h3>%s Details</h3><hr></div>`, m.DisplayName)

	htmlCode += `<div class="row"><div class="col-sm-12">`
	for _, fld := range m.Fields {
		htmlCode += fmt.Sprintf(`<div class="row"><div class="col-sm-12"><h3 ng-bind="%s"></h3><p ng-bind="%s.%s"></p></div></div>`,
			fld.DisplayName, modelData, fld.Name)
	}
	htmlCode += `</div></div>`
	htmlCode += `<div class="row"><hr></div>`
	return
}

func (m *Model) GetIndexViewClient() (fileName, htmlCode string) {
	fileName = m.Name + "Index.view.htm"
	modelData := m.Name + "List"
	modelDeleteFunc := "Delete" + m.Name
	modelLoadFunc := "Load" + m.Name

	htmlCode = fmt.Sprintf(`<div class="row" ng-include="'/app/error/errorhandler.view.html'"></div>`+
		`<div class="row text-center"><h3>%s List</h3><hr></div>`, m.DisplayName)

	htmlCode += fmt.Sprintf(`<div class="row text-center"><a href="" ng-click="%s()"><span class="glyphicon glyphicon-refresh"/> Refresh</a>`+
		`<a href="#%s/new" class="col-sm-offset-1"><span class="glyphicon glyphicon-plus"/> New %s</a></div><br/>`,
		modelLoadFunc, m.Name, m.DisplayName)

	htmlCode += fmt.Sprintf(`<div ng-if="%s.length==0" class="row"><div class="col-sm-12 text-center"><h3>0 Records Found.</h3></div></div>`,
		modelData)

	htmlCode += fmt.Sprintf(`<div ng-if="%s.length>0" class="row"><div class="col-sm-12">`, modelData)
	if m.ViewType == List {
		htmlCode += fmt.Sprintf(`<div class="row" ng-repeat="x in %s | orderBy:createdOn:reverse">`, modelData)
		htmlCode += fmt.Sprintf(`<div class="col-sm-1"><a href="/#%s/view/{{x.id}}" alt="View %s" title="View %s">`+
			`<span class="glyphicon glyphicon-folder-open"></span></a></div>`,
			m.Name, m.Name, m.Name)
		htmlCode += fmt.Sprintf(`<div class="col-sm-1"><a href="/#%s/edit/{{x.id}}" alt="Edit %s" title="Edit %s">`+
			`<span class="glyphicon glyphicon-edit"></span></a></div>`,
			m.Name, m.Name, m.Name)
		htmlCode += fmt.Sprintf(`<div class="col-sm-1"><a href="" alt="Delete %s" title="Delete %s"><span class="glyphicon glyphicon-remove" `+
			`ng-click="%s(x.Id);"></span></a></div>`,
			m.Name, m.Name, modelDeleteFunc)
		for _, fld := range m.Fields {
			if fld.ShowInIndex {
				htmlCode += fmt.Sprintf(`<div class="col-sm-1"><span ng-bind="%s.%s"></span></div>`,
					modelData, fld.Name)
			}
		}
		htmlCode += `</div>`
	}
	htmlCode += `</div></div>`
	htmlCode += `<div class="row"><hr></div>`
	return
}
