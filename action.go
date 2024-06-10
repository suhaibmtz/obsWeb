package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

func socket(w http.ResponseWriter, r *http.Request) {
	con, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		WriteLog(err.Error())
	} else {
		for {
			_, message, err := con.ReadMessage()
			if err != nil {
				fmt.Println("Conn err: ", err.Error())
				break
			}
			go HandleCommand(string(message))
		}
	}
}

type IndexData struct {
	Commands []command
}

func Index(w http.ResponseWriter, r *http.Request) {
	var Data IndexData
	Data.Commands = CommandsList()
	err := temp.ExecuteTemplate(w, "index.html", Data)
	if err != nil {
		fmt.Println(err.Error())
	}
}

type SettingsData struct {
	Commands []command
}

func settings(w http.ResponseWriter, r *http.Request) {
	var Data SettingsData
	Data.Commands = CommandsList()
	err := temp.ExecuteTemplate(w, "settings.html", Data)
	if err != nil {
		fmt.Println(err.Error())
	}
}

type SetData struct {
	ShowName bool
	Value    string
}

func SetPage(w http.ResponseWriter, r *http.Request) {
	var Data SetData
	if r.FormValue("command") == "" {
		Data.ShowName = true
	} else {
		Data.Value = Get(r.FormValue("command")).Command
	}
	if r.FormValue("set") != "" {
		key := r.FormValue("name")
		if key == "" {
			key = r.FormValue("command")
		}
		val := r.FormValue("value")
		success, err := SetCommand(key, val)
		if success {
			http.Redirect(w, r, "/settings", http.StatusTemporaryRedirect)
			return
		} else {
			fmt.Fprintf(w, "<p>%s</p>", err)
		}
	}
	err := temp.ExecuteTemplate(w, "set.html", Data)
	if err != nil {
		fmt.Println(err.Error())
	}
}
