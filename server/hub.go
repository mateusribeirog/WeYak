package main

import (
	"fmt"
)

type Hub struct {
	rooms      map[string]*Room
	register   chan *RegistrationInfo
	unregister chan *Client
	broadcast  chan Message
}

type RegistrationInfo struct {
	client          *Client
	desiredRoomName string
	name            string
}

func (h *Hub) Run() {
	for {
		select {
		case regInfo := <-h.register:

			regInfo.client.name = regInfo.name

			room, exists := h.rooms[regInfo.desiredRoomName]
			if !exists {
				msg := Message{
					Name:    "SERVER",
					Content: "Room doesn't exist",
				}
				regInfo.client.send <- msg
				break
			}

			err := room.AddClient(regInfo.client)
			if err != nil {
				msg := Message{
					Name:    "SERVER",
					Content: "Failed to join room: " + err.Error(),
				}
				regInfo.client.send <- msg
				break
			}

			regInfo.client.room = room
			fmt.Printf("\n%s joined the %s room\n", regInfo.client.name, regInfo.client.room.name)
			fmt.Println("Local Address: ", regInfo.client.conn.LocalAddr())
			fmt.Println("Remote Adress: ", regInfo.client.conn.RemoteAddr())

		case client := <-h.unregister:
			if client.room != nil {
				room := client.room
				msg := Message{
					Name:    "SERVER",
					Content: fmt.Sprintf("%s has left the room", client.name),
				}
				close(client.send)
				room.BroadcastMessage(msg)
				room.RemoveClient(client.name)
				fmt.Printf("\n%s left the %s room\n", client.name, client.room.name)

			}
		case message := <-h.broadcast:
			room := h.rooms[message.RoomName]
			room.BroadcastMessage(message)
			fmt.Printf("\n%s sent a message in the %s room\n", message.Name, message.RoomName)
			fmt.Println("Content: ", message.Content)
		}
	}

}

func NewHub() *Hub {
	h := &Hub{
		rooms:      make(map[string]*Room),
		register:   make(chan *RegistrationInfo),
		unregister: make(chan *Client),
		broadcast:  make(chan Message),
	}

	roomNames := []string{"Games", "Tech", "Movies", "General"}

	for _, room := range roomNames {
		h.rooms[room] = NewRoom(room)
	}
	return h
}
