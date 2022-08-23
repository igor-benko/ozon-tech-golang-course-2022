package consumer

import (
	"expvar"

	"gitlab.ozon.dev/igor.benko.1991/homework/internal/pkg/counter"
)

// counters
var (
	message_processed_ok   counter.Counter
	message_processed_fail counter.Counter
)

func init() {
	expvar.Publish("message_processed_ok", expvar.Func(func() any { return message_processed_ok.Get() }))
	expvar.Publish("message_processed_fail", expvar.Func(func() any { return message_processed_fail.Get() }))
}
