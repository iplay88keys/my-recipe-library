package services

import (
	"context"
	"errors"
	"net/http"

	"github.com/iplay88keys/my-recipe-library/pkg/token"
)

type TokenServiceInterface interface {
	CreateToken(userID int64) (*token.Details, error)
	ValidateToken(r *http.Request) (*token.AccessDetails, error)
}

type RedisRepositoryInterface interface {
	StoreTokenDetails(userID int64, details *token.Details) error
	DeleteTokenDetails(uuid string) error
}

type UsersRepositoryInterface interface {
	ExistsByUsername(username string) (bool, error)
	ExistsByEmail(email string) (bool, error)
	Insert(username, email, password string) (int64, error)
	Verify(login, password string) (bool, int64, error)
}

type UserService struct {
	usersRepo    UsersRepositoryInterface
	redisRepo    RedisRepositoryInterface
	tokenService TokenServiceInterface
}

func NewUserService(
	usersRepo UsersRepositoryInterface,
	redisRepo RedisRepositoryInterface,
	tokenService TokenServiceInterface,
) *UserService {
	return &UserService{
		usersRepo:    usersRepo,
		redisRepo:    redisRepo,
		tokenService: tokenService,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, username, email, password string) error {
	exists, err := s.usersRepo.ExistsByUsername(username)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("username already exists")
	}

	exists, err = s.usersRepo.ExistsByEmail(email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("email already exists")
	}

	_, err = s.usersRepo.Insert(username, email, password)
	return err
}

func (s *UserService) Verify(login, password string) (bool, int64, error) {
	return s.usersRepo.Verify(login, password)
}

func (s *UserService) CreateToken(userID int64) (*token.Details, error) {
	return s.tokenService.CreateToken(userID)
}

func (s *UserService) StoreTokenDetails(userID int64, details *token.Details) error {
	return s.redisRepo.StoreTokenDetails(userID, details)
}

func (s *UserService) ValidateToken(r *http.Request) (*token.AccessDetails, error) {
	return s.tokenService.ValidateToken(r)
}

func (s *UserService) DeleteTokenDetails(uuid string) error {
	return s.redisRepo.DeleteTokenDetails(uuid)
}
