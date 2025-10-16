package repositories

import (
	"fmt"
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

	err := r.client.Set("access:"+details.AccessUuid, strconv.Itoa(int(userID)), accessToken.Sub(now)).Err()
	if err != nil {
		return err
	}

	err = r.client.Set("refresh:"+details.RefreshUuid, strconv.Itoa(int(userID)), refreshToken.Sub(now)).Err()
	if err != nil {
		return err
	}

	// Store mapping from access UUID to refresh UUID for logout cleanup
	err = r.client.Set("access_to_refresh:"+details.AccessUuid, details.RefreshUuid, accessToken.Sub(now)).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisRepository) RetrieveTokenDetails(details *token.AccessDetails) (int64, error) {
	foundUserID, err := r.client.Get("access:" + details.AccessUuid).Result()
	if err != nil {
		return -1, err
	}

	userID, err := strconv.ParseInt(foundUserID, 10, 64)
	if err != nil {
		return -1, err
	}

	return userID, nil
}

func (r *RedisRepository) DeleteTokenDetails(uuid string) error {
	// Get the refresh UUID from the access UUID
	refreshUuid, err := r.client.Get("access_to_refresh:" + uuid).Result()
	if err != nil {
		// If we can't find the refresh UUID, just delete the access token
		_, err := r.client.Del("access:" + uuid).Result()
		if err != nil {
			fmt.Println("Error deleting access token:", err)
			return err
		}

		return nil
	}

	// Delete both access and refresh tokens atomically
	pipe := r.client.Pipeline()
	pipe.Del("access:" + uuid)
	pipe.Del("refresh:" + refreshUuid)
	pipe.Del("access_to_refresh:" + uuid) // Clean up the mapping too

	results, err := pipe.Exec()
	if err != nil {
		fmt.Println("Error deleting tokens:", err)
		return err
	}

	// Return the number of keys deleted (should be 3)
	totalDeleted := int64(0)
	for _, result := range results {
		if delResult, ok := result.(*redis.IntCmd); ok {
			totalDeleted += delResult.Val()
		}
	}

	return nil
}
