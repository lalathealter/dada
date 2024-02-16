package controllers

import "regexp"

func isJustLatin(name string) bool {
	nonLatinReg, _ := regexp.Compile("[^A-Za-z]")
	gotAny := nonLatinReg.Find([]byte(name))
	return gotAny == nil
}
