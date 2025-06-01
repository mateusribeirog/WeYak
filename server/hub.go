package main

import (
	"fmt"
	"sync"
)

type Hub struct {
	rooms      map[string]*Room
	register   chan *RegistrationInfo
	unregister chan *Client
	broadcast  chan Message
	roomsMutex sync.RWMutex
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

		case client := <-h.unregister:
			if client.room != nil {
				room := client.room
				msg := Message{
					Name:    "SERVER",
					Content: fmt.Sprintf("%s has left the room", client.name),
				}
				room.BroadcastMessage(msg)
				room.RemoveClient(client.name)
			}
		case message := <-h.broadcast:
			room := h.rooms[message.RoomName]
			room.BroadcastMessage(message)
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

func (h *Hub) ListRooms() []string {
	h.roomsMutex.RLock()
	defer h.roomsMutex.RUnlock()

	names := make([]string, 0, len(h.rooms))
	for room := range h.rooms {
		names = append(names, room)
	}
	return names
}
