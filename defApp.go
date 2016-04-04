package main

import (
	"errors"
	"fmt"
)

type App struct {
	Name        string   `json:"name"`
	DisplayName string   `json:"displayName"`
	CompanyName string   `json:"companyName"`
	VersionNo   string   `json:"versionNo"`
	AppDir      string   `json:"appDir"`
	Models      []*Model `json:"models"`
}

func (a *App) AddModel(model *Model) {
	mod := a.GetModel(model.Name)
	if mod != nil {
		errors.New(fmt.Sprintf("Add Failed as model %s exists already in the App", model.Name))
		return
	}
	model.AutoGenerateFields()
	model.appRef = a
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
