package main

import "fmt"

type Message struct {
	Name, Content, RoomName string
}

func FormatMessage(message Message) string {

	var formattedOutput string

	if message.Name == "SERVER" {
		formattedOutput = fmt.Sprintf("-%s-\n", message.Content)
	} else {
		formattedOutput = fmt.Sprintf("[%s] %s: %s\n", message.RoomName, message.Name, message.Content)
	}
	return formattedOutput
}
