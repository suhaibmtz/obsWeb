package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"
)

var client *redis.Client
var ctx = context.Background()

func init() {
	client = redis.NewClient(&redis.Options{Addr: ":6379"})
}

func GetCommandsList() (commands []command) {
	for key, val := range client.HGetAll(ctx, "obsWeb").Val() {
		val := strings.ReplaceAll(val, "'", `"`)
		var record command
		record.Name = key
		record.Command = val
		commands = append(commands, record)
	}
	return
}

func Set(name, value string) error {
	value = strings.ReplaceAll(value, `"`, `'`)
	fmt.Println(value)
	//{'command':'PlayVideo','params':['Scene','Memes/Hello There']}
	return client.HSet(ctx, "obsWeb", name, value).Err()
}

func Get(key string) command {
	val := client.HGet(ctx, "obsWeb", key).Val()
	val = strings.ReplaceAll(val, "'", `"`)
	var record command
	record.Name = key
	record.Command = val
	return record
}
