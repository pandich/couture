package model

// Message ...
// Message a message.
type Message string

// String ...
func (msg Message) String() string {
	return string(msg)
}
