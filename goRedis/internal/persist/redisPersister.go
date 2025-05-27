package persist

import (
	"context"
	"encoding/json"
	"github.com/folivorra/studyDir/tree/develop/goRedis/internal/model"
	"github.com/redis/go-redis/v9"
)

type RedisPersister struct {
	rdb *redis.Client
	key string
}

func NewRedisPersister(rdb *redis.Client, key string) *RedisPersister {
	return &RedisPersister{rdb: rdb, key: key}
}

func (p *RedisPersister) Dump(data map[int]model.Item) error {
	ctx := context.Background()

	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err = p.rdb.Set(ctx, p.key, string(bytes), 0).Err(); err != nil {
		return err
	}

	return nil
}

func (p *RedisPersister) Load() (map[int]model.Item, error) {
	ctx := context.Background()

	bytes, err := p.rdb.Get(ctx, p.key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	result := make(map[int]model.Item)
	if err = json.Unmarshal([]byte(bytes), &result); err != nil {
		return nil, err
	}

	return result, nil
}
