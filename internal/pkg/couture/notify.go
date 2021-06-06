package couture

import (
	"github.com/gen2brain/beeep"
	"go.uber.org/ratelimit"
	"time"
)

const notificationsPerMinute = 10
const notificationsToBuffer = 5

var osNotificationLimiter = ratelimit.New(
	notificationsPerMinute,
	ratelimit.Per(time.Minute),
	ratelimit.WithSlack(notificationsToBuffer),
)

// NotifyOS ...
func NotifyOS(title string, message string) error {
	const noIcon = ""
	osNotificationLimiter.Take()
	return beeep.Notify(title, message, noIcon)
}
