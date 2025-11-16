package main

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn   *websocket.Conn
	room   *Room
	code   string
	sendMu sync.Mutex
}

func (c *Client) Write(mt int, msg []byte) error {
	c.sendMu.Lock()
	defer c.sendMu.Unlock()
	return c.conn.WriteMessage(mt, msg)
}
func (c *Client) WriteJson(obj any) error {
	c.sendMu.Lock()
	defer c.sendMu.Unlock()
	return c.conn.WriteJSON(obj)
}

func (c *Client) Cleanup() {
	if c.room == nil {
		return
	}

	room := c.room
	roomsMu.Lock()
	delete(rooms, room.code)
	roomsMu.Unlock()

	if other := room.other(c); other != nil {
		other.WriteJson(map[string]any{
			"type": "partner_disconnected",
		})
		other.room = nil
	}
	c.room = nil
	c.code = ""
}
