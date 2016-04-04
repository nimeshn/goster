package main

import (
	"fmt"
	"strings"
)

func (fld *Field) GetServerValidation(a *ServerModelSettings) (goCode string) {
	if fld.Validator == nil {
		return
	}
	fieldName := fmt.Sprintf(`m.%s`, strings.Title(fld.Name))
	goCode += fmt.Sprintf(`//Validate %s`, fld.DisplayName)
	if fld.Validator.MinLen > 0 {
		goCode += fmt.Sprintln() + fmt.Sprintf(
			`if %s!="" && len(%s)<%d {
				modelErrors = append(modelErrors, "%s should have atleast %d characters.")
			}`, fieldName, fieldName, fld.Validator.MinLen, fld.DisplayName, fld.Validator.MinLen)
	}
	if fld.Validator.MaxLen > 0 {
		goCode += fmt.Sprintln() + fmt.Sprintf(
			`if %s!="" && len(%s)>%d{
				modelErrors = append(modelErrors, "%s should have atmost %d characters.")
			}`, fieldName, fieldName, fld.Validator.MaxLen, fld.DisplayName, fld.Validator.MaxLen)
	}
	if fld.Validator.MinValue > 0 {
		goCode += fmt.Sprintln() + fmt.Sprintf(
			`if %s!=nil && %s<%d{
				modelErrors = append(modelErrors, "%s should be greater than or equal to %d.")
			}`, fieldName, fieldName, fld.Validator.MinValue, fld.DisplayName, fld.Validator.MinValue)
	}
	if fld.Validator.MaxValue > 0 {
		goCode += fmt.Sprintln() + fmt.Sprintf(
			`if %s!=nil && %s>%d{
				modelErrors = append(modelErrors, "%s should be less than or equal to %d.")
			}`, fieldName, fieldName, fld.Validator.MaxValue, fld.DisplayName, fld.Validator.MaxValue)
	}
	if fld.Validator.Email {
		goCode += fmt.Sprintln() + fmt.Sprintf(
			`if %s!="" && !ValidateEmail(%s){
				modelErrors = append(modelErrors, "%s should be a valid Email Id.")
			}`, fieldName, fieldName, fld.DisplayName)
	}
	if fld.Validator.Url {
		goCode += fmt.Sprintln() + fmt.Sprintf(
			`if %s!="" && !ValidateUrl(%s){
				modelErrors = append(modelErrors, "%s should be a valid Url.")
			}`, fieldName, fieldName, fld.DisplayName)
	}
	if fld.Validator.IsAlpha {
		goCode += fmt.Sprintln() + fmt.Sprintf(
			`if %s!="" && !IsAlpha(%s){
				modelErrors = append(modelErrors, "%s should only contains alphabets.")
			}`, fieldName, fieldName, fld.DisplayName)
	}
	if fld.Validator.IsAlphaNumeric {
		goCode += fmt.Sprintln() + fmt.Sprintf(
			`if %s!="" && !IsAlphaNumeric(%s){
				modelErrors = append(modelErrors, "%s should only contains alphabets or numbers.")
			}`, fieldName, fieldName, fld.DisplayName)
	}
	if fld.Validator.Required {
		goCode += fmt.Sprintln() + fmt.Sprintf(
			`if %s==""{
				modelErrors = append(modelErrors, "%s is required.")
			}`, fieldName, fld.DisplayName)
	}
	return
}
