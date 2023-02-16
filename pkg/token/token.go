package token

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
)

type UserClaims struct {
	UserID      int64
	AccessUUID  string
	RefreshUUID string
	jwt.StandardClaims
}

func (u UserClaims) Valid() error {
	return nil
}

type Details struct {
	AccessToken    string
	RefreshToken   string
	AccessUuid     string
	RefreshUuid    string
	AccessExpires  int64
	RefreshExpires int64
}

type AccessDetails struct {
	AccessUuid string
	UserId     int64
}

type Service struct {
	accessSecret  string
	refreshSecret string
}

func NewService(accessSecret, refreshSecret string) *Service {
	return &Service{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
	}
}

func (s *Service) CreateToken(userID int64) (*Details, error) {
	var err error
	details := &Details{}
	details.AccessExpires = time.Now().Add(time.Minute * 15).Unix()
	details.AccessUuid = uuid.NewV4().String()

	details.RefreshExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	details.RefreshUuid = uuid.NewV4().String()

	accessClaims := UserClaims{
		UserID:     userID,
		AccessUUID: details.AccessUuid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: details.AccessExpires,
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	details.AccessToken, err = accessToken.SignedString([]byte(s.accessSecret))
	if err != nil {
		return nil, err
	}

	refreshClaims := UserClaims{
		UserID:      userID,
		RefreshUUID: details.RefreshToken,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: details.RefreshExpires,
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	details.RefreshToken, err = refreshToken.SignedString([]byte(s.accessSecret))
	if err != nil {
		return nil, err
	}

	return details, nil
}

func (s *Service) ValidateToken(r *http.Request) (*AccessDetails, error) {
	auth := strings.Split(r.Header.Get("Authorization"), " ")
	if len(auth) != 2 {
		return nil, errors.New("authorization header missing token")
	}

	token, err := jwt.ParseWithClaims(auth[1], &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.accessSecret), nil
	})
	if err != nil {
		return nil, err
	}

	var claims *UserClaims
	if tokenClaims, ok := token.Claims.(*UserClaims); !ok || !token.Valid {
		return nil, errors.New("invalid token")
	} else {
		claims = tokenClaims
	}

	return &AccessDetails{
		AccessUuid: claims.AccessUUID,
		UserId:     claims.UserID,
	}, nil
}
