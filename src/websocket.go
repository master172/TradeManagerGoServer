package main

import (
	"fmt"
	"net/http"
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Eroor upgrading", err)
		return
	}
	client := &Client{conn: conn}
	go client.ReadLoop()
}
