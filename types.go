package main

// A types specifies a data type allowed as a field
type Types int

const (
	String Types = iota + 1
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
func (t Types) String() string { return typesDef[t-1] }

type IndexView int

const (
	List IndexView = iota + 1
	Table
)

var indexViewDef = [...]string{
	"List",
	"Table",
}

// String returns the desc of the Type
func (i IndexView) String() string { return indexViewDef[i-1] }
