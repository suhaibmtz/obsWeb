package main

import (
	"fmt"

	"github.com/motaz/codeutils"
)

func WriteLog(event ...any) {
	fmt.Print("Log:")
	fmt.Println(event...)
	err := codeutils.WriteToLog(fmt.Sprint(event...), "obsWeb")
	if err != nil {
		fmt.Println("WriteLog Error: " + err.Error())
	}
}

func GetConfigValue(param, def string) (value string) {

	value = codeutils.GetConfigValue("config.ini", param)

	if value == "" {
		value = def
	}
	return
}

func SetConfigValue(param, value string) (success bool) {

	success = codeutils.SetConfigValue("config.ini", param, value)
	return
}

func GetMD5(text string) string {
	return codeutils.GetMD5(text)
}
