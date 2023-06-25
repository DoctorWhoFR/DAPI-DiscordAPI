package DiscordEvent

import (
	"azginfr/dapi/DiscordAPI"
)

var HandlerInteractionMessageComponent func(interaction DiscordAPI.InteractionMessageComponent)

func SetHandlerInteractionMessageComponent(f func(interaction DiscordAPI.InteractionMessageComponent)) {
	HandlerInteractionMessageComponent = f
}
