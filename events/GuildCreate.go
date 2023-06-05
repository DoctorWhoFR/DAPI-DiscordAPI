package events

import (
	"log"
	"time"
)

type GuildCreate struct {
	Stickers    []interface{} `json:"stickers"`
	Unavailable bool          `json:"unavailable"`
	Roles       []struct {
		Version      int64       `json:"version"`
		UnicodeEmoji interface{} `json:"unicode_emoji"`
		Tags         struct {
			BotId             string      `json:"bot_id,omitempty"`
			PremiumSubscriber interface{} `json:"premium_subscriber"`
		} `json:"tags"`
		Position       int         `json:"position"`
		PermissionsNew string      `json:"permissions_new"`
		Permissions    int         `json:"permissions"`
		Name           string      `json:"name"`
		Mentionable    bool        `json:"mentionable"`
		Managed        bool        `json:"managed"`
		Id             string      `json:"id"`
		Icon           interface{} `json:"icon"`
		Hoist          bool        `json:"hoist"`
		Flags          int         `json:"flags"`
		Color          int         `json:"color"`
	} `json:"roles"`
	Banner   interface{} `json:"banner"`
	Channels []struct {
		Version              int64       `json:"version"`
		UserLimit            int         `json:"user_limit,omitempty"`
		Type                 int         `json:"type"`
		RtcRegion            interface{} `json:"rtc_region"`
		RateLimitPerUser     int         `json:"rate_limit_per_user,omitempty"`
		Position             int         `json:"position"`
		PermissionOverwrites []struct {
			Type     string `json:"type"`
			Id       string `json:"id"`
			DenyNew  string `json:"deny_new"`
			Deny     int    `json:"deny"`
			AllowNew string `json:"allow_new"`
			Allow    int    `json:"allow"`
		} `json:"permission_overwrites"`
		ParentId                      *string     `json:"parent_id"`
		Nsfw                          bool        `json:"nsfw,omitempty"`
		Name                          string      `json:"name"`
		LastMessageId                 *string     `json:"last_message_id,omitempty"`
		Id                            string      `json:"id"`
		Flags                         int         `json:"flags"`
		Bitrate                       int         `json:"bitrate,omitempty"`
		Topic                         *string     `json:"topic,omitempty"`
		ThemeColor                    interface{} `json:"theme_color"`
		IconEmoji                     interface{} `json:"icon_emoji"`
		DefaultThreadRateLimitPerUser int         `json:"default_thread_rate_limit_per_user,omitempty"`
	} `json:"channels"`
	VanityUrlCode interface{} `json:"vanity_url_code"`
	Splash        string      `json:"splash"`
	Large         bool        `json:"large"`
	MemberCount   int         `json:"member_count"`
	Presences     []struct {
		User struct {
			Id string `json:"id"`
		} `json:"user"`
		Status string `json:"status"`
		Game   *struct {
			Type       int     `json:"type"`
			SessionId  *string `json:"session_id"`
			Name       string  `json:"name"`
			Id         string  `json:"id"`
			CreatedAt  int64   `json:"created_at"`
			Timestamps struct {
				Start int64 `json:"start"`
				End   int64 `json:"end"`
			} `json:"timestamps,omitempty"`
			SyncId string `json:"sync_id,omitempty"`
			State  string `json:"state,omitempty"`
			Party  struct {
				Id string `json:"id"`
			} `json:"party,omitempty"`
			Flags   int    `json:"flags,omitempty"`
			Details string `json:"details,omitempty"`
			Assets  struct {
				LargeText  string `json:"large_text"`
				LargeImage string `json:"large_image"`
			} `json:"assets,omitempty"`
		} `json:"game"`
		ClientStatus struct {
			Web     string `json:"web,omitempty"`
			Desktop string `json:"desktop,omitempty"`
			Mobile  string `json:"mobile,omitempty"`
		} `json:"client_status"`
		Broadcast  interface{} `json:"broadcast"`
		Activities []struct {
			Type       int    `json:"type"`
			Name       string `json:"name"`
			Id         string `json:"id"`
			CreatedAt  int64  `json:"created_at"`
			Timestamps struct {
				Start int64 `json:"start"`
				End   int64 `json:"end"`
			} `json:"timestamps,omitempty"`
			SyncId    string `json:"sync_id,omitempty"`
			State     string `json:"state,omitempty"`
			SessionId string `json:"session_id,omitempty"`
			Party     struct {
				Id string `json:"id"`
			} `json:"party,omitempty"`
			Flags   int    `json:"flags,omitempty"`
			Details string `json:"details,omitempty"`
			Assets  struct {
				LargeText  string `json:"large_text"`
				LargeImage string `json:"large_image"`
			} `json:"assets,omitempty"`
		} `json:"activities"`
	} `json:"presences"`
	OwnerId                     string        `json:"owner_id"`
	HomeHeader                  interface{}   `json:"home_header"`
	PremiumProgressBarEnabled   bool          `json:"premium_progress_bar_enabled"`
	Description                 interface{}   `json:"description"`
	VoiceStates                 []interface{} `json:"voice_states"`
	ExplicitContentFilter       int           `json:"explicit_content_filter"`
	Id                          string        `json:"id"`
	SystemChannelId             string        `json:"system_channel_id"`
	DiscoverySplash             interface{}   `json:"discovery_splash"`
	HubType                     interface{}   `json:"hub_type"`
	GuildScheduledEvents        []interface{} `json:"guild_scheduled_events"`
	PublicUpdatesChannelId      string        `json:"public_updates_channel_id"`
	SystemChannelFlags          int           `json:"system_channel_flags"`
	Threads                     []interface{} `json:"threads"`
	JoinedAt                    time.Time     `json:"joined_at"`
	MaxStageVideoChannelUsers   int           `json:"max_stage_video_channel_users"`
	LatestOnboardingQuestionId  interface{}   `json:"latest_onboarding_question_id"`
	StageInstances              []interface{} `json:"stage_instances"`
	DefaultMessageNotifications int           `json:"default_message_notifications"`
	EmbeddedActivities          []interface{} `json:"embedded_activities"`
	Lazy                        bool          `json:"lazy"`
	PreferredLocale             string        `json:"preferred_locale"`
	MaxVideoChannelUsers        int           `json:"max_video_channel_users"`
	Members                     []struct {
		User struct {
			Username         string  `json:"username"`
			PublicFlags      int     `json:"public_flags"`
			Id               string  `json:"id"`
			GlobalName       *string `json:"global_name"`
			DisplayName      *string `json:"display_name"`
			Discriminator    string  `json:"discriminator"`
			Bot              bool    `json:"bot"`
			AvatarDecoration *string `json:"avatar_decoration"`
			Avatar           *string `json:"avatar"`
		} `json:"user"`
		Roles                      []string    `json:"roles"`
		PremiumSince               *time.Time  `json:"premium_since"`
		Pending                    bool        `json:"pending"`
		Nick                       interface{} `json:"nick"`
		Mute                       bool        `json:"mute"`
		JoinedAt                   time.Time   `json:"joined_at"`
		Flags                      int         `json:"flags"`
		Deaf                       bool        `json:"deaf"`
		CommunicationDisabledUntil interface{} `json:"communication_disabled_until"`
		Avatar                     *string     `json:"avatar"`
	} `json:"members"`
	RulesChannelId           string      `json:"rules_channel_id"`
	SafetyAlertsChannelId    interface{} `json:"safety_alerts_channel_id"`
	ApplicationId            interface{} `json:"application_id"`
	Region                   string      `json:"region"`
	NsfwLevel                int         `json:"nsfw_level"`
	MfaLevel                 int         `json:"mfa_level"`
	AfkChannelId             string      `json:"afk_channel_id"`
	PremiumTier              int         `json:"premium_tier"`
	AfkTimeout               int         `json:"afk_timeout"`
	IncidentsData            interface{} `json:"incidents_data"`
	PremiumSubscriptionCount int         `json:"premium_subscription_count"`
	Features                 []string    `json:"features"`
	MaxMembers               int         `json:"max_members"`
	Name                     string      `json:"name"`
	GuildHashes              struct {
		Version int `json:"version"`
		Roles   struct {
			Omitted bool   `json:"omitted"`
			Hash    string `json:"hash"`
		} `json:"roles"`
		Metadata struct {
			Omitted bool   `json:"omitted"`
			Hash    string `json:"hash"`
		} `json:"metadata"`
		Channels struct {
			Omitted bool   `json:"omitted"`
			Hash    string `json:"hash"`
		} `json:"channels"`
	} `json:"guild_hashes"`
	Nsfw                     bool          `json:"nsfw"`
	Emojis                   []interface{} `json:"emojis"`
	ApplicationCommandCounts struct {
		Field1 int `json:"3"`
		Field2 int `json:"2"`
		Field3 int `json:"1"`
	} `json:"application_command_counts"`
	VerificationLevel int    `json:"verification_level"`
	Icon              string `json:"icon"`
}

func HandleGuildCreate(data GuildCreate) {
	// Handle the GuildCreate event
	log.Println("Received GuildCreate event:", data.Name)
}
