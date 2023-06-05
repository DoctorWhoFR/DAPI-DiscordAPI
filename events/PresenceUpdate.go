package events

type PresenceUpdate struct {
	User struct {
		Id string `json:"id"`
	} `json:"user"`
	Status       string      `json:"status"`
	Roles        []string    `json:"roles"`
	PremiumSince interface{} `json:"premium_since"`
	Nick         interface{} `json:"nick"`
	GuildId      string      `json:"guild_id"`
	Game         struct {
		Type      int         `json:"type"`
		SessionId interface{} `json:"session_id"`
		Name      string      `json:"name"`
		Id        string      `json:"id"`
		Emoji     struct {
			Name string `json:"name"`
		} `json:"emoji"`
		CreatedAt int64 `json:"created_at"`
	} `json:"game"`
	ClientStatus struct {
		Mobile string `json:"mobile"`
	} `json:"client_status"`
	Broadcast  interface{} `json:"broadcast"`
	Activities []struct {
		Type  int    `json:"type"`
		Name  string `json:"name"`
		Id    string `json:"id"`
		Emoji struct {
			Name string `json:"name"`
		} `json:"emoji"`
		CreatedAt int64 `json:"created_at"`
	} `json:"activities"`
}

var HandlePresenceUpdate func(event PresenceUpdate)

func SetHandlePresenceUpdate(f func(event PresenceUpdate)) {
	// Handle the HandleMessageReactionAdd event
	HandlePresenceUpdate = f
}
