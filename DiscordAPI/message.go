package DiscordAPI

import (
	"azginfr/dapi/DiscordInternal"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

type EmbedAuthor struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	IconURL string `json:"icon_url"`
}

type Embed struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	Color       int       `json:"color"`
	Timestamp   time.Time `json:"timestamp"`
	Footer      struct {
		IconURL string `json:"icon_url"`
		Text    string `json:"text"`
	} `json:"footer"`
	Thumbnail struct {
		URL string `json:"url"`
	} `json:"thumbnail"`
	Image struct {
		URL string `json:"url"`
	} `json:"image"`
	Author EmbedAuthor  `json:"author"`
	Fields []EmbedField `json:"fields"`
}

type MessageFlag int64

const (
	MF_EPHEMERAL MessageFlag = 64
	MF_NO        MessageFlag = 0
)

type Reaction struct {
	Emoji struct {
		ID   any    `json:"id"`
		Name string `json:"name"`
	} `json:"emoji"`
	Count        int `json:"count"`
	CountDetails struct {
		Burst  int `json:"burst"`
		Normal int `json:"normal"`
	} `json:"count_details"`
	BurstColors []any `json:"burst_colors"`
	MeBurst     bool  `json:"me_burst"`
	Me          bool  `json:"me"`
	BurstCount  int   `json:"burst_count"`
}

type MessageSubComponent struct {
}

const (
	MessageComponentActionRow = 1 + iota
	MessageComponentButton
	MessageComponentStringSelect
)

type MessageComponent struct {
	Type        int                `json:"type,omitempty"`
	Label       string             `json:"label,omitempty"`
	Style       int                `json:"style,omitempty"`
	CustomId    string             `json:"custom_id,omitempty"`
	Components  []MessageComponent `json:"components,omitempty"`
	Url         string             `json:"url,omitempty"`
	Disabled    bool               `json:"disabled,omitempty"`
	Emoji       string             `json:"emoji,omitempty"`
	Options     []MessageComponent `json:"options,omitempty"`
	Value       string             `json:"value,omitempty"`
	Description string             `json:"description,omitempty"`
	Default     bool               `json:"default,omitempty"`
	Required    bool               `json:"required,omitempty"`
}

type MessageCreate struct {
	Type              int                `json:"type,omitempty"`
	Tts               bool               `json:"tts,omitempty"`
	Timestamp         time.Time          `json:"timestamp,omitempty"`
	ReferencedMessage any                `json:"referenced_message,omitempty"`
	Pinned            bool               `json:"pinned,omitempty"`
	Nonce             string             `json:"nonce,omitempty"`
	Mentions          []any              `json:"mentions,omitempty"`
	MentionRoles      []any              `json:"mention_roles,omitempty"`
	MentionEveryone   bool               `json:"mention_everyone,omitempty"`
	Member            Member             `json:"member,omitempty"`
	ID                string             `json:"id,omitempty"`
	Flags             MessageFlag        `json:"flags,omitempty"`
	Embeds            []Embed            `json:"embeds"`
	EditedTimestamp   any                `json:"edited_timestamp,omitempty"`
	Content           string             `json:"content"`
	Components        []MessageComponent `json:"components"`
	ChannelID         string             `json:"channel_id,omitempty"`
	Author            User               `json:"author,omitempty"`
	Attachments       []interface{}      `json:"attachments"`
	GuildID           string             `json:"guild_id,omitempty"`
	Channel           Channel            `json:"channel,omitempty"`
	Reactions         []Reaction         `json:"reactions,omitempty"`
	CustomID          string             `json:"custom_id,omitempty"`
	Title             string             `json:"title,omitempty"`
}

// DeleteMessage WARNING should use Channel.BulkDeleteMessages if you want to remove more than 1 message
func (message *MessageCreate) DeleteMessage() error {
	body, _ := json.Marshal(message)

	answer := make(chan BucketRequestAnswer, 1)

	request := BucketRequest{
		AnswerQueue: answer,
		BucketName:  "channels",
		Url:         "/channels/" + message.ChannelID + "/messages/" + message.ID,
		Methode:     http.MethodDelete,
		Payload:     body,
	}

	addRequest(request)

	response := <-answer
	close(answer)

	err := json.Unmarshal(response.Body, &message)
	if err != nil {
		return err
	}
	return nil
}

func (message *MessageCreate) UpdateMessage() error {
	body, _ := json.Marshal(message)

	answer := make(chan BucketRequestAnswer, 1)

	request := BucketRequest{
		AnswerQueue: answer,
		BucketName:  "channels",
		Url:         "/channels/" + message.ChannelID + "/messages/" + message.ID,
		Methode:     http.MethodPatch,
		Payload:     body,
	}

	addRequest(request)

	response := <-answer
	close(answer)

	err := json.Unmarshal(response.Body, &message)
	if err != nil {
		return err
	}
	return nil
}

func (message *MessageCreate) CrossPostMessage() error {
	answer := make(chan BucketRequestAnswer, 1)
	request := BucketRequest{
		AnswerQueue: answer,
		BucketName:  "channels",
		Url:         "/channels/" + message.ChannelID + "/messages/" + message.ID + "/crosspost",
		Methode:     http.MethodPost,
		Payload:     nil,
	}

	addRequest(request)

	response := <-answer
	close(answer)

	err := json.Unmarshal(response.Body, &message)
	if err != nil {
		return err
	}

	return nil
}

