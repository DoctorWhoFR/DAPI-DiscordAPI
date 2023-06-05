package restapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

type InteractionResponseMessage struct {
	Content string       `json:"content,omitempty"`
	Flags   MessageFlags `json:"flags,omitempty"`
}

type InteractionResponsePayloadType int64

const (
	IRPT_LOADING InteractionResponsePayloadType = 5
	IRPT_MESSAGE InteractionResponsePayloadType = 4
)

type InteractionResponsePayload struct {
	Type InteractionResponsePayloadType `json:"type,omitempty"`
	Data MessageCreate                  `json:"data,omitempty"`
}

type Options struct {
	Value any    `json:"value,omitempty"`
	Type  int    `json:"type,omitempty"`
	Name  string `json:"name,omitempty"`
}

type InteractionBase struct {
	Version           int           `json:"version,omitempty"`
	Type              int           `json:"type,omitempty"`
	Token             string        `json:"token,omitempty"`
	Member            Member        `json:"member,omitempty"`
	Locale            string        `json:"locale,omitempty"`
	ID                string        `json:"id,omitempty"`
	GuildLocale       string        `json:"guild_locale,omitempty"`
	GuildId           string        `json:"guild_id,omitempty"`
	Entitlements      []interface{} `json:"entitlements,omitempty"`
	EntitlementSkuIds []interface{} `json:"entitlement_sku_ids,omitempty"`
	ChannelId         string        `json:"channel_id,omitempty"`
	Channel           Channel       `json:"channel,omitempty"`
	ApplicationID     string        `json:"application_id,omitempty"`
	AppPermissions    string        `json:"app_permissions,omitempty"`
}

const (
	BUTTON_PRIMARY int = 1 + iota
	BUTTON_SECONDARY
	BUTTON_SUCCESS
	BUTTON_DANGER
	BUTTON_LINK
)

type InteractionButton struct {
	InteractionBase
	Data struct {
		CustomId      string `json:"custom_id"`
		ComponentType int    `json:"component_type"`
	} `json:"data"`
}

type InteractionCommand struct {
	InteractionBase
	Data struct {
		Type    int       `json:"type"`
		Options []Options `json:"options"`
		Name    string    `json:"name"`
		ID      string    `json:"id"`
	} `json:"data"`
}

func (interaction *InteractionBase) InteractionResponseEdit(message MessageCreate) error {
	// /interactions/{interaction.id}/{interaction.token}/callback
	jsonned, _ := json.Marshal(message)

	answer := make(chan BucketRequestAnswer, 1)
	request := BucketRequest{AnswerQueue: answer,
		BucketName: "webhooks",
		Url:        "/webhooks/" + interaction.ApplicationID + "/" + interaction.Token + "/messages/@original",
		Methode:    http.MethodPatch,
		Payload:    jsonned,
	}
	addRequest(request)

	response := <-answer

	if response.Res.StatusCode >= 300 {
		log.Println(string(response.Body))
		return errors.New("not an 2xx response")
	}

	return nil
}

func (interaction *InteractionBase) InteractionResponse(interactionPayload InteractionResponsePayload) error {
	start := time.Now()
	jsonned, _ := json.Marshal(interactionPayload)

	answer := RequestDiscord("/interactions/"+interaction.ID+"/"+interaction.Token+"/callback", http.MethodPost, "interactions", jsonned, false)

	log.Println("testing", time.Since(start))
	log.Println(answer.Res.StatusCode)
	if answer.Res.StatusCode > 204 {
		return errors.New("not an sss2xx response")
	}

	return nil
}

func (interaction *InteractionBase) AnswerError(message string) {
	// TODO: search a way to answer interaction if error
	log.Println("error", message)
	/**
	if err := interaction.InteractionResponse(InteractionResponsePayload{}); err != nil {
		internal.LogError("error during response", err)
		return
	}
	**/
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
