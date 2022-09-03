package cache

import (
	"context"
	"time"

	redis "github.com/go-redis/redis/v9"
	"gitlab.ozon.dev/igor.benko.1991/homework/internal/config"
	"gitlab.ozon.dev/igor.benko.1991/homework/pkg/logger"
)

type RedisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Publish(ctx context.Context, channel string, message interface{}) *redis.IntCmd
	Subscribe(ctx context.Context, channels ...string) *redis.PubSub
}

type redisClient struct {
	client          RedisClient
	defaultExpireMs uint64
	channel         string
}

func NewRedisCache(cfg config.Config, client RedisClient) *redisClient {
	return &redisClient{
		defaultExpireMs: cfg.Cache.ExpireMs,
		channel:         cfg.Cache.Channel,
		client:          client,
	}
}

func (r *redisClient) Get(ctx context.Context, key string) ([]byte, error) {
	res := r.client.Get(ctx, key)
	if res.Err() != nil {
		if res.Err() == redis.Nil {
			cache_miss.Inc()
		}

		return nil, res.Err()
	}

	cache_hit.Inc()

	return res.Bytes()
}

func (r *redisClient) Set(ctx context.Context, key string, value []byte) error {
	return r.client.
		Set(ctx, key, value, time.Duration(r.defaultExpireMs)*time.Millisecond).
		Err()
}

func (r *redisClient) Publish(ctx context.Context, data *Publication) error {
	return r.client.Publish(ctx, r.channel, data).Err()
}

func (r *redisClient) Subscribe(ctx context.Context) <-chan *Publication {
	resultCh := make(chan *Publication, 1)
	sub := r.client.Subscribe(ctx, r.channel)
	msgCh := sub.Channel()

	go func() {
	loop:
		for {
			select {
			case <-ctx.Done():
				sub.Close()
				logger.Warnf("Context was cancelled")
				break loop
			case msg := <-msgCh:
				pubData := &Publication{}
				if err := pubData.UnmarshalBinary([]byte(msg.Payload)); err != nil {
					logger.Warnf("ReceiveMessage failed: %s", err.Error())
					continue
				}

				resultCh <- pubData
			}

		}

		close(resultCh)
	}()

	return resultCh
}
