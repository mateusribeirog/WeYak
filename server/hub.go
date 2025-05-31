package main

import "sync"

type Hub struct {
	rooms      map[string]*Room
	register   chan *RegistrationInfo
	unregister chan *Client
	broadcast  chan *Message
	roomsMutex sync.RWMutex
}

type RegistrationInfo struct {
	client          *Client
	desiredRoomName string
}

func NewHub() *Hub {
	h := &Hub{
		rooms:      make(map[string]*Room),
		register:   make(chan *RegistrationInfo),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
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
