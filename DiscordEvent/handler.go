package DiscordEvent

import (
	"azginfr/dapi/DiscordAPI"
	"azginfr/dapi/DiscordInternal"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"math/rand"
	"os"
	"time"
)

var DiscordWebSocket *websocket.Conn

const (
	OP_DISPATCH      = 0 // DISPATCH TYPE
	OP_HEARTHBEAT    = 1
	OP_IDENTIFY      = 2
	OP_HELLO         = 10
	OP_HERTHBEAT_ACK = 11
)

/*
DiscordGatewayMessage

Representation of a discord gateway message
*/
type DiscordGatewayMessage struct {
	OpCode int             `json:"op"`
	Data   json.RawMessage `json:"d"`
	Type   string          `json:"t"`
	Seq    int             `json:"s"`
}

// SendHeartbeats
// Currently Sending heartbeats to discord.
//   - Seed the random number generator with the current time.
//   - Generate a random number between 0 and 1.
//   - Wait for heartbeat interval : `time.Duration(float64(interval)*r) * time.Millisecond`
//   - Send a heartbeat message to the Gateway.
func SendHeartbeats(ws *websocket.Conn, interval int) {
	randG := rand.New(rand.NewSource(time.Now().UnixNano()))

	ackCh := make(chan struct{})
	defer close(ackCh)

	for {
		r := randG.Float64()

		duration := time.Duration(float64(interval)*r) * time.Millisecond

		<-time.After(duration)

		payload := map[string]interface{}{
			"op": OP_HEARTHBEAT,
			"d":  nil,
		}
		if err := ws.WriteJSON(payload); err != nil {
			DiscordInternal.LogTrace(err)
			return
		}

		continue
	}
}

// HandleGatewayMessage
func HandleGatewayMessage(msg DiscordGatewayMessage) {
	DiscordInternal.LogTrace("Handling message", msg.OpCode)
	switch msg.OpCode {
	case OP_HELLO:
		// Receive the heartbeat interval from the Gateway
		var payload struct {
			HeartbeatInterval int `json:"heartbeat_interval"`
		}
		if err := json.Unmarshal(msg.Data, &payload); err != nil {
			DiscordInternal.LogTrace(err)
			return
		}

		// Start sending heartbeat messages to the Gateway
		go SendHeartbeats(DiscordWebSocket, payload.HeartbeatInterval)
	case OP_HERTHBEAT_ACK:
		DiscordInternal.LogTrace("ack hearthbeat")
		// Gateway acknowledged our last heartbeat message
	case OP_DISPATCH:
		switch msg.Type {
		case "INTERACTION_CREATE":
			// // Handle the INTERACTION_CREATE event
			var base DiscordAPI.InteractionBase

			err := json.Unmarshal(msg.Data, &base)

			if err != nil {
				fmt.Println(err)
				return
			}
			switch base.Type {
			case 2:
				// APPLICATION_COMMAND
				var interaction DiscordAPI.InteractionCommand

				DiscordInternal.LogTrace(string(msg.Data))

				err := json.Unmarshal(msg.Data, &interaction)

				if err != nil {
					fmt.Println(err)
				}

				if HandlerInteractionCommandEvent != nil {
					go HandlerInteractionCommandEvent(interaction)
				}
			case 3:
				// MESSAGE_COMPONENT
				var interaction DiscordAPI.InteractionMessageComponent

				err = os.WriteFile("./data/interaction_component.json", msg.Data, 0644)

				err := json.Unmarshal(msg.Data, &interaction)

				if err != nil {
					fmt.Println(err)
				}
				go HandlerInteractionMessageComponent(interaction)
			case 5:
				var interaction DiscordAPI.InteractionModalSubmit

				err = os.WriteFile("./data/interaction_component.json", msg.Data, 0644)

				err := json.Unmarshal(msg.Data, &interaction)

				if err != nil {
					fmt.Println(err)
				}
				go HandlerInteractionModalSubmit(interaction)
			}
		default:
			saveEventActivated := os.Getenv("SAVE_EVENT")
			DiscordInternal.LogTrace("no event handler for", msg.Type)

			if saveEventActivated == "1" {
				marshalJSON, err := msg.Data.MarshalJSON()
				if err != nil {
					DiscordInternal.LogTrace(err)
				}

				err = os.WriteFile("./data/"+msg.Type+".json", marshalJSON, 0644)
				if err != nil {
					DiscordInternal.LogError("cant SaveCommand file for type", msg.Type)
					return
				}
			}
		}
	}
}

// MainEventHandler
// MainDiscord Event Handler
// Get & Read 1 message from websocket
// And send it to HandleGatewayMessage() function
func MainEventHandler() {
	done := make(chan struct{})
	defer close(done)

	for {
		var msg DiscordGatewayMessage

		if err := DiscordWebSocket.ReadJSON(&msg); err != nil {
			DiscordInternal.LogTrace(err)
			return
		}

		HandleGatewayMessage(msg)
	}

}
