package main

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
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
	for {
		select {
		case <-app.cancel:
			/*
				need to do the clean up here, but return I will return
			*/
			return
		default:
			/*
				we need to start reading here
			*/
			msgType, r, err := conn.NextReader()

			/*
				need to find a way to handle this
			*/
			if err != nil {
				/*
					!!!!TODO
					when we have the proper client information, we can output client information

				*/
				app.logger.PrintError(fmt.Errorf("error trying to read from the web socket %w", err), nil)

				// right now if there is an error, I will just return
				// TODO: must perform clean up
				return
			}

			switch msgType {
			case websocket.TextMessage:
				/*
					message structure 4 bytes to hold the size of the header
					get the convert the bytes to int which will basically tell use the size of the header

					the header should be a json
					header:{
						size: //size,
						//other information
					}
					body:{
						// body
					}
				*/
				headerBytes := make([]byte, 4)
				for {

					_, err := readBytes(r, headerBytes)
					if err != nil {
						app.logger.PrintError(err, nil)
						continue
					}

				}

			/*
				this is for audio, video, and other stuff
			*/
			case websocket.BinaryMessage:

			}

		}

	}

}

func readBytes(r io.Reader, buffer []byte) (int, error) {
	n, err := r.Read(buffer)
	if n < len(buffer) {
		return n, fmt.Errorf("read %d bytes, supposed to read %d bytes, %w", n, len(buffer), insufficientBytes)
	}
	return n, err

}

/*
for {
    msgType, r, err := conn.NextReader()
    if err != nil {
        log.Println("nextreader error:", err)
        break
    }

    switch msgType {
    case websocket.TextMessage:
        // Read text in chunks
        buf := make([]byte, 1024) // 1KB chunks
        for {
            n, err := r.Read(buf)
            if n > 0 {
                fmt.Println("Got text chunk:", string(buf[:n]))
            }
            if err == io.EOF {
                break // done with this message
            }
            if err != nil {
                log.Println("read error:", err)
                break
            }
        }

    case websocket.BinaryMessage:
        // Example: save binary data to a file
        f, _ := os.Create("upload.bin")
        defer f.Close()
        if _, err := io.Copy(f, r); err != nil {
            log.Println("copy error:", err)
        }
        fmt.Println("Saved binary message to file")

    default:
        log.Println("Other message type:", msgType)
    }
}




*/
