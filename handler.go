package dapi

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"test/dapi/events"
	"test/dapi/internal"
	"test/dapi/restapi"
	"time"

	"github.com/gorilla/websocket"
)

type GatewayMessage struct {
	OpCode int             `json:"op"`
	Data   json.RawMessage `json:"d"`
	Type   string          `json:"t"`
	Seq    int             `json:"s"`
}

// SendHeartbeats Currently Sending heartbeats discord.
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
			internal.LogTrace(err)
			return
		}

		continue
	}
}

// Handle Event
func handleEvents(msg GatewayMessage) {
	internal.LogTrace("Handling message", msg.OpCode)
	switch msg.OpCode {
	case OP_HELLO:
		// Receive the heartbeat interval from the Gateway
		var payload struct {
			HeartbeatInterval int `json:"heartbeat_interval"`
		}
		if err := json.Unmarshal(msg.Data, &payload); err != nil {
			internal.LogTrace(err)
			return
		}

		// Start sending heartbeat messages to the Gateway
		go SendHeartbeats(ws, payload.HeartbeatInterval)
	case OP_HERTHBEAT_ACK:
		internal.LogTrace("ack hearthbeat")
		// Gateway acknowledged our last heartbeat message
	case OP_DISPATCH:
		switch msg.Type {
		case "READY":
			// Handle the READY event
			var event events.ReadyEvent

			err := json.Unmarshal(msg.Data, &event)

			if err != nil {
				fmt.Println(err)
			}

			if events.HandleReady != nil {
				go events.HandleReady(event)
			}
		case "MESSAGE_CREATE":
			// Handle the MESSAGE_CREATE event
			var messageCreate restapi.MessageCreate

			err := json.Unmarshal(msg.Data, &messageCreate)

			if err != nil {
				fmt.Println(err)
			}

			go events.HandleMessageCreate(messageCreate)
		// Add more event handlers here...
		case "INTERACTION_CREATE":
			// // Handle the INTERACTION_CREATE event
			var base restapi.InteractionBase

			err := json.Unmarshal(msg.Data, &base)

			if err != nil {
				fmt.Println(err)
				return
			}

			switch base.Type {
			case 2:
				// APPLICATION_COMMAND
				var interaction restapi.InteractionCommand

				internal.LogTrace(string(msg.Data))

				err := json.Unmarshal(msg.Data, &interaction)

				if err != nil {
					fmt.Println(err)
				}

				if events.HandlerInteractionCreate != nil {
					go events.HandlerInteractionCreate(interaction)
				}
			case 3:
				// MESSAGE_COMPONENT
				var interaction restapi.InteractionButton

				internal.LogTrace(string(msg.Data))

				err := json.Unmarshal(msg.Data, &interaction)

				if err != nil {
					fmt.Println(err)
				}
				go events.HandlerInteractionButton(interaction)
			}
		case "MESSAGE_REACTION_ADD":
			// // Handle the INTERACTION_CREATE event
			var _MessageReactionAdd events.MessageReactionAddEvent

			err := json.Unmarshal(msg.Data, &_MessageReactionAdd)

			if err != nil {
				fmt.Println(err)
			}

			go events.HandleMessageReactionAdd(_MessageReactionAdd)
		case "GUILD_CREATE":
			// // Handle the INTERACTION_CREATE event
			var _GuildCreate events.GuildCreate

			err := json.Unmarshal(msg.Data, &_GuildCreate)

			if err != nil {
				fmt.Println(err)
			}

			go events.HandleGuildCreate(_GuildCreate)
		case "PRESENCE_UPDATE":
			// // Handle the INTERACTION_CREATE event
			var _PresenceUpdate events.PresenceUpdate

			err := json.Unmarshal(msg.Data, &_PresenceUpdate)

			if err != nil {
				fmt.Println(err)
			}

			if events.HandlePresenceUpdate != nil {
				go events.HandlePresenceUpdate(_PresenceUpdate)
			}
		default:
			marshalJSON, err := msg.Data.MarshalJSON()
			if err != nil {
				internal.LogTrace(err)
			}
			//internal.LogTrace(string(marshalJSON))
			internal.LogTrace("no event handler for", msg.Type)

			err = os.WriteFile("./data/"+msg.Type+".json", marshalJSON, 0644)
			if err != nil {
				internal.LogError("cant save file for type", msg.Type)
				return
			}
		}

	}
}

// MainHandler
// MainDiscord Event Handler
// Get & Read 1 message from websocket
// And send it to handleEvents() function
func MainHandler() {
	// Start a goroutine to read messages from the Gateway
	done := make(chan struct{})
	defer close(done)

	for {
		var msg GatewayMessage

		if err := ws.ReadJSON(&msg); err != nil {
			internal.LogTrace(err)
			return
		}

		handleEvents(msg)
	}

}
