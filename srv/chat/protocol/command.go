package protocol

import "errors"

// This package is inspired by nqbao-gosandbox
var (
	ErrUnknownCommand = errors.New("unknown command")
)

// SendCommand is for sending the message.
type SendCommand struct {
	Message string
}

// MessageCommand is for receiving the message.
type MessageCommand struct {
	Name    string
	Message string
}

// NameCommand is used for setting client display name
type NameCommand struct {
	Name string
}
