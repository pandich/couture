package manager

import (
	"github.com/gen2brain/beeep"
	"go.uber.org/ratelimit"
	"time"
)

// this uses the beeep project to sent OS-speciic notifications.
// to prevent the user being flooded, a rate limiter is in place.

// notificationsPerMinute to the OS.
const notificationsPerMinute = 6

// osNotificationRateLimiter to prevent the OS from being flooded.
var osNotificationRateLimiter = ratelimit.New(notificationsPerMinute, ratelimit.Per(time.Minute))

func notifyOS(title string, message string) error {
	// TODO can we specify a Couture icon?
	const noIcon = ""

	osNotificationRateLimiter.Take()

	return beeep.Notify(title, message, noIcon)
}
