package DiscordAPI

import "log"

type InteractionModalSubmit struct {
	InteractionBase
	Message  MessageCreate    `json:"message,omitempty"`
	Data     MessageComponent `json:"data,omitempty"`
	CustomID string           `json:"customID,omitempty"`
}

func (interaction *InteractionModalSubmit) AnswerSimple(message MessageCreate) {
	err := interaction.InteractionResponse(InteractionResponsePayload{
		Type: 7,
		Data: message,
	})
	if err != nil {
		log.Println(err)
	}
}

func (interaction *InteractionModalSubmit) AnswerWait(flags MessageFlag) {
	err := interaction.InteractionResponse(InteractionResponsePayload{
		Type: IRPT_UPDATE_MESSAGE_WAITING,
		Data: MessageCreate{Flags: flags},
	})
	if err != nil {
		log.Println(err)
	}
}
