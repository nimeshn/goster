package main

import (
	"fmt"
	"path"
	"strings"
)

type DBGenerator interface {
	GetFileName(a *App) (fileName string)
	CreateDB(a *App) (sql string)
	CreateTable(m *Model) (sql string)
	CreateColumn(f *Field) (sql string)
}

type MySQLDBGenerator struct {
	Name string
}

func (g *MySQLDBGenerator) GetFileName(a *App) (fileName string) {
	t := a.GetServerSettings()
	fileName = fmt.Sprintf("%s.%s.sql", path.Join(t.directories["db"], a.Name), strings.ToLower(g.Name))
	return
}

func (g *MySQLDBGenerator) CreateDB(a *App) (sql string) {
	t := a.GetServerSettings()
	sql = fmt.Sprintf(
		`DROP DATABASE IF EXISTS %s;
		CREATE DATABASE %s;
		GRANT ALL ON sampleapp.* TO '%s'@'localhost' IDENTIFIED BY '%s';`,
		t.dbName, t.dbName, t.dbUser, t.dbUserPassword)
	return
}

/*
func (g *MySQLDBGenerator) CreateTable(m *Model) (sql string) {
	t := m.GetServerSettings()
	sql = fmt.Sprintf(
		`DROP DATABASE IF EXISTS %s;
		CREATE DATABASE %s;
		GRANT ALL ON sampleapp.* TO '%s'@'localhost' IDENTIFIED BY '%s';`,
		t.dbName, t.dbName, t.dbUser, t.dbUserPassword)
	return
}

func (g *MySQLDBGenerator) CreateColumn(f *Field) (sql string) {
	t := a.GetServerSettings()
	sql = fmt.Sprintf(
		`DROP DATABASE IF EXISTS %s;
		CREATE DATABASE %s;
		GRANT ALL ON sampleapp.* TO '%s'@'localhost' IDENTIFIED BY '%s';`,
		t.dbName, t.dbName, t.dbUser, t.dbUserPassword)
	return
}
*/
