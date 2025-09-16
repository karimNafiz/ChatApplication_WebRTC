package chat

/*
	the interface needs to be implemented by all the message types
*/

type Sendable interface {
	GetType() MessageType
}
