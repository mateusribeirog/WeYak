package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"sync"
)

const bufferSize = 256

type Client struct {
	conn      net.Conn
	hub       *Hub
	room      *Room
	name      string
	send      chan Message
	connected bool
	mutex     sync.Mutex
}

func newClient(conn net.Conn, h *Hub) *Client {
	return &Client{
		conn:      conn,
		hub:       h,
		send:      make(chan Message, bufferSize),
		connected: true,
		mutex:     sync.Mutex{},
	}
}

func (c *Client) Disconnect() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if !c.connected {
		return
	}
	c.connected = false
	close(c.send)
	c.conn.Close()
	if c.name != "" {
		c.hub.unregister <- c
	}
}

func (c *Client) readFromClient() {
	defer c.Disconnect()
	scanner := bufio.NewScanner(c.conn)
	for scanner.Scan() { //scanner only return false when the user close the connection, e. g.: sent EOF
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		if c.room == nil {
			continue
		}
		msg := Message{
			Name:     c.name,
			Content:  text,
			RoomName: c.room.name,
		}
		c.hub.broadcast <- msg
	}
}

func (c *Client) sendToClient() {
	defer c.Disconnect()
	for msg := range c.send {
		formattedMsg := formatMessage(msg)
		_, err := c.conn.Write([]byte(formattedMsg))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func HandleConnection(conn net.Conn, h *Hub) {

	client := newClient(conn, h)

	_, err := client.conn.Write([]byte("Welcome to WeYak, please enter your name: "))
	if err != nil {
		fmt.Println(err)
		return
	}

	reader := bufio.NewReader(client.conn)
	username, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = client.conn.Write([]byte("What room do you wanna join: "))
	if err != nil {
		fmt.Println(err)
		return
	}

	rooms := h.ListRooms()
	roomList := strings.Join(rooms, "\n")
	_, err = client.conn.Write([]byte(roomList))
	if err != nil {
		fmt.Println(err)
		return
	}

	desiredroom, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return
	}

	regRequest := &RegistrationInfo{
		client:          client,
		desiredRoomName: strings.TrimSpace(desiredroom),
		name:            username,
	}
	h.register <- regRequest

	go client.sendToClient()
	client.readFromClient()
}
