package utils

import (
	"time"

	tb "gopkg.in/tucnak/telebot.v2"

	"github.com/kitanoyoru/kitaDriveBot/pkg/logger"
)

func logDuration(f func(*tb.Message)) func(*tb.Message) {
	return func(m *tb.Message) {
		start := time.Now()
		f(m)
		diff := time.Now().Sub(start)
		if diff.Seconds() > 1 {
			logger.Warnf("Took %s time to complete update processing", time.Time{}.Add(diff).Format("04:05:000"))
		}
	}
}
