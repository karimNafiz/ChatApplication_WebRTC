package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/websocket"

	// TODO: change this to your actual module path for the data package
	"github.com/karimNafiz/ChatApplication_WebRTC/internal/data"
)

var upgrader = websocket.Upgrader{

	/*
		this is mainly for development
		need to decide if I want this to be within the app struct

	*/
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	insufficientBytes error = fmt.Errorf("insufficient bytes")
)

func (app *application) clientRegisterHandler(w http.ResponseWriter, r *http.Request) {

}

func (app *application) clientLoginHandler(w http.ResponseWriter, r *http.Request) {

}

/*

	clients need to be authenticated by jwt tokens

*/

func (app *application) clientEstablishWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		app.logError(r, fmt.Errorf("could not establish web socket connection %w ", err))
		return
	}
	/*
		need to add more information about the client in this log
	*/
	app.logger.PrintInfo("established webSocket connection with client", nil)

	go app.manageClientSocket(conn)
}

/*
	this function will handle the web socket connection between the client and the server

*/

func (app *application) manageClientSocket(conn *websocket.Conn) {
	// Optional: set read limits to avoid abuse (tune for your app)
	const (
		maxHeaderBytes = 64 * 1024       // 64KB header cap
		maxBodyBytes   = 4 * 1024 * 1024 // 4MB body cap (tune)
	)

	/*
		important stuff
	*/
	defer func() {
		_ = conn.Close()
	}()

	for {
		select {
		case <-app.cancel:
			// TODO: any conn-specific cleanup here
			return
		default:
			// keeping this empty for beauty
		}

		msgType, r, err := conn.NextReader()
		if err != nil {
			app.logger.PrintError(fmt.Errorf("websocket next reader: %w", err), nil)
			return
		}

		switch msgType {
		case websocket.TextMessage:

			// read header size
			headerLen, err := readUint32BE(r)
			if err != nil {
				app.logger.PrintError(fmt.Errorf("reading header length: %w", err), nil)
				/*
					need to handle this better
				*/
				return
			}
			if headerLen == 0 || headerLen > maxHeaderBytes {
				app.logger.PrintError(fmt.Errorf("invalid header length: %d", headerLen), nil)
				/*
					need to handle this better, instead of returning I need to clean out the buffer so that I can continue getting more data
				*/
				return
			}

			// parse header json
			headerBuf := make([]byte, int(headerLen))
			if err := readFull(r, headerBuf); err != nil {
				app.logger.PrintError(fmt.Errorf("reading header: %w", err), nil)
				return
			}

			var hdr data.TCPHeader
			if err := json.Unmarshal(headerBuf, &hdr); err != nil {
				app.logger.PrintError(fmt.Errorf("unmarshal header JSON: %w", err), nil)
				// TODO: need to handle this better
				return
			}

			// house keeping
			if hdr.BodySize < 0 || hdr.BodySize > maxBodyBytes {
				app.logger.PrintError(fmt.Errorf("invalid body size: %d", hdr.BodySize), nil)
				// TODO: need to handle this better
				return
			}

			// parse body, for text it should be json
			bodyBuf := make([]byte, hdr.BodySize)
			if err := readFull(r, bodyBuf); err != nil {
				app.logger.PrintError(fmt.Errorf("reading body: %w", err), nil)
				return
			}

			switch hdr.MessageType {

			default:
				// Assume text payload encoded as JSON matching TCPBody_Text
				var body data.TCPBody_Text
				if err := json.Unmarshal(bodyBuf, &body); err != nil {
					app.logger.PrintError(fmt.Errorf("unmarshal body JSON: %w", err), nil)
					return
				}

				// right now just logging it out
				app.logger.PrintInfo("received text body from client", map[string]string{
					"len":  fmt.Sprintf("%d", len(body.Body)),
					"text": body.Body, // careful with logging PII; trim in prod
				})

				// simple ack
				w, err := conn.NextWriter(websocket.TextMessage)
				if err == nil {
					ack := map[string]any{"ok": true, "type": hdr.MessageType}
					if enc, e := json.Marshal(ack); e == nil {
						// Frame format for reply (simple JSON, no length prefix)
						_, _ = w.Write(enc)
					}
					_ = w.Close()
				}
			}

			// 6) Drain to end-of-frame if there’s leftover (shouldn’t be, but safe)
			_, _ = io.Copy(io.Discard, r)

		case websocket.BinaryMessage:
			// You said to ignore binary for now; just drain it.
			_, _ = io.Copy(io.Discard, r)

		default:
			// Control or unsupported types; drain and continue
			_, _ = io.Copy(io.Discard, r)
		}
	}
}

// readFull reads exactly len(buf) bytes into buf (or returns an error).
/*
	the io.ReadFull function gave me so much pain in the FileUploadMicroservice

*/
func readFull(r io.Reader, buf []byte) error {
	total := 0
	for total < len(buf) {
		n, err := r.Read(buf[total:])
		if n > 0 {
			total += n
		}
		if err != nil {
			if errors.Is(err, io.EOF) && total == len(buf) {
				return nil
			}
			return fmt.Errorf("readFull: read %d/%d: %w", total, len(buf), err)
		}
	}
	return nil
}

// readUint32BE reads 4 bytes and returns a big-endian uint32.
/*
	in the documentation need to mention, we will be using big-endian
*/
func readUint32BE(r io.Reader) (uint32, error) {
	var b [4]byte
	if err := readFull(r, b[:]); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(b[:]), nil
}
