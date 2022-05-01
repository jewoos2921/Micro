package main

import (
	"github.com/gorilla/websocket"
	"time"
)

// client 는 한명의 채팅 사용자를 나타낸다.
type client struct {
	// socket은 이 클라이언트의 웹 소켓이다.
	socket *websocket.Conn
	// send는 메시지가 전송되는 채널이다.
	send chan *message
	// room은 클라이언트가 채팅방이다.
	room *room
	// 사용자에 대한 정보를 보유한다.
	userData map[string]interface{}
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		var msg *message
		if err := c.socket.ReadJSON(&msg); err != nil {
			return
		}

		msg.When = time.Now()
		msg.Name = c.userData["name"].(string)
		if avatarUrl, ok := c.userData["avatar_url"]; ok {
			msg.AvatarURL = avatarUrl.(string)
		}

		c.room.forward <- msg
	}
}
func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.send {
		err := c.socket.WriteJSON(msg)
		if err != nil {
			break
		}
	}
}
