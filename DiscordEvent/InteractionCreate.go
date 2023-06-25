package DiscordEvent

import (
	"azginfr/dapi/DiscordAPI"
	"azginfr/dapi/DiscordInternal"
	"encoding/json"
	"errors"
	"os"
)

type DiscordApplicationCommandOptionType int

const (
	DiscordApplicationCommandOptionTypeSubCommand DiscordApplicationCommandOptionType = 1 + iota
	DiscordApplicationCommandOptionTypeSubCommandGroup
	DiscordApplicationCommandOptionTypeString
	DiscordApplicationCommandOptionTypeInteger
	DiscordApplicationCommandOptionTypeBoolean
	DiscordApplicationCommandOptionTypeUser
	DiscordApplicationCommandOptionTypeChannel
	DiscordApplicationCommandOptionTypeRole
	DiscordApplicationCommandOptionTypeMentionable
	DiscordApplicationCommandOptionTypeNumber
	DiscordApplicationCommandOptionTypeAttachment
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
	Options []DiscordCommandOption `json:"options,omitempty"`
}

type DiscordApplicationCommandType int

const (
	DiscordApplicationCommandTypeChatInput DiscordApplicationCommandType = 1 + iota
	DiscordApplicationCommandTypeUser
	DiscordApplicationCommandTypeMessage
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

func (command *DiscordCommand) SaveCommand() error {
	body, err := json.Marshal(command)

	if err != nil {
		DiscordInternal.LogDebug("error creating command", err)
	}

	response := DiscordAPI.RequestDiscord("/applications/1106890236757278770/commands", "POST", "applications", body, false)

	if response.Res.StatusCode >= 300 {
		DiscordInternal.LogError(response)
		return errors.New("can't create command")
	}
	return nil
}

func (command *DiscordCommand) UpdateCommand(commandId string) error {
	body, err := json.Marshal(command)

	if err != nil {
		DiscordInternal.LogDebug("error creating command", err)
	}

	response := DiscordAPI.RequestDiscord("applications/1106890236757278770/commands/"+commandId, "PATCH", "applications", body, false)

	if response.Res.StatusCode >= 300 {
		DiscordInternal.LogDebug("log")
		return errors.New("can't create command")
	}
	return nil
}

type BotCommand struct {
	Execution func(interaction DiscordAPI.InteractionCommand)
	DCommand  DiscordCommand
	Name      string
}

var CommandsLists = make(map[string]BotCommand, 0)

// FindDiscordCommand TODO: find why ^^ there is DiscordCommand instead of BotCommand
func FindDiscordCommand(key string, commands []DiscordCommand) (DiscordCommand, bool) {
	for _, command := range commands {
		if command.Name == key {
			return command, true
		}
	}
	return DiscordCommand{}, false
}

func AddDiscordCommand(command BotCommand) {
	DiscordInternal.LogInfo("loading command", command.Name)

	// add command to list
	CommandsLists[command.Name] = command

	//// REGISTER COMMAND & UPDATE COMMAND
	if os.Getenv("CHECK_COMMAND") != "true" {
		DiscordInternal.LogInfo("CHECK_COMMAND", false)
		return
	}

	// Get alls actuals commands
	answer := DiscordAPI.RequestDiscord("/applications/1106890236757278770/commands", "GET", "applications", nil, true)

	var commands []DiscordCommand

	err := json.Unmarshal(answer.Body, &commands)
	if err != nil {
		return
	}

	updateCommand := os.Getenv("UPDATE_COMMAND")

	discordCommand, present := FindDiscordCommand(command.Name, commands)

	if !present {
		DiscordInternal.LogInfo(command.Name, "Request is not present in application global command creating it.")
		err := command.DCommand.SaveCommand()
		if err != nil {
			DiscordInternal.LogError("cant' create", err)
		}
		DiscordInternal.LogInfo("created")
	}

	if updateCommand == "true" {
		err := command.DCommand.UpdateCommand(discordCommand.ID)
		if err != nil {
			DiscordInternal.LogError("can't UpdateCommand")
			return
		}
		DiscordInternal.LogInfo("updated")
	}

}

var HandlerInteractionCommandEvent func(interaction DiscordAPI.InteractionCommand)

func SetHandlerInteractionCommandEvent(f func(interaction DiscordAPI.InteractionCommand)) {
	HandlerInteractionCommandEvent = f
}
