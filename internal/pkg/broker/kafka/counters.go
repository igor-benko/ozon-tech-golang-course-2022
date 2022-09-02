package kafka

import (
	"expvar"

	"gitlab.ozon.dev/igor.benko.1991/homework/internal/pkg/counter"
)

// counters
var (
	send_to_kafka_ok   counter.Counter
	send_to_kafka_fail counter.Counter
	got_from_kafka     counter.Counter
)

func init() {
	expvar.Publish("send_to_kafka_ok", expvar.Func(func() any { return send_to_kafka_ok.Get() }))
	expvar.Publish("send_to_kafka_fail", expvar.Func(func() any { return send_to_kafka_fail.Get() }))
	expvar.Publish("got_from_kafka", expvar.Func(func() any { return got_from_kafka.Get() }))
}
