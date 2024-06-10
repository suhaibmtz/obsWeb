package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

// todo SetText

var commands = map[string]func(params ...any){
	"SetShow":    hideShowItem,
	"ToggleShow": ToggleItemVisibility,
	"Wait":       Wait,
	"PlayVideo":  PlayVideo,
}

type runCommand struct {
	Command string
	Params  []any
}

type command struct {
	Name    string
	Command string
}

func HandleCommand(message string) {
	for _, CmdString := range strings.Split(message, "\n") {
		var Cmd runCommand
		err := json.Unmarshal([]byte(CmdString), &Cmd)
		if err != nil {
			WriteLog("Error in HandleCommand: ", err.Error())
			break
		}
		if cm, ok := commands[Cmd.Command]; ok {
			cm(Cmd.Params...)
		}
	}
}

func Wait(params ...any) {
	time.Sleep(time.Second * time.Duration(params[0].(float64)))
}

func hideShowItem(params ...any) {
	c := connect()
	defer c.Close()
	hideShowItemConn(c, params...)
}

func hideShowItemConn(c *websocket.Conn, params ...any) {
	val := fmt.Sprintf(`{
		"request-type": "SetSceneItemProperties",
		"message-id": "show_hide_%s",
		"scene-name": "%s",
		"item": "%s",
		"visible": %s
	}`, params[1].(string), params[0].(string), params[1].(string), strconv.FormatBool(params[2].(bool)))
	err := c.WriteMessage(1, []byte(val))
	if err != nil {
		WriteLog("Error: ", err.Error())
		return
	}
	_, resp, err := c.ReadMessage()
	if err == nil {
		WriteLog(string(resp))
	} else {
		WriteLog("Error: ", err.Error())
	}
	_, resp, err = c.ReadMessage()
	if err == nil {
		WriteLog(fmt.Sprint(params[3:]...), string(resp))
	} else {
		WriteLog("Error at ShowHideItem: ", err.Error())
	}
}

func PlayVideo(params ...any) {
	// hideShowItem(append(params, true)...)
	c := connect()
	val := fmt.Sprintf(`{
		"request-type": "RestartMedia",
		"message-id": "play_%s",
		"scene-name": "%s",
        "sourceName": "%s"
	}`, params[1].(string), params[0].(string), params[1].(string))
	err := c.WriteMessage(1, []byte(val))
	if err != nil {
		WriteLog("Error: ", err.Error())
		return
	}
	_, resp, err := c.ReadMessage()
	if err == nil {
		WriteLog(string(resp))
	} else {
		WriteLog("Error: ", err.Error())
	}
}

type TransformType struct {
	Visible bool
}

type GetTransform struct {
	TransformType `json:"transform"`
}

func ToggleItemVisibility(params ...any) {
	c := connect()
	defer c.Close()
	Scene := params[0].(string)
	Item := params[1].(string)
	val := fmt.Sprintf(`{
		"request-type": "SetSceneItemProperties",
		"message-id": "toggle_%s",
		"scene-name": "%s",
		"item": "%s"
	}`, Item, Scene, Item)
	c.WriteMessage(1, []byte(val))
	_, resp, err := c.ReadMessage()
	if err == nil {
		var trans GetTransform
		err := json.Unmarshal(resp, &trans)
		if err != nil {
			fmt.Println("ERROR: ", err.Error())
		}
		hideShowItemConn(c, Scene, Item, !trans.Visible, "Toggle")
	} else {
		WriteLog("Error at ToggleItemVisibility: ", err.Error())
	}
}

func connect() *websocket.Conn {
	var err error
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:4444", nil)
	if err != nil {
		WriteLog(err)
	}
	return c
}

func CommandsList() []command {
	return GetCommandsList()
}

func SetCommand(name, value string) (success bool, message string) {
	err := Set(name, value)
	success = true
	if err != nil {
		message = "Error: " + err.Error()
		WriteLog("Error in Set: ", err.Error())
		success = false
	}
	return
}
