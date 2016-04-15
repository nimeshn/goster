package main

import (
	"fmt"
	"strconv"
	"strings"
)

type DBGenerator interface {
	GetFileName(a *App) (fileName string)
	CreateDB(a *App) (sql string, err error)
	CreateTable(m *Model) (sql string, err error)
	CreateColumn(f *Field) (sql string, err error)
}

type MySQLDBGenerator struct {
	Name string
}

func (m *MySQLDBGenerator) GetFileName(a *App) (fileName string) {
	a.Name
}
