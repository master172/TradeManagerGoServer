package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Room struct {
	code string
	a    *Client
	b    *Client
}
type Client struct {
	conn   *websocket.Conn
	room   *Room
	code   string
	sendMu sync.Mutex
}

var (
	rooms   = make(map[string]*Room)
	roomsMu sync.Mutex
)

func CreateRoom(c *Client) {
	code := GenCode()
	room := &Room{
		code: code,
		a:    c,
		b:    nil,
	}
	c.room = room
	c.code = code

	roomsMu.Lock()
	rooms[code] = room
	roomsMu.Unlock()

	c.safeWriteJson(map[string]any{
		"type": "room_created",
		"code": code,
	})
}

func (c *Client) safeWrite(mt int, msg []byte) error {
	c.sendMu.Lock()
	defer c.sendMu.Unlock()
	return c.conn.WriteMessage(mt, msg)
}
func (c *Client) safeWriteJson(obj any) error {
	c.sendMu.Lock()
	defer c.sendMu.Unlock()
	return c.conn.WriteJSON(obj)
}
