package main

import (
	"regexp"
)

const longTimeForm = "Apr 3, 2016 at 1:41pm (IST)"

func ValidateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,3}$`)
	return Re.MatchString(email)
}

func ValidateUrl(email string) bool {
	Re := regexp.MustCompile(`https?://.+`)
	return Re.MatchString(email)
}

func IsAlpha(val string) bool {
	Re := regexp.MustCompile(`^[a-zA-Z]+$`)
	return Re.MatchString(val)
}

func IsAlphaNumeric(val string) bool {
	Re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	return Re.MatchString(val)
}
