package model

type Message struct {
	// The action to take
	Message string `json:"message"`
	// The contents of the message
	Data string `json:"data"`
}
