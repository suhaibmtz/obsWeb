package main

import (
	"fmt"
	"html/template"
	"net/http"
)

var temp *template.Template

func main() {
	temp, _ = template.ParseGlob("templates/*")
	println("http://localhost:8822/")
	http.HandleFunc("/settings", settings)
	http.HandleFunc("/socket", socket)
	http.HandleFunc("/set", SetPage)
	http.HandleFunc("/", Index)
	err := http.ListenAndServe(":8822", nil)
	fmt.Println(err)
}
