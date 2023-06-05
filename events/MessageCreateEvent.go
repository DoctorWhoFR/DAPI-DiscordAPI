package events

import "test/dapi/restapi"

var HandleMessageCreate func(event restapi.MessageCreate)

func SetHandleMessageCreate(f func(event restapi.MessageCreate)) {
	// Handle the HandleMessageReactionAdd event
	HandleMessageCreate = f
}
