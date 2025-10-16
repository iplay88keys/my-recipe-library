package repositories_test

import (
	"errors"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/elliotchance/redismock"
	"github.com/go-redis/redis"
	"github.com/stretchr/testify/mock"

	"github.com/iplay88keys/my-recipe-library/pkg/repositories"
	"github.com/iplay88keys/my-recipe-library/pkg/token"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Redis Repository", func() {
	var (
		redisClient *redismock.ClientMock
	)

	BeforeEach(func() {
		mr, err := miniredis.Run()
		if err != nil {
			panic(err)
		}

		client := redis.NewClient(&redis.Options{
			Addr: mr.Addr(),
		})

		redisClient = redismock.NewNiceMock(client)
	})

	Describe("StoreTokenDetails", func() {
		It("stores token details", func() {
			redisRepo := repositories.NewRedisRepository(redisClient)

			accessExpiration := time.Now().Add(time.Minute).Unix()
			refreshExpiration := time.Now().Add(time.Hour).Unix()

			redisClient.On("Set", "access:access uuid", "10", mock.AnythingOfType("time.Duration")).
				Return(redis.NewStatusResult("", nil))
			redisClient.On("Set", "refresh:refresh uuid", "10", mock.AnythingOfType("time.Duration")).
				Return(redis.NewStatusResult("", nil))
			redisClient.On("Set", "access_to_refresh:access uuid", "refresh uuid", mock.AnythingOfType("time.Duration")).
				Return(redis.NewStatusResult("", nil))

			details := &token.Details{
				AccessToken:    "access token",
				RefreshToken:   "refresh token",
				AccessUuid:     "access uuid",
				RefreshUuid:    "refresh uuid",
				AccessExpires:  accessExpiration,
				RefreshExpires: refreshExpiration,
			}

			err := redisRepo.StoreTokenDetails(10, details)
			Expect(err).ToNot(HaveOccurred())

			redisClient.AssertNumberOfCalls(GinkgoT(), "Set", 3)
		})

		It("returns an error if the access token set fails", func() {
			redisRepo := repositories.NewRedisRepository(redisClient)

			accessExpiration := time.Now().Add(time.Minute).Unix()
			refreshExpiration := time.Now().Add(time.Hour).Unix()

			redisClient.On("Set", "access:access uuid", "10", mock.AnythingOfType("time.Duration")).
				Return(redis.NewStatusResult("", errors.New("some redis error")))
			redisClient.On("Set", "refresh:refresh uuid", "10", mock.AnythingOfType("time.Duration")).
				Return(redis.NewStatusResult("", nil))
			redisClient.On("Set", "access_to_refresh:access uuid", "refresh uuid", mock.AnythingOfType("time.Duration")).
				Return(redis.NewStatusResult("", nil))

			details := &token.Details{
				AccessToken:    "access token",
				RefreshToken:   "refresh token",
				AccessUuid:     "access uuid",
				RefreshUuid:    "refresh uuid",
				AccessExpires:  accessExpiration,
				RefreshExpires: refreshExpiration,
			}

			err := redisRepo.StoreTokenDetails(10, details)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("some redis error"))
		})

		It("returns an error if the refresh token set fails", func() {
			redisRepo := repositories.NewRedisRepository(redisClient)

			accessExpiration := time.Now().Add(time.Minute).Unix()
			refreshExpiration := time.Now().Add(time.Hour).Unix()

			redisClient.On("Set", "access:access uuid", "10", mock.AnythingOfType("time.Duration")).
				Return(redis.NewStatusResult("", nil))
			redisClient.On("Set", "refresh:refresh uuid", "10", mock.AnythingOfType("time.Duration")).
				Return(redis.NewStatusResult("", errors.New("some redis error")))
			redisClient.On("Set", "access_to_refresh:access uuid", "refresh uuid", mock.AnythingOfType("time.Duration")).
				Return(redis.NewStatusResult("", nil))

			details := &token.Details{
				AccessToken:    "access token",
				RefreshToken:   "refresh token",
				AccessUuid:     "access uuid",
				RefreshUuid:    "refresh uuid",
				AccessExpires:  accessExpiration,
				RefreshExpires: refreshExpiration,
			}

			err := redisRepo.StoreTokenDetails(10, details)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("some redis error"))
		})
	})

	Describe("RetrieveTokenDetails", func() {
		It("retrieves access token details", func() {
			redisRepo := repositories.NewRedisRepository(redisClient)

			redisClient.On("Get", "access:access uuid").
				Return(redis.NewStringResult("10", nil))

			details := &token.AccessDetails{
				AccessUuid: "access uuid",
				UserId:     10,
			}

			userID, err := redisRepo.RetrieveTokenDetails(details)
			Expect(err).ToNot(HaveOccurred())
			Expect(userID).To(Equal(details.UserId))

			redisClient.AssertNumberOfCalls(GinkgoT(), "Get", 1)
		})

		It("returns an error if the access uuid is not found", func() {
			redisRepo := repositories.NewRedisRepository(redisClient)

			redisClient.On("Get", "access:access uuid").
				Return(redis.NewStringResult("10", errors.New("some redis error")))

			details := &token.AccessDetails{
				AccessUuid: "access uuid",
				UserId:     10,
			}

			_, err := redisRepo.RetrieveTokenDetails(details)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("some redis error"))
		})

		It("returns an error if the user id returned cannot be converted to an int", func() {
			redisRepo := repositories.NewRedisRepository(redisClient)

			redisClient.On("Get", "access:access uuid").
				Return(redis.NewStringResult("incorrect response", nil))

			details := &token.AccessDetails{
				AccessUuid: "access uuid",
				UserId:     10,
			}

			_, err := redisRepo.RetrieveTokenDetails(details)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("invalid syntax"))
		})
	})

	Describe("DeleteTokenDetails", func() {
		It("deletes access token details when refresh UUID not found", func() {
			redisRepo := repositories.NewRedisRepository(redisClient)

			redisClient.On("Get", "access_to_refresh:access uuid").
				Return(redis.NewStringResult("", errors.New("not found")))
			redisClient.On("Del", []string{"access:access uuid"}).
				Return(redis.NewIntResult(1, nil))

			err := redisRepo.DeleteTokenDetails("access uuid")
			Expect(err).ToNot(HaveOccurred())

			redisClient.AssertNumberOfCalls(GinkgoT(), "Get", 1)
			redisClient.AssertNumberOfCalls(GinkgoT(), "Del", 1)
		})

		It("returns an error if the delete fails", func() {
			redisRepo := repositories.NewRedisRepository(redisClient)

			redisClient.On("Get", "access_to_refresh:access uuid").
				Return(redis.NewStringResult("", errors.New("not found")))
			redisClient.On("Del", []string{"access:access uuid"}).
				Return(redis.NewIntResult(1, errors.New("some redis error")))

			err := redisRepo.DeleteTokenDetails("access uuid")
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("some redis error"))
		})
	})
})
