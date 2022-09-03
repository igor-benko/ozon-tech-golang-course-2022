//go:generate mockgen -source ./cache.go -destination ./cache_mock.go -package=cache
package cache

import (
	"context"
	"encoding/json"
	"expvar"

	"gitlab.ozon.dev/igor.benko.1991/homework/internal/pkg/counter"
)

type Publication struct {
	Key     string
	Payload []byte
}

func (p Publication) MarshalBinary() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Publication) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, p)
}

type CacheClient interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte) error
	Publish(ctx context.Context, pubData *Publication) error
	Subscribe(ctx context.Context) <-chan *Publication
}

// counters
var (
	cache_hit  counter.Counter
	cache_miss counter.Counter
)

func init() {
	expvar.Publish("cache_hit", expvar.Func(func() any { return cache_hit.Get() }))
	expvar.Publish("cache_miss", expvar.Func(func() any { return cache_miss.Get() }))
}
