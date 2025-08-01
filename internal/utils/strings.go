package utils

import (
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func TitleString(s string) string {
	caser := cases.Title(language.AmericanEnglish)
	return string(caser.String(s))
}
