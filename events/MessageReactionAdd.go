package events

type MessageReactionAddEvent struct {
	V            int `json:"v"`
	UserSettings struct {
	} `json:"user_settings"`
	User struct {
		Verified      bool   `json:"verified"`
		Username      string `json:"username"`
		MfaEnabled    bool   `json:"mfa_enabled"`
		ID            string `json:"id"`
		GlobalName    any    `json:"global_name"`
		Flags         int    `json:"flags"`
		Email         any    `json:"email"`
		DisplayName   any    `json:"display_name"`
		Discriminator string `json:"discriminator"`
		Bot           bool   `json:"bot"`
		Avatar        any    `json:"avatar"`
	} `json:"user"`
	SessionType      string `json:"session_type"`
	SessionID        string `json:"session_id"`
	ResumeGatewayURL string `json:"resume_gateway_url"`
	Relationships    []any  `json:"relationships"`
	PrivateChannels  []any  `json:"private_channels"`
	Presences        []any  `json:"presences"`
	Guilds           []struct {
		Unavailable bool   `json:"unavailable"`
		ID          string `json:"id"`
	} `json:"guilds"`
	GuildJoinRequests    []any    `json:"guild_join_requests"`
	GeoOrderedRtcRegions []string `json:"geo_ordered_rtc_regions"`
	Application          struct {
		ID    string `json:"id"`
		Flags int    `json:"flags"`
	} `json:"application"`
	Trace []string `json:"_trace"`
}

var HandleMessageReactionAdd func(event MessageReactionAddEvent)

func SetHandleMessageReactionAdd(f func(event MessageReactionAddEvent)) {
	// Handle the HandleMessageReactionAdd event
	HandleMessageReactionAdd = f
}
