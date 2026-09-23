package state

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type Redis struct{ client *redis.Client }

func New(addr string) *Redis  { return &Redis{client: redis.NewClient(&redis.Options{Addr: addr})} }
func (s *Redis) Close() error { return s.client.Close() }
func (s *Redis) Ready(ctx context.Context, key string, cooldown time.Duration) (bool, error) {
	ok, err := s.client.SetNX(ctx, "notified:"+key, time.Now().UTC().Format(time.RFC3339), cooldown).Result()
	return ok, err
}
func (s *Redis) Clear(ctx context.Context, key string) error {
	return s.client.Del(ctx, "notified:"+key).Err()
}
