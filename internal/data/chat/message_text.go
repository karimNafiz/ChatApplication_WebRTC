package chat

/*
	the ID, ChatRoomID, and SenderID can be put in a small struct private to the chat package

*/

type MessageText struct {
	/*
		Sendables that are sent later will have an ID bigger than Sendables sent before
		time(sendable1) > time(sendable2) -> id(sendable1) > id(sendable2)
		we can use this fact along with chat room ID to see if a client is caught up with all the messages
	*/
	ID         int64 // we need to use a number based id to syncing clients when they reconnect
	Body       string
	SenderID   string
	ChatRoomID string
}

func (m *MessageText) GetType() MessageType {
	return EMessageText
}

func NewMessageText(body string, senderID string, chatRoomID string) *MessageText {
	/*

		need to find out whats the best way to create an id for the message
	*/
	return &MessageText{
		Body:       body,
		SenderID:   senderID,
		ChatRoomID: chatRoomID,
	}
}
