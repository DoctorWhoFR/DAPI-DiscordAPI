package DiscordAPI

import (
	"azginfr/dapi/DiscordInternal"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type ChannelPermissionOverwriteType string

const (
	CPO_ROLE ChannelPermissionOverwriteType = "0"
	CPO_USER ChannelPermissionOverwriteType = "1"
)

type ChannelInvite struct {
	Code      string    `json:"code"`
	Type      int       `json:"type"`
	ExpiresAt time.Time `json:"expires_at"`
	Guild     Guild
	Channel   Channel
	Inviter   struct {
		ID               string      `json:"id"`
		Username         string      `json:"username"`
		GlobalName       interface{} `json:"global_name"`
		Avatar           string      `json:"avatar"`
		Discriminator    string      `json:"discriminator"`
		PublicFlags      int         `json:"public_flags"`
		AvatarDecoration string      `json:"avatar_decoration"`
	} `json:"inviter"`
	Uses      int       `json:"uses"`
	MaxUses   int       `json:"max_uses"`
	MaxAge    int       `json:"max_age"`
	Temporary bool      `json:"temporary"`
	CreatedAt time.Time `json:"created_at"`
}

type ChannelPermissionOverwrite struct {
	ID       string                         `json:"id"`
	Type     ChannelPermissionOverwriteType `json:"type"`
	Allow    int                            `json:"allow"`
	Deny     int                            `json:"deny"`
	AllowNew string                         `json:"allow_new"`
	DenyNew  string                         `json:"deny_new"`
}

type Channel struct {
	ID                   string                       `json:"id"`
	LastMessageID        string                       `json:"last_message_id"`
	Type                 int                          `json:"type"`
	Name                 string                       `json:"name"`
	Position             int                          `json:"position"`
	Flags                int                          `json:"flags"`
	ParentID             any                          `json:"parent_id"`
	Topic                any                          `json:"topic"`
	GuildID              string                       `json:"guild_id"`
	PermissionOverwrites []ChannelPermissionOverwrite `json:"permission_overwrites"`
	RateLimitPerUser     int                          `json:"rate_limit_per_user"`
	Nsfw                 bool                         `json:"nsfw"`
}

func (content *Channel) FollowAnnouncementChannel(channelId string) error {
	test := struct {
		WebhookChannelId string `json:"webhook_channel_id"`
	}{WebhookChannelId: channelId}

	body, err := json.Marshal(test)

	if err != nil {
		panic(err)
	}

	answer := RequestDiscord(fmt.Sprintf("/channels/%s/followers", content.ID), http.MethodPost, "channels", body, false)

	if answer.Err != nil {
		return errors.New("can't  based on technical error" + err.Error())
	}

	if answer.Res.StatusCode > http.StatusResetContent {
		return errors.New("can't  based on discord error: " + string(answer.Body))
	}

	return nil
}

// GetChannel
// https://discord.com/developers/docs/resources/channel#get-channel
func GetChannel(channelId string) (Channel, error) {
	response := RequestDiscord("/channels/"+channelId, http.MethodGet, "channels", nil, true)

	var channel Channel

	err := json.Unmarshal(response.Body, &channel)
	if err != nil {
		return Channel{}, errors.New("can't GetChannel based on technical error" + err.Error())
	}

	if response.Err != nil {
		return Channel{}, errors.New("can't GetChannel based on technical error: " + response.Err.Error())
	}

	if response.Res.StatusCode > http.StatusResetContent {
		return Channel{}, errors.New("can't GetChannel based on discord error: " + string(response.Body))
	}

	return channel, nil
}

func (content *Channel) UpdateChannel() error {
	body, err := json.Marshal(content)

	if err != nil {
		return errors.New("can't  based on technical error" + err.Error())
	}

	answer := RequestDiscord("/channels/"+content.ID, http.MethodPatch, "channels", body, false)

	if answer.Err != nil {
		return errors.New("can't  based on technical error" + answer.Err.Error())
	}

	if answer.Res.StatusCode > http.StatusResetContent {
		return errors.New("can't  based on discord error: " + string(answer.Body))
	}
	return nil
}

func (content *Channel) UpdateChannelPermission(overwrite ChannelPermissionOverwrite) error {
	body, err := json.Marshal(overwrite)

	if err != nil {
		return errors.New("can't  based on technical error" + err.Error())
	}

	answer := RequestDiscord("/channels/"+content.ID+"/permissions/"+overwrite.ID, http.MethodPut, "channels", body, false)

	if answer.Err != nil {
		return errors.New("can't  based on technical error" + answer.Err.Error())
	}

	if answer.Res.StatusCode > http.StatusResetContent {
		return errors.New("can't  based on discord error: " + string(answer.Body))
	}

	return nil
}

func (content *Channel) DeleteChannelPermission(overwrite ChannelPermissionOverwrite) error {
	body, err := json.Marshal(overwrite)

	if err != nil {
		return errors.New("can't  based on technical error" + err.Error())
	}

	answer := RequestDiscord("/channels/"+content.ID+"/permissions/"+overwrite.ID, http.MethodDelete, "channels", body, false)

	if answer.Err != nil {
		return errors.New("can't  based on technical error" + answer.Err.Error())
	}

	if answer.Res.StatusCode > http.StatusResetContent {
		return errors.New("can't  based on discord error: " + string(answer.Body))
	}

	return nil
}

func (content *Channel) DeleteChannel() error {
	answer := RequestDiscord("/channels/"+content.ID, "channels", http.MethodDelete, nil, false)

	if answer.Err != nil {
		return errors.New("can't based on technical error" + answer.Err.Error())
	}

	if answer.Res.StatusCode > http.StatusResetContent {
		return errors.New("can't based on discord error: " + string(answer.Body))
	}
	return nil
}

/*
PostMessage
*/
func (content *Channel) PostMessage(message *MessageCreate) error {
	body, _ := json.Marshal(message)

	answer := RequestDiscord("/channels/"+content.ID+"/messages", http.MethodPost, "channels", body, true)

	if answer.Err != nil {
		return errors.New("can't  based on technical error" + answer.Err.Error())
	}

	err := json.Unmarshal(answer.Body, &message)
	if err != nil {
		return err
	}

	if answer.Res.StatusCode > http.StatusResetContent {
		return errors.New("can't  based on discord error: " + string(answer.Body))
	}

	DiscordInternal.LogTrace("new message", message, "response message", string(answer.Body))

	message.Channel = *content

	return nil
}

func (content *Channel) PostMessageImage(message *MessageCreate, file string) error {
	messageJson, err := json.Marshal(message)

	if err != nil {
		panic(err)
	}

	fileContent, err := os.Open(file)
	if err != nil {
		panic(err)
		return err
	}

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	part1, _ := writer.CreateFormFile("files[0]", filepath.Base(fileContent.Name()))
	_, _ = io.Copy(part1, fileContent)
	_ = writer.WriteField("payload_json", string(messageJson))
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

	answer := RequestDiscordForm("/channels/"+content.ID+"/messages", http.MethodPost, "channels", messageJson, true, payload, writer.FormDataContentType())

	if answer.Err != nil {
		return errors.New("can't  based on technical error" + answer.Err.Error())
	}

	err = json.Unmarshal(answer.Body, &message)
	if err != nil {
		return err
	}

	if answer.Res.StatusCode > http.StatusResetContent {
		return errors.New("can't  based on discord error: " + string(answer.Body))
	}

	DiscordInternal.LogTrace("new message", message, "response message", string(answer.Body))

	message.Channel = *content

	return nil
}

func (content *Channel) GetInvites() ([]ChannelInvite, error) {
	answer := RequestDiscord("/channels/"+content.ID+"/invites", http.MethodGet, "channels", nil, true)

	if answer.Err != nil {
		return []ChannelInvite{}, errors.New("can't  based on technical error" + answer.Err.Error())
	}

	var invites []ChannelInvite

	err := json.Unmarshal(answer.Body, &invites)

	if err != nil {
		return []ChannelInvite{}, err
	}

	if answer.Res.StatusCode > http.StatusResetContent {
		return []ChannelInvite{}, errors.New("can't  based on discord error: " + string(answer.Body))
	}

	return invites, nil
}

func (content *Channel) CreateInvite(invite *ChannelInvite) error {
	body, _ := json.Marshal(invite)

	answer := RequestDiscord("/channels/"+content.ID+"/invites", http.MethodPost, "channels", body, true)

	if answer.Err != nil {
		return errors.New("can't  based on technical error" + answer.Err.Error())
	}

	err := json.Unmarshal(answer.Body, &invite)

	if err != nil {
		return err
	}

	if answer.Res.StatusCode > http.StatusResetContent {
		return errors.New("can't  based on discord error: " + string(answer.Body))
	}

	return nil
}

// GetMessage Get discord message using channelID + messageID
func (content *Channel) GetMessage(messageID string) (MessageCreate, error) {
	answer := RequestDiscord("/channels/"+content.ID+"/messages/"+messageID, http.MethodGet, "channels", nil, true)

	var message MessageCreate

	err := json.Unmarshal(answer.Body, &message)

	if err != nil {
		return MessageCreate{}, err
	}

	if answer.Err != nil {
		return MessageCreate{}, errors.New("can't  based on technical error" + answer.Err.Error())
	}

	if answer.Res.StatusCode > http.StatusResetContent {
		return MessageCreate{}, errors.New("can't  based on discord error: " + string(answer.Body))
	}
	message.Channel = *content

	DiscordInternal.LogTrace("get message", message, "response message", string(answer.Body))

	return message, nil
}

func (content *Channel) GetAllMessages(limit string) ([]MessageCreate, error) {
	answer := RequestDiscord("/channels/"+content.ID+"/messages?limit="+limit, http.MethodGet, "channels", nil, true)

	if limit == "" {
		limit = "50"
	}

	if answer.Err != nil {
		return []MessageCreate{}, errors.New("can't PostMessage based on technical error" + answer.Err.Error())
	}

	var messages []MessageCreate

	err := json.Unmarshal(answer.Body, &messages)

	if err != nil {
		return []MessageCreate{}, err
	}

	for _, m := range messages {
		m.Channel = *content
	}

	if answer.Res.StatusCode > http.StatusResetContent {
		return []MessageCreate{}, errors.New("can't  based on discord error: " + string(answer.Body))
	}

	DiscordInternal.LogTrace("get message", messages, "response message", string(answer.Body))

	return messages, nil
}

type BulkDeleteMessageType struct {
	Messages []string `json:"messages"`
}

func (content *Channel) BulkDeleteMessages(entries BulkDeleteMessageType) error {
	marshal, err := json.Marshal(entries)
	if err != nil {
		return err
	}

	answer := RequestDiscord("/channels/"+content.ID+"/messages/bulk-delete", http.MethodPost, "channels", marshal, false)

	if answer.Err != nil {
		return errors.New("can't PostMessage based on technical error" + answer.Err.Error())
	}

	if answer.Res.StatusCode > http.StatusResetContent {
		return errors.New("can't  based on discord error: " + string(answer.Body))
	}

	return nil
}

func (content *Channel) TriggerTypingIndicator() error {
	answer := RequestDiscord(fmt.Sprintf("/channels/%s/typing", content.ID), http.MethodPost, "channels", nil, false)

	if answer.Err != nil {
		return errors.New("can't  based on technical error" + answer.Err.Error())
	}

	if answer.Res.StatusCode > http.StatusResetContent {
		return errors.New("can't  based on discord error: " + string(answer.Body))
	}

	return nil
}

func (content *Channel) GetPinnedMessages() ([]MessageCreate, error) {
	answer := RequestDiscord(fmt.Sprintf("/channels/%s/pins", content.ID), http.MethodGet, "channels", nil, false)

	if answer.Err != nil {
		return nil, errors.New("can't  based on technical error" + answer.Err.Error())
	}

	if answer.Res.StatusCode > http.StatusResetContent {
		return nil, errors.New("can't  based on discord error: " + string(answer.Body))
	}

	var messages []MessageCreate

	err := json.Unmarshal(answer.Body, &messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

type StartThreadPayload struct {
	Name                string `json:"name"`                            // string	1-100 character channel name
	AutoArchiveDuration int    `json:"auto_archive_duration,omitempty"` // integer	the type of thread to create
	RateLimitPerUser    int    `json:"rate_limit_per_user,omitempty"`   // amount of seconds a user has to wait before sending another message (0-21600)
}

/*
StartThread
Start Thread without Message
POST/channels/{channel.id}/threads
Creates a new thread that is not connected to an existing message. Returns a channel on success, and a 400 BAD REQUEST on invalid parameters. Fires a Thread Create Gateway event.

This endpoint supports the X-Audit-Log-Reason header.
JSON Params
FIELD	TYPE	DESCRIPTION
name
auto_archive_duration?	integer	the thread will stop showing in the channel list after auto_archive_duration minutes of inactivity, can be set to: 60, 1440, 4320, 10080
type?*
invitable?	boolean	whether non-moderators can add other non-moderators to a thread; only available when creating a private thread
rate_limit_per_user?	?integer
* type currently defaults to PRIVATE_THREAD in order to match the behavior when thread documentation was first published. In a future API version this will be changed to be a required field, with no default.
*/
func (content *Channel) StartThread(thread StartThreadPayload) (Channel, error) {
	marshal, err := json.Marshal(thread)
	if err != nil {
		return Channel{}, err
	}

	answer := RequestDiscord(fmt.Sprintf("/channels/%s/threads", content.ID), http.MethodPost, "channels", marshal, true)

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

/*
JoinThreads
Join Thread
PUT/channels/{channel.id}/thread-members/@me
Adds the current user to a thread. Also requires the thread is not archived. Returns a 204 empty response on success. Fires a Thread Members UpdateCommand and a Thread Create Gateway event.
*/
func (content *Channel) JoinThreads() error {
	answer := RequestDiscord(fmt.Sprintf("/channels/%s/thread-members/@me", content.ID), http.MethodPut, "channels", nil, false)

	if answer.Err != nil {
		return errors.New("can't  based on technical error" + answer.Err.Error())
	}

	if answer.Res.StatusCode > http.StatusResetContent {
		return errors.New("can't  based on discord error: " + string(answer.Body))
	}

	return nil
}

/*
LeaveThreads
Leave Thread
DELETE/channels/{channel.id}/thread-members/@me
Removes the current user from a thread.
Also requires the thread is not archived.
Returns a 204 empty response on success. Fires a Thread Members UpdateCommand Gateway event.
*/
func (content *Channel) LeaveThreads() error {
	answer := RequestDiscord(fmt.Sprintf("/channels/%s/thread-members/@me", content.ID), http.MethodDelete, "channels", nil, false)

	if answer.Err != nil {
		return errors.New("can't  based on technical error" + answer.Err.Error())
	}

	if answer.Res.StatusCode > http.StatusResetContent {
		return errors.New("can't  based on discord error: " + string(answer.Body))
	}

	return nil
}

/*
AddThreadMember
Add Thread Member
PUT/channels/{channel.id}/thread-members/{user.id}
Adds another member to a thread.
Requires the ability to send messages in the thread. Also requires the thread is not archived.
Returns a 204 empty response if the member is successfully added or was already a member of the thread.
Fires a Thread Members UpdateCommand Gateway event.
*/
func (content *Channel) AddThreadMember(memberId string) error {
	answer := RequestDiscord(fmt.Sprintf("/channels/%s/thread-members/%s", content.ID, memberId), http.MethodPut, "channels", nil, false)

	if answer.Err != nil {
		return errors.New("can't  based on technical error" + answer.Err.Error())
	}

	if answer.Res.StatusCode > http.StatusResetContent {
		return errors.New("can't  based on discord error: " + string(answer.Body))
	}

	return nil
}

/*
DeleteThreadMember
Remove Thread Member
DELETE/channels/{channel.id}/thread-members/{user.id}
Removes another member from a thread.
Requires the MANAGE_THREADS permission, or the creator of the thread if it is a PRIVATE_THREAD.
Also requires the thread is not archived. Returns a 204 empty response on success.
Fires a Thread Members UpdateCommand Gateway event.
*/
func (content *Channel) DeleteThreadMember(memberId string) error {
	answer := RequestDiscord(fmt.Sprintf("/channels/%s/thread-members/%s", content.ID, memberId), http.MethodDelete, "channels", nil, false)

	if answer.Err != nil {
		return errors.New("can't  based on technical error" + answer.Err.Error())
	}

	if answer.Res.StatusCode > http.StatusResetContent {
		return errors.New("can't  based on discord error: " + string(answer.Body))
	}

	return nil
}

/*
GetThreadMember
TODO: fix this functions return nothing and should return an thread member
Get Thread Member
GET/channels/{channel.id}/thread-members/{user.id}
Returns a thread member object for the specified user if they are a member of the thread, returns a 404 response otherwise.

When with_member is set to true, the thread member object will include a member field containing a guild member object.

Query String Params
FIELD	TYPE	DESCRIPTION
with_member?	boolean	Whether to include a guild member object for the thread member
*/
func (content *Channel) GetThreadMember(memberId string) error {
	answer := RequestDiscord(fmt.Sprintf("/channels/%s/thread-members/%s", content.ID, memberId), http.MethodGet, "channels", nil, false)

	if answer.Err != nil {
		return errors.New("can't  based on technical error" + answer.Err.Error())
	}

	if answer.Res.StatusCode > http.StatusResetContent {
		return errors.New("can't  based on discord error: " + string(answer.Body))
	}

	return nil
}

/*
GetThreadMembers
List Thread Members
GET/channels/{channel.id}/thread-members
Starting in API v11, this endpoint will always return paginated results. Paginated results can be enabled before API v11 by setting with_member to true. Read the changelog for details.
Returns array of thread members objects that are members of the thread.

When with_member is set to true, the results will be paginated and each thread member object will include a member field containing a guild member object.

This endpoint is restricted according to whether the GUILD_MEMBERS Privileged Intent is enabled for your application.

# Query String Params

# FIELD	TYPE	DESCRIPTION
  - with_member?	boolean	Whether to include a guild member object for each thread member
  - after?	snowflake	Get thread members after this user ID
  - limit?	integer	Max number of thread members to return (1-100). Defaults to 100.
*/
func (content *Channel) GetThreadMembers() ([]ThreadMember, error) {
	answer := RequestDiscord(fmt.Sprintf("/channels/%s/thread-members", content.ID), http.MethodGet, "channels", nil, true)

	if answer.Err != nil {
		return nil, errors.New("can't  based on technical error" + answer.Err.Error())
	}

	if answer.Res.StatusCode > http.StatusResetContent {
		return nil, errors.New("can't  based on discord error: " + string(answer.Body))
	}

	var members []ThreadMember

	err := json.Unmarshal(answer.Body, &members)
	if err != nil {
		return nil, err
	}

	return members, nil
}

type ListPublicThread struct {
	Threads []Channel `json:"threads"`
	Members []Member  `json:"members"`
	HasMore bool      `json:"has_more"`
}

func (content *Channel) ListPublicArchivedThreads() (ListPublicThread, error) {
	answer := RequestDiscord(fmt.Sprintf("/channels/%s/threads/archived/public", content.ID), http.MethodGet, "channels", nil, true)

	if answer.Err != nil {
		return ListPublicThread{}, errors.New("can't  based on technical error" + answer.Err.Error())
	}

	if answer.Res.StatusCode > http.StatusResetContent {
		return ListPublicThread{}, errors.New("can't  based on discord error: " + string(answer.Body))
	}

	var listPublicThread ListPublicThread

	err := json.Unmarshal(answer.Body, &listPublicThread)

	if err != nil {
		return ListPublicThread{}, err
	}

	return listPublicThread, nil
}
