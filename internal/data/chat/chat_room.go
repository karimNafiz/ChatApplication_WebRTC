package chat

import "time"

/*
	To send messages or attachements, the server will first create a ChatRoom
	example, A and B are clients, if A wants to send a message to B, the server will ->
	create a chat room, so next time A wants to send a message to B, the server will ->
	check if a chat room exists or not, if not it will create it
	the chat room will persist in database

*/

type ChatRoom struct {
	ID        string
	MemberIDs []string
	CreatedAt time.Time
}
