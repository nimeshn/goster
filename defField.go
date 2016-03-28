package main

type FieldValidation struct {
	IsAlpha        bool
	IsAlphaNumeric bool
	Url            bool
	Email          bool
	Unique         bool
	Required       bool
	MinLen         int
	MaxLen         int
	MinValue       int
	MaxValue       int
}

type Field struct {
	Name          string
	DisplayName   string
	ShowInIndex   bool
	AutoGenerated bool
	Type          Types
	Validator     *FieldValidation
}

func NewField() *Field {
	return &Field{Validator: &FieldValidation{}}
}