func (message *MessageCreate) CreateSelfReaction(emoji string) error {
	answer := make(chan BucketRequestAnswer, 1)
	request := BucketRequest{
		AnswerQueue: answer,
		BucketName:  "channels",
		Url:         "/channels/" + message.ChannelID + "/messages/" + message.ID + "/reactions/" + emoji + "/@me",
		Methode:     http.MethodPut,
		Payload:     nil,
	}

	addRequest(request)

	<-answer
	close(answer)

	getMessage, err := message.Channel.GetMessage(message.ID)
	if err != nil {
		return err
	}

	*message = getMessage

	return nil
}

func (message *MessageCreate) DeleteSelfReaction(emoji string) error {
	answer := make(chan BucketRequestAnswer, 1)
	request := BucketRequest{
		AnswerQueue: answer,
		BucketName:  "channels",
		Url:         "/channels/" + message.ChannelID + "/messages/" + message.ID + "/reactions/" + emoji + "/@me",
		Methode:     http.MethodDelete,
		Payload:     nil,
	}

	addRequest(request)

	<-answer
	close(answer)

	getMessage, err := message.Channel.GetMessage(message.ID)
	if err != nil {
		return err
	}

	*message = getMessage

	return nil
}

func (message *MessageCreate) DeleteUserReaction(emoji, userid string) error {
	answer := make(chan BucketRequestAnswer, 1)
	request := BucketRequest{
		AnswerQueue: answer,
		BucketName:  "channels",
		Url:         "/channels/" + message.ChannelID + "/messages/" + message.ID + "/reactions/" + emoji + "/" + userid,
		Methode:     http.MethodPost,
		Payload:     nil,
	}

	addRequest(request)

	<-answer
	close(answer)

	return nil
}

func (message *MessageCreate) GetReactionUsers(emoji string) ([]User, error) {
	answer := make(chan BucketRequestAnswer, 1)
	request := BucketRequest{
		AnswerQueue: answer,
		BucketName:  "channels",
		Url:         "/channels/" + message.ChannelID + "/messages/" + message.ID + "/reactions/" + emoji,
		Methode:     http.MethodGet,
		Payload:     nil,
	}

	addRequest(request)

	response := <-answer
	close(answer)

	var usersLists []User

	err := json.Unmarshal(response.Body, &usersLists)
	if err != nil {
		return usersLists, err
	}

	DiscordInternal.LogDebug(response.Body, usersLists)

	return usersLists, nil
}

func (message *MessageCreate) DeleteAllReaction() error {
	answer := make(chan BucketRequestAnswer, 1)
	request := BucketRequest{
		AnswerQueue: answer,
		BucketName:  "channels",
		Url:         "/channels/" + message.ChannelID + "/messages/" + message.ID + "/reactions",
		Methode:     http.MethodDelete,
		Payload:     nil,
	}

	addRequest(request)

	<-answer
	close(answer)

	return nil
}

func (message *MessageCreate) DeleteAllEmojiReaction(emoji string) error {
	answer := make(chan BucketRequestAnswer, 1)
	request := BucketRequest{
		AnswerQueue: answer,
		BucketName:  "channels",
		Url:         "/channels/" + message.ChannelID + "/messages/" + message.ID + "/reactions/" + emoji,
		Methode:     http.MethodDelete,
		Payload:     nil,
	}

	addRequest(request)

	<-answer
	close(answer)

	return nil
}

func (message *MessageCreate) PinMessage() error {
	_url := fmt.Sprintf("/channels/%s/pins/%s", message.ChannelID, message.ID)

	answer := RequestDiscord(_url, http.MethodPut, "channels", nil, false)

	if answer.Err != nil {
		return errors.New("can't  based on technical error" + answer.Err.Error())
	}

	if answer.Res.StatusCode > http.StatusResetContent {
		return errors.New("can't  based on discord error: " + string(answer.Body))
	}

	return nil
}

func (message *MessageCreate) UnPinMessage() error {
	_url := fmt.Sprintf("/channels/%s/pins/%s", message.ChannelID, message.ID)

	answer := RequestDiscord(_url, http.MethodDelete, "channels", nil, false)

	if answer.Err != nil {
		return errors.New("can't  based on technical error" + answer.Err.Error())
	}

	if answer.Res.StatusCode > http.StatusResetContent {
		return errors.New("can't  based on discord error: " + string(answer.Body))
	}

	return nil
}

/*
StartThread

# Start Thread from Message

POST/channels/{channel.id}/messages/{message.id}/threads

Creates a new thread from an existing message.

Returns a channel on success, and a 400 BAD REQUEST on invalid parameters.

Fires a Thread Create and a Message UpdateCommand Gateway event.

When called on a GUILD_TEXT channel, creates a PUBLIC_THREAD.

When called on a GUILD_ANNOUNCEMENT channel, creates a ANNOUNCEMENT_THREAD.

Does not work on a GUILD_FORUM channel.

# The id of the created thread will be the same as the id of the source message

, and as such a message can only have a single thread created from it.

This endpoint supports the X-Audit-Log-Reason header.
*/
func (message *MessageCreate) StartThread(thread StartThreadPayload) (Channel, error) {
	marshal, err := json.Marshal(thread)
	if err != nil {
		return Channel{}, err
	}

	answer := RequestDiscord(fmt.Sprintf("/channels/%s/messages/%s/threads", message.ChannelID, message.ID), http.MethodPost, "channels", marshal, true)

	if answer.Err != nil {
		return Channel{}, errors.New("can't  based on technical error" + answer.Err.Error())
	}

	if answer.Res.StatusCode > http.StatusResetContent {
		return Channel{}, errors.New("can't  based on discord error: " + string(answer.Body))
	}

	var channel Channel

	err = json.Unmarshal(answer.Body, &channel)
	if err != nil {
		return Channel{}, err
	}

	return channel, nil
}
