package DiscordEvent

import (
	"azginfr/dapi/DiscordAPI"
)

var HandlerInteractionModalSubmit func(interaction DiscordAPI.InteractionModalSubmit)

func SetHandlerInteractionModalSubmit(f func(interaction DiscordAPI.InteractionModalSubmit)) {
	HandlerInteractionModalSubmit = f
}
