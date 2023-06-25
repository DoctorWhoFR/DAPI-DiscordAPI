package DiscordClient

import (
	"azginfr/dapi/DiscordAPI"
	"azginfr/dapi/DiscordEvent"
	"azginfr/dapi/DiscordInternal"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

/*
ConnectBot

Entry point function used to lunch discord base code.

Example :

	func main() {
		// login to discord & start everything
		go DiscordClient.ConnectBot()

		DiscordEvent.AddDiscordCommand(commands.RegisterCommandDiscord())
	}

// TODO add parameter token to ConnectBot function
// make the user able to provide to token in self in reverse of getting this one in env variable
// like
// 	`ConnectBot(token string)`
// <!-- #DAPI -->
*/
func ConnectBot() {
	err := godotenv.Load()
	if err != nil {
		return
	}

	DiscordInternal.GetEnvLogLevel()
	DiscordInternal.LogDebug("starting ")

	// Discord Gateway endpoint URL
	gatewayURL := "wss://gateway.discord.gg/?encoding=json"

	var err2 error

	// Create a new WebSocket connection to the Gateway
	DiscordEvent.WebSocket, _, err2 = websocket.DefaultDialer.Dial(gatewayURL, nil)
	if err2 != nil {
		log.Fatal(err2)
	}

	DiscordInternal.LogDebug("connected, sending identify ")

	// Send the identify payload to authenticate with the Gateway
	payload := map[string]interface{}{
		"op": DiscordEvent.OP_IDENTIFY,
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
	if err := DiscordEvent.WebSocket.WriteJSON(payload); err != nil {
		log.Fatal(err)
	}
	DiscordInternal.LogDebug("sent auth, start commandHandler and DiscordEvents handler, also as main handler")

	go DiscordAPI.NewBucketHandler()

	time.Sleep(time.Second)

	go DiscordEvent.MainEventHandler()

}
