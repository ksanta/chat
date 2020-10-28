package model

type Message struct {
	// The action to take
	Type string `json:"type"`
	// The contents of the message
	Data string `json:"data"`
	// Name of the room user is chatting in, blank for the global room
	RoomName string `json:"roomName"`
}
