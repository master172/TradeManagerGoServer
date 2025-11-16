package main

import (
	"sync"
)

var (
	rooms   = make(map[string]*Room)
	roomsMu sync.Mutex
)

func CreateRoom(Player *Client) {
	code := GenCode()
	room := &Room{
		code:    code,
		Player1: Player,
		Player2: nil,
	}
	Player.room = room
	Player.code = code

	roomsMu.Lock()
	rooms[code] = room
	roomsMu.Unlock()

	Player.WriteJson(map[string]any{
		"type": "room_created",
		"code": code,
	})
}

func JoinRoom(Player *Client, code string) {
	roomsMu.Lock()
	room, ok := rooms[code]

	if !ok {
		Player.WriteJson(map[string]any{
			"type":   "join_failed",
			"reason": "Room not found",
		})
	} else if room.Player2 != nil {
		Player.WriteJson(map[string]any{
			"type":   "join_failed",
			"reason": "room_full",
		})
	}

	room.Player2 = Player
	Player.room = room
	Player.code = code

	room.Player1.WriteJson(map[string]any{
		"type": "paired",
		"code": "code",
	})

}
