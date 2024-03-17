// neutralino-extension
//
// Neutralino GoExtension
//
// (c)2024 Harald Schneider - marketmix.com

package neutralino_extension

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"io"
	"net/url"
	"os"
	"os/signal"
)

const Version = "1.0.6"

type Config struct {
	NlPort         string `json:"nlPort"`
	NlToken        string `json:"nlToken"`
	NlExtensionId  string `json:"nlExtensionId"`
	NlConnectToken string `json:"nlConnectToken"`
}

type EventMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type DataPacket struct {
	Id          string       `json:"id"`
	Method      string       `json:"method"`
	AccessToken string       `json:"accessToken"`
	Data        EventMessage `json:"data"`
}

type WSClient struct {
	url    url.URL
	socket *websocket.Conn
	debug  bool
}

var ExtConfig = Config{}

// Send : Send a websocket message
func (wsclient *WSClient) Send(event string, data map[string]interface{}) {

	// Prep data packet
	//
	var msg = DataPacket{}
	msg.Id = uuid.New().String()
	msg.Method = "app.broadcast"
	msg.AccessToken = ExtConfig.NlToken
	msg.Data.Event = event
	msg.Data.Data = data

	var d, err = json.Marshal(msg)
	if err != nil {
		fmt.Println("Error in marshaling data-packet.")
		return
	}

	if wsclient.debug {
		fmt.Printf("%sSent: %s%s\n", "\u001B[32m", string(d), "\u001B[0m")
	}

	// Send
	//
	err = wsclient.socket.WriteMessage(websocket.TextMessage, []byte(d))
	if err != nil {
		if wsclient.debug {
			fmt.Println("Error in Send(): ", err)
		}
		return
	}
}

// SendMessageString : Send a WebSocket message with a string in a result field
func (wsclient *WSClient) SendMessageString(event string, data string) {
	msg := make(map[string]interface{})
	msg["result"] = data
	wsclient.Send(event, msg)
}

// Run : Start Server main loop
func (wsclient *WSClient) Run(callback func(message EventMessage), debug bool) {

	wsclient.debug = debug

	decoder := json.NewDecoder(os.Stdin)

	err := decoder.Decode(&ExtConfig)
	if err != nil {
		if err != io.EOF {
			fmt.Println(err)
		}
	}

	// Listen to keyboard interrupts when in debug mode
	//
	sigInt := make(chan os.Signal, 1)
	if debug {
		signal.Notify(sigInt, os.Interrupt)
	}

	// Connect to server
	//
	var addr = "127.0.0.1:" + ExtConfig.NlPort
	var path = "?extensionId=" + ExtConfig.NlExtensionId + "&connectToken=" + ExtConfig.NlConnectToken

	wsclient.url = url.URL{Scheme: "ws", Host: addr, Path: path}
	if wsclient.debug {
		fmt.Printf("Connecting to %s\n", wsclient.url.String())
	}

	wsclient.socket, _, err = websocket.DefaultDialer.Dial(wsclient.url.String(), nil)
	if err != nil {
		if wsclient.debug {
			fmt.Println("Connect: ", err)
		}
	}
	defer wsclient.socket.Close()

	// WebSocket read loop
	//
	go func() {
		for {
			_, msg, err := wsclient.socket.ReadMessage()
			if err != nil {
				if wsclient.debug {
					fmt.Println("ERROR in read loop: ", err)
					continue
				}
			}
			if wsclient.debug {
				fmt.Printf("%sReceived: %s%s\n", "\u001B[91m", msg, "\u001B[0m")
			}

			// Parse JSON and forward to callback:
			//
			var d EventMessage
			err = json.Unmarshal([]byte(msg), &d)
			if err != nil {
				if wsclient.debug {
					fmt.Println("ERROR in read loop, while unmarshalling JSON: ", err)
					continue
				}
			}

			if wsclient.IsEvent(d, "windowClose") || wsclient.IsEvent(d, "appClose") {
				wsclient.quit()
				continue
			}

			callback(d)
		}
	}()

	// Signal-listener loop
	//
	for {
		if <-sigInt != nil {
			fmt.Println("Interrupted by keyboard interaction ...")
			wsclient.quit()
		}
	}
}

// IsEvent : Return true, if data matches event-name
func (wsclient *WSClient) IsEvent(data EventMessage, event string) bool {
	if data.Event == event {
		return true
	}
	return false
}

// quit: Do Harakiri
func (wsclient *WSClient) quit() {
	var pid = os.Getpid()
	fmt.Println("Killing own process with PID ", pid)
	process, _ := os.FindProcess(pid)
	err := process.Signal(os.Kill)
	if err != nil {
		return
	}
}
