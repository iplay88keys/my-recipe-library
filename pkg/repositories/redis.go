package repositories

import (
	"strconv"
	"time"

	"github.com/go-redis/redis"

	"github.com/iplay88keys/my-recipe-library/pkg/token"
)

type RedisRepository struct {
	client redis.Cmdable
}

func NewRedisRepository(client redis.Cmdable) *RedisRepository {
	return &RedisRepository{client: client}
}

func (r *RedisRepository) StoreTokenDetails(userID int64, details *token.Details) error {
	accessToken := time.Unix(details.AccessExpires, 0)
	refreshToken := time.Unix(details.RefreshExpires, 0)
	now := time.Now()

	err := r.client.Set(details.AccessUuid, strconv.Itoa(int(userID)), accessToken.Sub(now)).Err()
	if err != nil {
		return err
	}

	err = r.client.Set(details.RefreshUuid, strconv.Itoa(int(userID)), refreshToken.Sub(now)).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisRepository) RetrieveTokenDetails(details *token.AccessDetails) (int64, error) {
	foundUserID, err := r.client.Get(details.AccessUuid).Result()
	if err != nil {
		return -1, err
	}

	userID, err := strconv.ParseInt(foundUserID, 10, 64)
	if err != nil {
		return -1, err
	}

	return userID, nil
}

func (r *RedisRepository) DeleteTokenDetails(uuid string) (int64, error) {
	deleted, err := r.client.Del(uuid).Result()
	if err != nil {
		return 0, err
	}

	return deleted, nil
}
