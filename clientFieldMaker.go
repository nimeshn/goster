package main

import (
	"fmt"
	"github.com/yosssi/gohtml"
)

func (fld *Field) GetClientValidation(a *ClientModelSettings) (htmlCode string) {
	fieldName := fmt.Sprintf(`$scope.%s.%s`, a.formData, fld.Name)

	if fld.Validator.MinLen > 0 {
		htmlCode += fmt.Sprintln() + fmt.Sprintf(
			`if (%s!=null && %s.length<%d){
				$scope.errors.push('%s should have atleast %d characters.');
			}`, fieldName, fieldName, fld.Validator.MinLen, fld.DisplayName, fld.Validator.MinLen)
	}
	if fld.Validator.MaxLen > 0 {
		htmlCode += fmt.Sprintln() + fmt.Sprintf(
			`if (%s!=null && %s.length>%d){
				$scope.errors.push('%s should have atmost %d characters.');
			}`, fieldName, fieldName, fld.Validator.MaxLen, fld.DisplayName, fld.Validator.MaxLen)
	}
	if fld.Validator.MinValue > 0 {
		htmlCode += fmt.Sprintln() + fmt.Sprintf(
			`if (%s!=null && %s<%d){
				$scope.errors.push('%s should be greater than or equal to %d.');
			}`, fieldName, fieldName, fld.Validator.MinValue, fld.DisplayName, fld.Validator.MinValue)
	}
	if fld.Validator.MaxValue > 0 {
		htmlCode += fmt.Sprintln() + fmt.Sprintf(
			`if (%s!=null && %s>%d){
				$scope.errors.push('%s should be less than or equal to %d.');
			}`, fieldName, fieldName, fld.Validator.MaxValue, fld.DisplayName, fld.Validator.MaxValue)
	}
	if fld.Validator.Email {
		htmlCode += fmt.Sprintln() + fmt.Sprintf(
			`if (%s!=null && !ValidateEmail(%s)){
				$scope.errors.push('%s should be a valid Email Id.');
			}`, fieldName, fieldName, fld.DisplayName)
	}
	if fld.Validator.Url {
		htmlCode += fmt.Sprintln() + fmt.Sprintf(
			`if (%s!=null && !ValidateUrl(%s)){
				$scope.errors.push('%s should be a valid Url.');
			}`, fieldName, fieldName, fld.DisplayName)
	}
	if fld.Validator.IsAlpha {
		htmlCode += fmt.Sprintln() + fmt.Sprintf(
			`if (%s!=null && !IsAlpha(%s)){
				$scope.errors.push('%s should only contains alphabets.');
			}`, fieldName, fieldName, fld.DisplayName)
	}
	if fld.Validator.IsAlphaNumeric {
		htmlCode += fmt.Sprintln() + fmt.Sprintf(
			`if (%s!=null && !IsAlphaNumeric(%s)){
				$scope.errors.push('%s should only contains alphabets or numbers.');
			}`, fieldName, fieldName, fld.DisplayName)
	}
	if fld.Validator.Required {
		htmlCode += fmt.Sprintln() + fmt.Sprintf(
			`if (%s==null || %s==""){
				$scope.errors.push('%s is required.');
			}`, fieldName, fieldName, fld.DisplayName)
	}
	htmlCode = gohtml.Format(htmlCode)
	return
}
