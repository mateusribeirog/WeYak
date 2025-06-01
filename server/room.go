package main

import (
	"errors"
	"sync"
)

type Room struct {
	name    string
	clients map[string]*Client
	mutex   sync.Mutex
}

func NewRoom(name string) *Room {
	return &Room{
		name:    name,
		clients: make(map[string]*Client),
	}
}

func (room *Room) AddClient(c *Client) error {
	room.mutex.Lock()
	defer room.mutex.Unlock()
	_, check := room.clients[c.name]
	if check {
		return errors.New("user already exists in this room")
	}
	room.clients[c.name] = c
	return nil
}

func (room *Room) RemoveClient(name string) {
	room.mutex.Lock()
	defer room.mutex.Unlock()
	delete(room.clients, name)
}

func (room *Room) BroadcastMessage(m Message) {
	room.mutex.Lock()
	defer room.mutex.Unlock()
	for name, client := range room.clients {
		if name != m.Name {
			client.send <- m
		}
	}
}
