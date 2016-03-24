package main

import (
	"errors"
	"fmt"
)

type ViewOptions struct {
	Index  bool
	Create bool
	View   bool
	Edit   bool
	Remove bool
}

type FieldValidation struct {
	Required  bool
	Url       bool
	Date      bool
	Email     bool
	Integer   bool
	Float     bool
	MinLen    int
	MaxLen    int
	MinValue  int
	MaxValue  int
	BelongsTo string
}

type Field struct {
	Name        string
	DisplayName string
	Validator   *FieldValidation
	ShowInList  bool
}

type Model struct {
	Name             string
	DisplayName      string
	Fields           []*Field
	Options          ViewOptions
	IncludeTimeStamp bool
	IsRestricted     bool
}

type App struct {
	Name        string
	DisplayName string
	CompanyName string
	Models      []*Model
}

func (a *App) AddModel(model *Model) {
	mod := a.GetModel(model.Name)
	if mod != nil {
		errors.New(fmt.Sprintf("Add Failed as model %s exists already in the App", model.Name))
		return
	}
	a.Models = append(a.Models, model)
}

func (a *App) DeleteModel(name string) {
	for index, val := range a.Models {
		if val.Name == name {
			a.Models = append(a.Models[:index], a.Models[index+1:]...)
			return
		}
	}
}

func (a *App) GetModel(name string) *Model {
	for _, val := range a.Models {
		if val.Name == name {
			return val
		}
	}
	return nil
}

func (a *App) SaveModel(model *Model) {
	for index, mod := range a.Models {
		if mod.Name == model.Name {
			a.Models[index] = model
			return
		}
	}
}

func (m *Model) AddField(field *Field) {
	fld := m.GetField(field.Name)
	if fld != nil {
		errors.New(fmt.Sprintf("Add Failed as field %s exists already in the model %s", field.Name, m.Name))
		return
	}
	m.Fields = append(m.Fields, field)
}

func (m *Model) DeleteField(name string) {
	for index, val := range m.Fields {
		if val.Name == name {
			m.Fields = append(m.Fields[:index], m.Fields[index+1:]...)
			return
		}
	}
}

func (m *Model) GetField(name string) *Field {
	for _, val := range m.Fields {
		if val.Name == name {
			return val
		}
	}
	return nil
}

func (m *Model) SaveField(field *Field) {
	for index, fld := range m.Fields {
		if fld.Name == field.Name {
			m.Fields[index] = field
			return
		}
	}
}
