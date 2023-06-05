package events

import (
	"encoding/json"
	"errors"
	"os"
	"test/dapi/internal"
	"test/dapi/restapi"
)

type DiscordApplicationCommandOptionType int

const (
	DACOT_SUB_COMMAND DiscordApplicationCommandOptionType = 1 + iota
	DACOT_SUB_COMMAND_GROUP
	DACOT_STRING
	DACOT_INTEGER
	DACOT_BOOLEAN
	DACOT_USER
	DACOT_CHANNEL
	DACOT_ROLE
	DACOT_MENTIONABLE
	DACOT_NUMBER
	DACOT_ATTACHMENT
)

type DiscordCommandOption struct {
	Type                     DiscordApplicationCommandOptionType `json:"type"`
	Name                     string                              `json:"name"`
	NameLocalizations        interface{}                         `json:"name_localizations"`
	Description              string                              `json:"description"`
	DescriptionLocalizations interface{}                         `json:"description_localizations"`
	Required                 bool                                `json:"required,omitempty"`
	Choices                  []struct {
		Name              string      `json:"name"`
		NameLocalizations interface{} `json:"name_localizations"`
		Value             string      `json:"value"`
	} `json:"choices,omitempty"`
}

type DiscordApplicationCommandType int

const (
	DACT_CHAT_INPUT DiscordApplicationCommandType = 1 + iota
	DACT_USER
	DACT_MESSAGE
)

type DiscordCommand struct {
	ID                       string                        `json:"id"`
	ApplicationID            string                        `json:"application_id"`
	Version                  string                        `json:"version"`
	DefaultMemberPermissions interface{}                   `json:"default_member_permissions"`
	Type                     DiscordApplicationCommandType `json:"type"`
	Name                     string                        `json:"name"`
	NameLocalizations        interface{}                   `json:"name_localizations"`
	Description              string                        `json:"description"`
	DescriptionLocalizations interface{}                   `json:"description_localizations"`
	DmPermission             bool                          `json:"dm_permission"`
	Contexts                 interface{}                   `json:"contexts"`
	Options                  []DiscordCommandOption        `json:"options"`
	Nsfw                     bool                          `json:"nsfw"`
}

func (command DiscordCommand) save() error {
	body, err := json.Marshal(command)

	if err != nil {
		internal.LogDebug("error creating command", err)
	}

	response := restapi.RequestDiscord("/applications/1106890236757278770/commands", "POST", "applications", body, false)

	if response.Res.StatusCode >= 300 {
		internal.LogError(response)
		return errors.New("can't create command")
	}
	return nil
}

func (command DiscordCommand) update(commandId string) error {
	body, err := json.Marshal(command)

	if err != nil {
		internal.LogDebug("error creating command", err)
	}

	response := restapi.RequestDiscord("applications/1106890236757278770/commands/"+commandId, "PATCH", "applications", body, false)

	if response.Res.StatusCode >= 300 {
		internal.LogDebug("log")
		return errors.New("can't create command")
	}
	return nil
}

type DiscordCommands struct {
	Commands []DiscordCommand
}

type Command struct {
	Execution func(interaction restapi.InteractionCommand)
	DCommand  DiscordCommand
	Name      string
}

var CommandsLists = make(map[string]Command, 0)

func findCommand(key string, commands []DiscordCommand) (DiscordCommand, bool) {
	for _, command := range commands {
		if command.Name == key {
			return command, true
		}
	}
	return DiscordCommand{}, false
}

func AddCommand(command Command) {
	internal.LogInfo("loading command", command.Name)

	// technique add command to list
	CommandsLists[command.Name] = command

	//// REGISTER COMMAND & UPDATE COMMAND
	if os.Getenv("CHECK_COMMAND") != "true" {
		internal.LogInfo("CHECK_COMMAND", false)
		return
	}

	// Get alls actuals commands
	answer := restapi.RequestDiscord("/applications/1106890236757278770/commands", "GET", "applications", nil, true)

	var commands []DiscordCommand

	err := json.Unmarshal(answer.Body, &commands)
	if err != nil {
		return
	}

	updateCommand := os.Getenv("UPDATE_COMMAND")

	discordCommand, present := findCommand(command.Name, commands)

	if !present {
		internal.LogInfo(command.Name, "Request is not present in application global command creating it.")
		err := command.DCommand.save()
		if err != nil {
			internal.LogError("cant' create", err)
		}
		internal.LogInfo("created")
	}

	if updateCommand == "true" {
		err := command.DCommand.update(discordCommand.ID)
		if err != nil {
			internal.LogError("can't update")
			return
		}
		internal.LogInfo("updated")
	}

}

var HandlerInteractionCreate func(interaction restapi.InteractionCommand)

func SetHandlerInteractionCreate(f func(interaction restapi.InteractionCommand)) {

	HandlerInteractionCreate = f
}

var HandlerInteractionButton func(interaction restapi.InteractionButton)

func SetHandlerInteractionButton(f func(interaction restapi.InteractionButton)) {
	HandlerInteractionButton = f
}
