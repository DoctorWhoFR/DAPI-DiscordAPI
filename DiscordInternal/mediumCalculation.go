package DiscordInternal

import "time"

var HandlingTimeDiscord = []time.Duration{0}
var HandlingTimeBot = []time.Duration{0}

// MediumValueAPI return medium values of all api call
func MediumValueAPI() int64 {
	var total int64
	for _, timed := range HandlingTimeDiscord {
		total = total + timed.Milliseconds()
	}

	return total / int64(len(HandlingTimeDiscord))
}

// MediumValueBOT return medium values of all bot process time
func MediumValueBOT() int64 {
	var total int64

	for _, timed := range HandlingTimeBot {
		total = total + timed.Milliseconds()
	}

	return total / int64(len(HandlingTimeBot))
}

// MediumValueWithoutAPI comparison from discord call time and real bot process time
func MediumValueWithoutAPI() int64 {
	mediumApi := MediumValueAPI()
	mediumBot := MediumValueBOT()

	return mediumBot - mediumApi
}
