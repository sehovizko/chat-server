package main

import "time"

type message struct {
	AvatarURL string
	Name      string
	Message   string
	When      time.Time
}
