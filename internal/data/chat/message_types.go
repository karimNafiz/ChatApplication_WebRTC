package chat

type MessageType int

/*
	represents all the different types of messages that can be sent

*/

const (
	MessageText MessageType = iota
	MessageImage
	MessageFile
	MessageVideo
	MessageAudio
)
