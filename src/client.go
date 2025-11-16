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
	room.mu.Lock()
	defer room.mu.Unlock()
	roomsMu.Lock()
	if existing, ok := rooms[room.code]; ok && existing == room {
		delete(rooms, room.code)
	}
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
			fmt.Println("Read error:", err)
			return
		}

		if mt == websocket.TextMessage {
			var obj map[string]any
			if err := json.Unmarshal(msg, &obj); err == nil {
				fmt.Println("Error parsing json")
				Player.WriteJson(map[string]any{
					"type":   "parse_error",
					"detail": "invalid_json",
				})
				continue
			}
		}

		if Player.room != nil {
			room := Player.room
			room.mu.Lock()
			other := room.other(Player)
			if other != nil {
				other.Write(mt, msg)
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
