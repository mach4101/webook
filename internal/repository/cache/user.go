package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/mach4101/geek_go_camp/webook/internal/domain"
)

var ErrKeyNotExist = redis.Nil

type UserCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewUserCache(client redis.Cmdable) *UserCache {
	return &UserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

// 构造key
func (cache *UserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}

// 若无数据，返回一个特定的error
func (cache *UserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.key(id)

	val, err := cache.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}

	var u domain.User
	err = json.Unmarshal(val, &u)
	return u, err
}

// 放置key
func (cache *UserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}

	key := cache.key(u.Id)
	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}
