package main

import "time"

// 단일 메시지를 나타낸다.
type message struct {
	Name      string
	Message   string
	When      time.Time
	AvatarURL string
}
