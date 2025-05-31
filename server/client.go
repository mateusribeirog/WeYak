package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	conn     net.Conn
	hub      *Hub
	room     *Room
	name     string
	outgoing chan Message
}

func newClient(conn net.Conn, h *Hub) *Client {
	return &Client{
		conn:     conn,
		hub:      h,
		outgoing: make(chan Message),
	}
}

func (c *Client) Disconnect() {

}

func (c *Client) ReadFromClient() {

}

func (c *Client) SendToClient() {

}

func HandleConnection(conn net.Conn, h *Hub) {

	client := newClient(conn, h)
	defer func() {

		close(client.outgoing)
		client.conn.Close()
	}()

	_, err := client.conn.Write([]byte("Welcome to WeYak, please enter your name: "))
	if err != nil {
		fmt.Println(err)
		return
	}

	reader := bufio.NewReader(client.conn)
	name, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println(err)
		return
	}

	client.name = strings.TrimSpace(name)

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
	desiredroom = strings.TrimSpace(desiredroom)
	h.rooms[desiredroom].AddClient(client)

	regRequest := &RegistrationInfo{
		client:          client,
		desiredRoomName: desiredroom,
	}
	h.register := <-regRequest
}
