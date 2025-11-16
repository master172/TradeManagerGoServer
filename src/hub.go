package main

import (
	"sync"
)

var (
	rooms   = make(map[string]*Room)
	roomsMu sync.Mutex
)

func CreateRoom(Player *Client) {
	var code string
	for {
		code = GenCode()

		roomsMu.Lock()
		_, exists := rooms[code]
		roomsMu.Unlock()

		if !exists {
			break
		}
	}
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
	roomsMu.Unlock()

	if !ok {
		Player.WriteJson(map[string]any{
			"type":   "join_failed",
			"reason": "Room not found",
		})
		return
	}

	room.mu.Lock()
	defer room.mu.Unlock()

	if room.Player2 != nil {
		Player.WriteJson(map[string]any{
			"type":   "join_failed",
			"reason": "room_full",
		})
		return
	}

	room.Player2 = Player
	Player.room = room
	Player.code = code

	room.Player1.WriteJson(map[string]any{
		"type": "paired",
		"code": code,
	})
	room.Player2.WriteJson(map[string]any{
		"type": "paired",
		"code": code,
	})

}

func GetActiveRoomCodes() []string {
	roomsMu.Lock()
	defer roomsMu.Unlock()
	codes := make([]string, 0, len(rooms))
	for k := range rooms {
		codes = append(codes, k)
	}
	return codes
}
