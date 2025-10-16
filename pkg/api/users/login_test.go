package users_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/iplay88keys/my-recipe-library/pkg/api"
	"github.com/iplay88keys/my-recipe-library/pkg/api/users"
	"github.com/iplay88keys/my-recipe-library/pkg/token"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("login", func() {
	It("logs a user in", func() {
		fakeLoginService := &mockLoginService{
			verify: func(login, password string) (bool, int64, error) {
				return true, 1, nil
			},
			createToken: func(userID int64) (*token.Details, error) {
				return &token.Details{
					AccessToken:  "access token",
					RefreshToken: "refresh token",
				}, nil
			},
			storeTokenDetails: func(userID int64, details *token.Details) error {
				return nil
			},
		}

		body := []byte(`{
            "login": "username",
            "password": "Pa3$12345"
        }`)

		req, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Login(fakeLoginService).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		respBody, err := json.Marshal(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(respBody).To(MatchJSON(`{
            "access_token": "access token",
            "refresh_token": "refresh token"
        }`))
	})

	It("returns validation info", func() {
		fakeLoginService := &mockLoginService{
			verify: func(login, password string) (bool, int64, error) {
				return false, 0, nil
			},
			createToken: func(userID int64) (*token.Details, error) {
				return &token.Details{
					AccessToken:  "access token",
					RefreshToken: "refresh token",
				}, nil
			},
			storeTokenDetails: func(userID int64, details *token.Details) error {
				return nil
			},
		}

		body := []byte(`{
            "login": "",
            "password": ""
        }`)

		req, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Login(fakeLoginService).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
	})

	It("returns unauthorized for invalid credentials", func() {
		fakeLoginService := &mockLoginService{
			verify: func(login, password string) (bool, int64, error) {
				return false, 0, nil
			},
			createToken: func(userID int64) (*token.Details, error) {
				return &token.Details{
					AccessToken:  "access token",
					RefreshToken: "refresh token",
				}, nil
			},
			storeTokenDetails: func(userID int64, details *token.Details) error {
				return nil
			},
		}

		body := []byte(`{
            "login": "username",
            "password": "wrongpassword"
        }`)

		req, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Login(fakeLoginService).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))

		respBody, err := json.Marshal(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(respBody).To(MatchJSON(`{
            "errors": {
                "alert": "Invalid login credentials"
            }
        }`))
	})

	It("returns internal server error if verification fails", func() {
		fakeLoginService := &mockLoginService{
			verify: func(login, password string) (bool, int64, error) {
				return false, 0, errors.New("verification failed")
			},
			createToken: func(userID int64) (*token.Details, error) {
				return &token.Details{
					AccessToken:  "access token",
					RefreshToken: "refresh token",
				}, nil
			},
			storeTokenDetails: func(userID int64, details *token.Details) error {
				return nil
			},
		}

		body := []byte(`{
            "login": "username",
            "password": "Pa3$12345"
        }`)

		req, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Login(fakeLoginService).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
	})

	It("returns internal server error if token creation fails", func() {
		fakeLoginService := &mockLoginService{
			verify: func(login, password string) (bool, int64, error) {
				return true, 1, nil
			},
			createToken: func(userID int64) (*token.Details, error) {
				return nil, errors.New("token creation failed")
			},
			storeTokenDetails: func(userID int64, details *token.Details) error {
				return nil
			},
		}

		body := []byte(`{
            "login": "username",
            "password": "Pa3$12345"
        }`)

		req, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Login(fakeLoginService).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
	})

	It("returns internal server error if token storage fails", func() {
		fakeLoginService := &mockLoginService{
			verify: func(login, password string) (bool, int64, error) {
				return true, 1, nil
			},
			createToken: func(userID int64) (*token.Details, error) {
				return &token.Details{
					AccessToken:  "access token",
					RefreshToken: "refresh token",
				}, nil
			},
			storeTokenDetails: func(userID int64, details *token.Details) error {
				return errors.New("token storage failed")
			},
		}

		body := []byte(`{
            "login": "username",
            "password": "Pa3$12345"
        }`)

		req, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Login(fakeLoginService).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
	})
})

type mockLoginService struct {
	verify            func(login, password string) (bool, int64, error)
	createToken       func(userID int64) (*token.Details, error)
	storeTokenDetails func(userID int64, details *token.Details) error
}

func (m *mockLoginService) Verify(login, password string) (bool, int64, error) {
	return m.verify(login, password)
}

func (m *mockLoginService) CreateToken(userID int64) (*token.Details, error) {
	return m.createToken(userID)
}

func (m *mockLoginService) StoreTokenDetails(userID int64, details *token.Details) error {
	return m.storeTokenDetails(userID, details)
}
