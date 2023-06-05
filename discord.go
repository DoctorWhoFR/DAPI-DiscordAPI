package dapi

import (
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"log"
	"os"
	"test/dapi/internal"
	"test/dapi/restapi"
	"time"
)

var ws *websocket.Conn

var BotStarting time.Time

// Main entry function
// Currently start eventHandler & DiscordAPI Handler
func Logging() {
	BotStarting = time.Now()

	err := godotenv.Load()
	if err != nil {
		return
	}

	internal.GetEnvLogLevel()
	internal.LogDebug("starting ")

	// Discord Gateway endpoint URL
	gatewayURL := "wss://gateway.discord.gg/?encoding=json"

	var err2 error

	// Create a new WebSocket connection to the Gateway
	ws, _, err2 = websocket.DefaultDialer.Dial(gatewayURL, nil)
	if err2 != nil {
		log.Fatal(err2)
	}

	/**
	defer func(ws *websocket.Conn) {
		err := ws.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(ws)
	**/
	internal.LogDebug("connected, sending identify ")

	// Send the identify payload to authenticate with the Gateway
	payload := map[string]interface{}{
		"op": OP_IDENTIFY,
		"d": map[string]interface{}{
			"token":   os.Getenv("TOKEN"),
			"intents": 32767,
			"properties": map[string]interface{}{
				"$os":      "linux",
				"$browser": "myapp",
				"$device":  "myapp",
			},
		},
	}
	if err := ws.WriteJSON(payload); err != nil {
		log.Fatal(err)
	}
	internal.LogDebug("sended auth, start commandHandler and events handler, also as main handler")

	go restapi.DiscordAPIHandler()

	go MainHandler()
}
