package data

type TCPHeader struct {
	MessageType int `json:"message_type"`
	BodySize    int `json:"body_size"`
}

type TCPBody_Text struct {
	Body string `json:"body"`
}
