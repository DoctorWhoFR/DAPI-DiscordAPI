package DiscordAPI

import "log"

type InteractionMessageComponent struct {
	InteractionBase
	Data struct {
		Values        []string `json:"values"`
		CustomID      string   `json:"custom_id"`
		ComponentType int      `json:"component_type"`
	} `json:"data"`
	Message MessageCreate `json:"message,omitempty"`
}

func (interaction *InteractionMessageComponent) AnswerSimple(message MessageCreate) {
	err := interaction.InteractionResponse(InteractionResponsePayload{
		Type: 7,
		Data: message,
	})
	if err != nil {
		log.Println(err)
	}
}

func (interaction *InteractionMessageComponent) AnswerWait(flags MessageFlag) {
	err := interaction.InteractionResponse(InteractionResponsePayload{
		Type: IRPT_UPDATE_MESSAGE_WAITING,
		Data: MessageCreate{Flags: flags},
	})
	if err != nil {
		log.Println(err)
	}
}

func (interaction *InteractionMessageComponent) AnswerModal(message MessageCreate) {
	err := interaction.InteractionResponse(InteractionResponsePayload{
		Type: 9,
		Data: message,
	})
	if err != nil {
		log.Println(err)
	}
}
