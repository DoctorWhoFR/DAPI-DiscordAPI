package DiscordInternal

import "time"

const BaseSleepTime = time.Millisecond * 10

func SimpleSleep() {
	time.Sleep(time.Millisecond)
}
