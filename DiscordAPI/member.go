package DiscordAPI

import (
	"azginfr/dapi/DiscordInternal"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type User struct {
	ID               string      `json:"id"`
	Username         string      `json:"username"`
	Avatar           string      `json:"avatar"`
	Discriminator    string      `json:"discriminator"`
	PublicFlags      int         `json:"public_flags"`
	Flags            int         `json:"flags"`
	Banner           string      `json:"banner"`
	AccentColor      int         `json:"accent_color"`
	GlobalName       interface{} `json:"global_name"`
	AvatarDecoration interface{} `json:"avatar_decoration"`
	DisplayName      interface{} `json:"display_name"`
	BannerColor      string      `json:"banner_color"`
}

type ThreadMember struct {
	Id            string    `json:"id"`
	Flags         int       `json:"flags"`
	JoinTimestamp time.Time `json:"join_timestamp"`
	UserId        string    `json:"user_id"`
}

type Member struct {
	Avatar                     string      `json:"avatar"`
	CommunicationDisabledUntil interface{} `json:"communication_disabled_until"`
	Flags                      int         `json:"flags"`
	JoinedAt                   time.Time   `json:"joined_at"`
	Nick                       string      `json:"nick"`
	Pending                    bool        `json:"pending"`
	PremiumSince               time.Time   `json:"premium_since"`
	Roles                      []string    `json:"roles"`
	User                       User        `json:"user"`
	Mute                       bool        `json:"mute"`
	Deaf                       bool        `json:"deaf"`
}

type GuildMemberUpdate struct {
	Nick       string   `json:"nick"`
	Roles      []string `json:"roles"`
	Mute       bool     `json:"mute"`
	Deaf       bool     `json:"deaf"`
	Channel_id string   `json:"channel_id"`
	Flags      int      `json:"flags"`
}

func (content Member) AddRole(guildID, roleID string) bool {
	answer := make(chan BucketRequestAnswer, 1)
	defer close(answer)

	request := BucketRequest{
		AnswerQueue: answer,
		BucketName:  "guilds",
		Url:         fmt.Sprintf("/guilds/%s/members/%s/roles/%s", guildID, content.User.ID, roleID),
		Methode:     http.MethodPut,
		Payload:     nil,
	}
	addRequest(request)

	// Get response
	response := <-answer

	if response.Res.StatusCode >= 300 {
		DiscordInternal.LogDebug(request, response)
		DiscordInternal.LogDebug("error adding role")
		return false
	}

	return true
}

func (content Member) UpdateMember(guildID string, update GuildMemberUpdate) Member {
	body, err := json.Marshal(update)

	if err != nil {
		return Member{}
	}

	answer := make(chan BucketRequestAnswer, 1)
	defer close(answer)

	request := BucketRequest{
		AnswerQueue: answer,
		BucketName:  "guilds",
		Url:         fmt.Sprintf("/guilds/%s/members/%s", guildID, content.User.ID),
		Methode:     http.MethodPatch,
		Payload:     body,
	}
	addRequest(request)

	// Get response
	response := <-answer

	var member Member

	err = json.Unmarshal(response.Body, &member)

	if err != nil {
		return Member{}
	}

	return member
}
