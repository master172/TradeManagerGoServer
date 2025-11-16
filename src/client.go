package main

import (
	"encoding/json"
	"fmt"
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

func (Player *Client) ReadLoop() {
	defer func() {
		Player.Cleanup()
		Player.conn.Close()
	}()
	for {
		mt, msg, err := Player.conn.ReadMessage()
		if err != nil {
			return
		}

		if mt == websocket.TextMessage {
			var obj map[string]any
			if err := json.Unmarshal(msg, &obj); err == nil {
				fmt.Println("Error parsing json")
				continue
			}
		}

		if Player.room != nil {
			if other := Player.room.other(Player); other != nil {
				other.WriteJson(msg)
			} else {
				Player.WriteJson(map[string]any{
					"type": "no_partner",
				})
			}
		} else {
			Player.WriteJson(map[string]any{
				"type": "not_in_room",
			})
		}
	}
}
