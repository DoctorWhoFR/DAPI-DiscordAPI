package DiscordAPI

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type InteractionResponseMessage struct {
	Content string      `json:"content,omitempty"`
	Flags   MessageFlag `json:"flags,omitempty"`
}

type InteractionResponsePayloadType int64

const (
	IRPT_UPDATE_MESSAGE_WAITING InteractionResponsePayloadType = 6
	IRPT_UPDATE_MESSAGE         InteractionResponsePayloadType = 7
	IRPT_LOADING                InteractionResponsePayloadType = 5
	IRPT_MESSAGE                InteractionResponsePayloadType = 4
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
	Message           MessageCreate `json:"message"`
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
		log.Println(string(jsonned))
		return errors.New("not an 2xx response")
	}

	return nil
}

func (interaction *InteractionBase) InteractionResponseEditImage(message MessageCreate, file string) error {
	// /interactions/{interaction.id}/{interaction.token}/callback
	jsonned, _ := json.Marshal(message)

	fileContent, err := os.Open(file)
	if err != nil {
		panic(err)
		return err
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	part1, _ := writer.CreateFormFile("files[0]", filepath.Base(fileContent.Name()))
	_, _ = io.Copy(part1, fileContent)
	_ = writer.WriteField("payload_json", string(jsonned))
	err = writer.Close()
	if err != nil {
		fmt.Println(err)
	}

	err = fileContent.Close()
	if err != nil {
		return err
	}

	go func(f *os.File) {
		timer1 := time.NewTimer(5 * time.Second)
		<-timer1.C
		err := os.Remove(f.Name())
		if err != nil {
			panic(err)
			return
		}
		fmt.Println("Timer 2 fired")
	}(fileContent)

	RequestDiscordForm("/webhooks/"+interaction.ApplicationID+"/"+interaction.Token+"/messages/@original", http.MethodPatch, "webhooks", jsonned, false, payload, writer.FormDataContentType())

	return nil
}

func (interaction *InteractionBase) InteractionResponseImage(interactionPayload InteractionResponsePayload, file string) error {
	jsonned, _ := json.Marshal(interactionPayload)

	fileContent, err := os.Open(file)
	if err != nil {
		panic(err)
		return err
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	part1, _ := writer.CreateFormFile("files[0]", filepath.Base(fileContent.Name()))
	_, _ = io.Copy(part1, fileContent)
	_ = writer.WriteField("payload_json", string(jsonned))
	err = writer.Close()
	if err != nil {
		fmt.Println(err)
	}

	err = fileContent.Close()
	if err != nil {
		return err
	}

	go func(f *os.File) {
		timer1 := time.NewTimer(5 * time.Second)
		<-timer1.C
		err := os.Remove(f.Name())
		if err != nil {
			panic(err)
			return
		}
		fmt.Println("Timer 2 fired")
	}(fileContent)

	answer := RequestDiscordForm("/interactions/"+interaction.ID+"/"+interaction.Token+"/callback", http.MethodPost, "interactions", jsonned, false, payload, writer.FormDataContentType())

	if answer.Res.StatusCode > 204 {
		return errors.New("not an sss2xx response")
	}

	return nil
}

func (interaction *InteractionBase) InteractionResponse(interactionPayload InteractionResponsePayload) error {
	start := time.Now()
	jsonned, _ := json.Marshal(interactionPayload)

	answer := RequestDiscord("/interactions/"+interaction.ID+"/"+interaction.Token+"/callback", http.MethodPost, "interactions", jsonned, false)

	log.Println("testing", time.Since(start))
	if answer.Res.StatusCode > 204 {
		return errors.New("not an sss2xx response")
	}

	return nil
}

func (interaction *InteractionBase) AnswerWaitSimple(flags MessageFlag) {
	err := interaction.InteractionResponse(InteractionResponsePayload{
		Type: IRPT_LOADING,
		Data: MessageCreate{Flags: flags},
	})
	if err != nil {
		log.Println(err)
	}
}

func (interaction *InteractionBase) AnswerEdit(message MessageCreate) {
	err := interaction.InteractionResponseEdit(message)
	if err != nil {
		log.Println(err)
	}
}

func (interaction *InteractionBase) AnswerEditImage(message MessageCreate, imagePath string) {
	err := interaction.InteractionResponseEditImage(message, imagePath)
	if err != nil {
		return
	}
}
