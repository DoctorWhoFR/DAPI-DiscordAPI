package restapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Role struct {
	ID             string      `json:"id"`
	Name           string      `json:"name"`
	Description    interface{} `json:"description"`
	Permissions    int         `json:"permissions"`
	Position       int         `json:"position"`
	Color          int         `json:"color"`
	Hoist          bool        `json:"hoist"`
	Managed        bool        `json:"managed"`
	Mentionable    bool        `json:"mentionable"`
	Icon           interface{} `json:"icon"`
	UnicodeEmoji   interface{} `json:"unicode_emoji"`
	Flags          int         `json:"flags"`
	PermissionsNew string      `json:"permissions_new"`
	Tags           struct {
		BotID string `json:"bot_id"`
	} `json:"tags,omitempty"`
}

type Guild struct {
	ID                          string        `json:"id"`
	Name                        string        `json:"name"`
	Icon                        string        `json:"icon"`
	Description                 interface{}   `json:"description"`
	HomeHeader                  interface{}   `json:"home_header"`
	Splash                      string        `json:"splash"`
	DiscoverySplash             interface{}   `json:"discovery_splash"`
	Features                    []string      `json:"features"`
	Emojis                      []interface{} `json:"emojis"`
	Stickers                    []interface{} `json:"stickers"`
	Banner                      interface{}   `json:"banner"`
	OwnerID                     string        `json:"owner_id"`
	ApplicationID               interface{}   `json:"application_id"`
	Region                      string        `json:"region"`
	AfkChannelID                string        `json:"afk_channel_id"`
	AfkTimeout                  int           `json:"afk_timeout"`
	SystemChannelID             string        `json:"system_channel_id"`
	WidgetEnabled               bool          `json:"widget_enabled"`
	WidgetChannelID             interface{}   `json:"widget_channel_id"`
	VerificationLevel           int           `json:"verification_level"`
	Roles                       []Role        `json:"roles"`
	DefaultMessageNotifications int           `json:"default_message_notifications"`
	MfaLevel                    int           `json:"mfa_level"`
	ExplicitContentFilter       int           `json:"explicit_content_filter"`
	MaxPresences                interface{}   `json:"max_presences"`
	MaxMembers                  int           `json:"max_members"`
	MaxStageVideoChannelUsers   int           `json:"max_stage_video_channel_users"`
	MaxVideoChannelUsers        int           `json:"max_video_channel_users"`
	VanityURLCode               interface{}   `json:"vanity_url_code"`
	PremiumTier                 int           `json:"premium_tier"`
	PremiumSubscriptionCount    int           `json:"premium_subscription_count"`
	SystemChannelFlags          int           `json:"system_channel_flags"`
	PreferredLocale             string        `json:"preferred_locale"`
	RulesChannelID              string        `json:"rules_channel_id"`
	SafetyAlertsChannelID       interface{}   `json:"safety_alerts_channel_id"`
	PublicUpdatesChannelID      string        `json:"public_updates_channel_id"`
	HubType                     interface{}   `json:"hub_type"`
	PremiumProgressBarEnabled   bool          `json:"premium_progress_bar_enabled"`
	LatestOnboardingQuestionID  interface{}   `json:"latest_onboarding_question_id"`
	IncidentsData               interface{}   `json:"incidents_data"`
	Nsfw                        bool          `json:"nsfw"`
	NsfwLevel                   int           `json:"nsfw_level"`
	EmbedEnabled                bool          `json:"embed_enabled"`
	EmbedChannelID              interface{}   `json:"embed_channel_id"`
}

// GetGuild Get discord Guild using channelID + GuildID
func GetGuild(GuildID string) Guild {
	answer := make(chan BucketRequestAnswer, 1)
	defer close(answer)

	request := BucketRequest{
		AnswerQueue: answer,
		BucketName:  "guilds",
		Url:         "/guilds/" + GuildID,
		Methode:     http.MethodGet,
		Payload:     nil,
	}
	addRequest(request)

	// Get response
	response := <-answer

	var guild Guild

	err := json.Unmarshal(response.Body, &guild)

	if err != nil {
		return Guild{}
	}

	return guild
}

func (content *Guild) SearchMember(count int, query string) []Member {
	answer := make(chan BucketRequestAnswer, 1)
	request := BucketRequest{
		AnswerQueue: answer,
		BucketName:  "guilds",
		Url:         fmt.Sprintf("guilds/%s/members/search?limit=%d&query=%s", content.ID, count, query),
		Methode:     http.MethodGet,
		Payload:     nil,
	}

	addRequest(request)

	response := <-answer
	close(answer)

	var members []Member

	err := json.Unmarshal(response.Body, &members)
	if err != nil {
		return nil
	}
	return members
}

func (content *Guild) GetMembersLists(count int, afterUserId string) []Member {
	answer := make(chan BucketRequestAnswer, 1)
	request := BucketRequest{
		AnswerQueue: answer,
		BucketName:  "guilds",
		Url:         fmt.Sprintf("guilds/%s/members?limit=%d&after=%s", content.ID, count, afterUserId),
		Methode:     http.MethodGet,
		Payload:     nil,
	}

	addRequest(request)

	response := <-answer
	close(answer)

	var members []Member

	err := json.Unmarshal(response.Body, &members)
	if err != nil {
		return nil
	}
	return members
}

func (content *Guild) UpdateGuild() {
	body, _ := json.Marshal(content)

	answer := make(chan BucketRequestAnswer, 1)

	request := BucketRequest{
		AnswerQueue: answer,
		BucketName:  "guilds",
		Url:         "/guilds/" + content.ID,
		Methode:     http.MethodPatch,
		Payload:     body,
	}

	addRequest(request)

	response := <-answer
	close(answer)

	var Guild Guild

	err := json.Unmarshal(response.Body, &Guild)
	if err != nil {
		return
	}
}

func (content *Guild) PostGuild() error {
	body, _ := json.Marshal(content)

	answer := make(chan BucketRequestAnswer, 1)
	request := BucketRequest{
		AnswerQueue: answer,
		BucketName:  "guilds",
		Url:         "/guilds",
		Methode:     http.MethodPost,
		Payload:     body,
	}

	addRequest(request)

	response := <-answer
	close(answer)

	err := json.Unmarshal(response.Body, content)
	if err != nil {
		return err
	}

	return nil
}

func (content *Guild) CreateChannel(channel Channel) bool {
	body, err := json.Marshal(channel)

	if err != nil {
		panic(err)
	}

	RequestDiscord("/guilds/"+content.ID+"/channels", http.MethodPost, "channels", body, false)

	return true
}

// GetMember Get discord Member using guildsID + UserID
func (content *Guild) GetMember(UserID string) Member {
	answer := make(chan BucketRequestAnswer, 1)
	defer close(answer)

	request := BucketRequest{
		AnswerQueue: answer,
		BucketName:  "guilds",
		Url:         "/guilds/" + content.ID + "/members/" + UserID,
		Methode:     http.MethodGet,
		Payload:     nil,
	}
	addRequest(request)

	// Get response
	response := <-answer

	var member Member

	err := json.Unmarshal(response.Body, &member)

	if err != nil {
		return Member{}
	}

	return member
}
