package DiscordAPI

import (
	"errors"
	"log"
)

const (
	BUTTON_PRIMARY int = 1 + iota
	BUTTON_SECONDARY
	BUTTON_SUCCESS
	BUTTON_DANGER
	BUTTON_LINK
)

type InteractionCommandData struct {
	Type    int       `json:"type"`
	Options []Options `json:"options"`
	Name    string    `json:"name"`
	ID      string    `json:"id"`
}

type InteractionCommand struct {
	InteractionBase
	Data InteractionCommandData `json:"data"`
}

func (interaction *InteractionCommand) GetOptionBool(key string) (bool, error) {
	for _, option := range interaction.Data.Options {
		if option.Name == key {
			optionString, ok := option.Value.(bool)
			if !ok {
				return false, errors.New("can' get option value")
			}
			return optionString, nil
		}
	}
	return false, errors.New("can' find option value")
}

func (interaction *InteractionCommand) GetOptionString(key string) (string, error) {
	for _, option := range interaction.Data.Options {
		if option.Name == key {
			optionString, ok := option.Value.(string)
			if !ok {
				return "", errors.New("can' get option value")
			}
			return optionString, nil
		}
	}
	return "", errors.New("can' find option value")
}

func (interaction *InteractionCommand) GetOptionInt(key string) (int64, error) {
	for _, option := range interaction.Data.Options {
		if option.Name == key {
			optionString, ok := option.Value.(int64)
			if !ok {
				return -1, errors.New("can' get option value")
			}
			return optionString, nil
		}
	}
	return -1, errors.New("can' find option value")
}

func (interaction *InteractionCommand) AnswerSimple(message MessageCreate) {
	err := interaction.InteractionResponse(InteractionResponsePayload{
		Type: IRPT_MESSAGE,
		Data: message,
	})
	if err != nil {
		log.Println(err)
	}
}
