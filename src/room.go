package main

type Room struct {
	code    string
	Player1 *Client
	Player2 *Client
}

func (r *Room) other(c *Client) *Client {
	if c == r.Player1 {
		return r.Player2
	}
	if c == r.Player2 {
		return r.Player1
	}
	return nil
}
