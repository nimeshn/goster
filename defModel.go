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

type Model struct {
	Name             string `json:"name"`
	DisplayName      string
	Fields           []*Field
	Options          ViewOptions
	IncludeTimeStamp bool
	IsRestricted     bool
	HasMany          string
	BelongsTo        string
	ViewType         IndexView
	appRef           *App
}

func (m *Model) AddField(field *Field) {
	fld := m.GetField(field.Name)
	if fld != nil {
		errors.New(fmt.Sprintf("Add Failed as field %s exists already in the model %s", field.Name, m.Name))
		return
	}
	m.Fields = append(m.Fields, field)
}

func (m *Model) ValidateRelations() (valid bool, err error) {
	if m.HasMany != "" {
		mod := m.appRef.GetModel(m.HasMany)
		valid = (mod != nil && mod.BelongsTo == m.HasMany)
		if !valid {
			err = errors.New(fmt.Sprintf("%s has many relation %s does not exist", m.Name, m.HasMany))
		}
	}
	if m.BelongsTo != "" {
		mod := m.appRef.GetModel(m.BelongsTo)
		valid = (mod != nil && mod.HasMany == m.BelongsTo)
		if !valid {
			err = errors.New(fmt.Sprintf("%s belongs to relation %s does not exist", m.Name, m.BelongsTo))
		}
	}
	return
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

func NewModel() *Model {
	return &Model{Fields: []*Field{}}
}
