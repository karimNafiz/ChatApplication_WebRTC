package chat

type MessageType int

/*
	represents all the different types of messages that can be sent

*/

const (
	EMessageText MessageType = iota
	EMessageImage
	EMessageFile
	EMessageVideo
	EMessageAudio
)
