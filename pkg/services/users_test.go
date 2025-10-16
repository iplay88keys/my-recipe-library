package services_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/iplay88keys/my-recipe-library/pkg/services"
	"github.com/iplay88keys/my-recipe-library/pkg/token"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("UserService", func() {
	var (
		userService      *services.UserService
		mockUsersRepo    *MockUsersRepository
		mockTokenService *MockTokenService
		mockRedisRepo    *MockRedisRepository
		ctx              context.Context
		userID           int64
	)

	BeforeEach(func() {
		mockUsersRepo = &MockUsersRepository{}
		mockTokenService = &MockTokenService{}
		mockRedisRepo = &MockRedisRepository{}
		userService = services.NewUserService(mockUsersRepo, mockRedisRepo, mockTokenService)

		ctx = context.Background()
		userID = 1
	})

	Describe("RegisterUser", func() {
		var username, email, password string

		BeforeEach(func() {
			username = "testuser"
			email = "test@example.com"
			password = "password123"
		})

		Context("when registration is successful", func() {
			It("registers a new user", func() {
				mockUsersRepo.ExistsByUsernameFunc = func(username string) (bool, error) {
					Expect(username).To(Equal(username))
					return false, nil
				}

				mockUsersRepo.ExistsByEmailFunc = func(email string) (bool, error) {
					Expect(email).To(Equal(email))
					return false, nil
				}

				mockUsersRepo.InsertFunc = func(username, email, password string) (int64, error) {
					Expect(username).To(Equal(username))
					Expect(email).To(Equal(email))
					Expect(password).ToNot(BeEmpty())
					return 1, nil
				}

				err := userService.RegisterUser(ctx, username, email, password)

				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when username already exists", func() {
			It("returns an error", func() {
				mockUsersRepo.ExistsByUsernameFunc = func(username string) (bool, error) {
					Expect(username).To(Equal(username))
					return true, nil
				}

				err := userService.RegisterUser(ctx, username, email, password)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("username already exists"))
			})
		})

		Context("when email already exists", func() {
			It("returns an error", func() {
				mockUsersRepo.ExistsByUsernameFunc = func(username string) (bool, error) {
					return false, nil
				}

				mockUsersRepo.ExistsByEmailFunc = func(email string) (bool, error) {
					Expect(email).To(Equal(email))
					return true, nil
				}

				err := userService.RegisterUser(ctx, username, email, password)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("email already exists"))
			})
		})

		Context("when username check fails", func() {
			It("returns an error", func() {
				mockUsersRepo.ExistsByUsernameFunc = func(username string) (bool, error) {
					return false, errors.New("failed to query for user by username")
				}

				err := userService.RegisterUser(ctx, username, email, password)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to query for user by username"))
			})
		})

		Context("when email check fails", func() {
			It("returns an error", func() {
				mockUsersRepo.ExistsByUsernameFunc = func(username string) (bool, error) {
					return false, nil
				}

				mockUsersRepo.ExistsByEmailFunc = func(email string) (bool, error) {
					return false, errors.New("failed to query for user by email")
				}

				err := userService.RegisterUser(ctx, username, email, password)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to query for user by email"))
			})
		})

		Context("when user insert fails", func() {
			It("returns an error", func() {
				mockUsersRepo.ExistsByUsernameFunc = func(username string) (bool, error) {
					return false, nil
				}

				mockUsersRepo.ExistsByEmailFunc = func(email string) (bool, error) {
					return false, nil
				}

				mockUsersRepo.InsertFunc = func(username, email, password string) (int64, error) {
					return 0, errors.New("user could not be added")
				}

				err := userService.RegisterUser(ctx, username, email, password)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("user could not be added"))
			})
		})
	})

	Describe("Verify", func() {
		var login, password string

		BeforeEach(func() {
			login = "testuser"
			password = "password123"
		})

		Context("when credentials are valid", func() {
			It("returns true and user ID", func() {
				mockUsersRepo.VerifyFunc = func(login, password string) (bool, int64, error) {
					Expect(login).To(Equal(login))
					Expect(password).To(Equal(password))
					return true, userID, nil
				}

				valid, returnedUserID, err := userService.Verify(login, password)

				Expect(err).ToNot(HaveOccurred())
				Expect(valid).To(BeTrue())
				Expect(returnedUserID).To(Equal(userID))
			})
		})

		Context("when credentials are invalid", func() {
			It("returns false and zero user ID", func() {
				mockUsersRepo.VerifyFunc = func(login, password string) (bool, int64, error) {
					return false, -1, nil
				}

				valid, returnedUserID, err := userService.Verify(login, password)

				Expect(err).ToNot(HaveOccurred())
				Expect(valid).To(BeFalse())
				Expect(returnedUserID).To(Equal(int64(-1)))
			})
		})

		Context("when database query fails", func() {
			It("returns an error", func() {
				mockUsersRepo.VerifyFunc = func(login, password string) (bool, int64, error) {
					return false, -1, errors.New("database error")
				}

				valid, returnedUserID, err := userService.Verify(login, password)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("database error"))
				Expect(valid).To(BeFalse())
				Expect(returnedUserID).To(Equal(int64(-1)))
			})
		})
	})

	Describe("CreateToken", func() {
		Context("when token creation is successful", func() {
			It("returns token details", func() {
				expectedDetails := &token.Details{
					AccessToken:    "access-token",
					RefreshToken:   "refresh-token",
					AccessUuid:     "access-uuid",
					RefreshUuid:    "refresh-uuid",
					AccessExpires:  1234567890,
					RefreshExpires: 1234567890,
				}

				mockTokenService.CreateTokenFunc = func(userID int64) (*token.Details, error) {
					Expect(userID).To(Equal(userID))
					return expectedDetails, nil
				}

				result, err := userService.CreateToken(userID)

				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(Equal(expectedDetails))
			})
		})

		Context("when token creation fails", func() {
			It("returns an error", func() {
				mockTokenService.CreateTokenFunc = func(userID int64) (*token.Details, error) {
					return nil, errors.New("token creation failed")
				}

				result, err := userService.CreateToken(userID)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("token creation failed"))
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("StoreTokenDetails", func() {
		Context("when storing token details is successful", func() {
			It("stores token details in Redis", func() {
				details := &token.Details{
					AccessToken:    "access-token",
					RefreshToken:   "refresh-token",
					AccessUuid:     "access-uuid",
					RefreshUuid:    "refresh-uuid",
					AccessExpires:  1234567890,
					RefreshExpires: 1234567890,
				}

				mockRedisRepo.StoreTokenDetailsFunc = func(userID int64, details *token.Details) error {
					Expect(userID).To(Equal(userID))
					Expect(details).To(Equal(details))
					return nil
				}

				err := userService.StoreTokenDetails(userID, details)

				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when storing token details fails", func() {
			It("returns an error", func() {
				details := &token.Details{}

				mockRedisRepo.StoreTokenDetailsFunc = func(userID int64, details *token.Details) error {
					return errors.New("redis error")
				}

				err := userService.StoreTokenDetails(userID, details)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("redis error"))
			})
		})
	})

	Describe("ValidateToken", func() {
		Context("when token validation is successful", func() {
			It("returns access details", func() {
				expectedDetails := &token.AccessDetails{
					AccessUuid: "access-uuid",
					UserId:     userID,
				}

				req := httptest.NewRequest("GET", "/test", nil)

				mockTokenService.ValidateTokenFunc = func(r *http.Request) (*token.AccessDetails, error) {
					Expect(r).To(Equal(req))
					return expectedDetails, nil
				}

				result, err := userService.ValidateToken(req)

				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(Equal(expectedDetails))
			})
		})

		Context("when token validation fails", func() {
			It("returns an error", func() {
				req := httptest.NewRequest("GET", "/test", nil)

				mockTokenService.ValidateTokenFunc = func(r *http.Request) (*token.AccessDetails, error) {
					return nil, errors.New("invalid token")
				}

				result, err := userService.ValidateToken(req)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("invalid token"))
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("DeleteTokenDetails", func() {
		Context("when deleting token details is successful", func() {
			It("deletes token details from Redis", func() {
				uuid := "access-uuid"

				mockRedisRepo.DeleteTokenDetailsFunc = func(uuid string) error {
					Expect(uuid).To(Equal(uuid))
					return nil
				}

				err := userService.DeleteTokenDetails(uuid)

				Expect(err).ToNot(HaveOccurred())
			})
		})

		Context("when deleting token details fails", func() {
			It("returns an error", func() {
				uuid := "access-uuid"

				mockRedisRepo.DeleteTokenDetailsFunc = func(uuid string) error {
					return errors.New("redis error")
				}

				err := userService.DeleteTokenDetails(uuid)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("redis error"))
			})
		})
	})
})

type MockUsersRepository struct {
	ExistsByUsernameFunc func(username string) (bool, error)
	ExistsByEmailFunc    func(email string) (bool, error)
	InsertFunc           func(username, email, password string) (int64, error)
	VerifyFunc           func(login, password string) (bool, int64, error)
}

func (m *MockUsersRepository) ExistsByUsername(username string) (bool, error) {
	if m.ExistsByUsernameFunc != nil {
		return m.ExistsByUsernameFunc(username)
	}
	return false, nil
}

func (m *MockUsersRepository) ExistsByEmail(email string) (bool, error) {
	if m.ExistsByEmailFunc != nil {
		return m.ExistsByEmailFunc(email)
	}
	return false, nil
}

func (m *MockUsersRepository) Insert(username, email, password string) (int64, error) {
	if m.InsertFunc != nil {
		return m.InsertFunc(username, email, password)
	}
	return 1, nil
}

func (m *MockUsersRepository) Verify(login, password string) (bool, int64, error) {
	if m.VerifyFunc != nil {
		return m.VerifyFunc(login, password)
	}
	return false, -1, nil
}

type MockTokenService struct {
	CreateTokenFunc   func(userID int64) (*token.Details, error)
	ValidateTokenFunc func(r *http.Request) (*token.AccessDetails, error)
}

func (m *MockTokenService) CreateToken(userID int64) (*token.Details, error) {
	if m.CreateTokenFunc != nil {
		return m.CreateTokenFunc(userID)
	}
	return &token.Details{
		AccessToken:    "mock-access-token",
		RefreshToken:   "mock-refresh-token",
		AccessUuid:     "mock-access-uuid",
		RefreshUuid:    "mock-refresh-uuid",
		AccessExpires:  1234567890,
		RefreshExpires: 1234567890,
	}, nil
}

func (m *MockTokenService) ValidateToken(r *http.Request) (*token.AccessDetails, error) {
	if m.ValidateTokenFunc != nil {
		return m.ValidateTokenFunc(r)
	}
	return &token.AccessDetails{
		AccessUuid: "mock-access-uuid",
		UserId:     1,
	}, nil
}

type MockRedisRepository struct {
	StoreTokenDetailsFunc  func(userID int64, details *token.Details) error
	DeleteTokenDetailsFunc func(uuid string) error
}

func (m *MockRedisRepository) StoreTokenDetails(userID int64, details *token.Details) error {
	if m.StoreTokenDetailsFunc != nil {
		return m.StoreTokenDetailsFunc(userID, details)
	}
	return nil
}

func (m *MockRedisRepository) DeleteTokenDetails(uuid string) error {
	if m.DeleteTokenDetailsFunc != nil {
		return m.DeleteTokenDetailsFunc(uuid)
	}
	return nil
}
