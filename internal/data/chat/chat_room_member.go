package chat

import "time"

/*
	a chat room member is a wrapper around the client
	the chat room member struct holds information regarding chat rooms, messages ->
	which seem inappropriate to put into the client struct
	so I'm creating a wrapper around the client to hold these information

*/

type ChatRoomMember struct {
	ID         string // same as the client ID
	ChatRoomID string

	/*
		this variable will tell me what was the last message received by the client in a chat room
	*/
	LastRecievedID int64
	CreatedAt      time.Time
}
