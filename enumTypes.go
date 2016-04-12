package main

// A types specifies a data type allowed as a field
type Types int

const (
	String Types = iota
	Integer
	Float
	Date
	Boolean
)

var typesDef = [...]string{
	"String",
	"Integer",
	"Float",
	"Date",
	"Boolean",
}

// String returns the desc of the Type
func (t Types) String() string { return typesDef[t] }

type IndexView int

const (
	Table IndexView = iota
	List
)

var indexViewDef = [...]string{
	"Table",
	"List",
}

// String returns the desc of the Type
func (i IndexView) String() string { return indexViewDef[i] }
