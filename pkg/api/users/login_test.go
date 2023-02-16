package users_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/iplay88keys/my-recipe-library/pkg/api"

	"github.com/iplay88keys/my-recipe-library/pkg/token"

	"github.com/iplay88keys/my-recipe-library/pkg/api/users"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("login", func() {
	It("logs a user in", func() {
		verify := func(login, password string) (bool, int64, error) {
			return true, 0, nil
		}

		createToken := func(userid int64) (*token.Details, error) {
			return &token.Details{
				AccessToken:  "access token",
				RefreshToken: "refresh token",
			}, nil
		}

		storeToken := func(userid int64, details *token.Details) error {
			return nil
		}

		body := []byte(`{
            "login": "username",
            "password": "Pa3$12345"
        }`)

		req, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Login(verify, createToken, storeToken).Handle(&api.Request{
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

	It("returns unauthorized if the login fails due to bad credentials", func() {
		verify := func(login, password string) (bool, int64, error) {
			return false, 0, nil
		}

		createToken := func(userid int64) (*token.Details, error) {
			return &token.Details{}, nil
		}

		storeToken := func(userid int64, details *token.Details) error {
			return nil
		}

		body := []byte(`{
            "login": "username",
            "password": "bad-password"
        }`)

		req, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Login(verify, createToken, storeToken).Handle(&api.Request{
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

	It("returns a bad request if the body is empty", func() {
		verify := func(login, password string) (bool, int64, error) {
			return false, 0, nil
		}

		createToken := func(userid int64) (*token.Details, error) {
			return &token.Details{}, nil
		}

		storeToken := func(userid int64, details *token.Details) error {
			return nil
		}

		req, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer([]byte("")))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Login(verify, createToken, storeToken).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
	})

	It("returns an error if the login check fails", func() {
		verify := func(login, password string) (bool, int64, error) {
			return false, 0, errors.New("some error")
		}

		createToken := func(userid int64) (*token.Details, error) {
			return &token.Details{}, nil
		}

		storeToken := func(userid int64, details *token.Details) error {
			return nil
		}

		body := []byte(`{
            "login": "username",
            "password": "bad-password"
        }`)

		req, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Login(verify, createToken, storeToken).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
	})

	It("returns an error if the token creation fails", func() {
		verify := func(login, password string) (bool, int64, error) {
			return true, 0, nil
		}

		createToken := func(userid int64) (*token.Details, error) {
			return nil, errors.New("some error")
		}

		storeToken := func(userid int64, details *token.Details) error {
			return nil
		}

		body := []byte(`{
            "login": "username",
            "password": "bad-password"
        }`)

		req, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Login(verify, createToken, storeToken).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
	})

	It("returns an error if the token storing fails", func() {
		verify := func(login, password string) (bool, int64, error) {
			return true, 0, nil
		}

		createToken := func(userid int64) (*token.Details, error) {
			return &token.Details{}, nil
		}

		storeToken := func(userid int64, details *token.Details) error {
			return errors.New("some error")
		}

		body := []byte(`{
            "login": "username",
            "password": "bad-password"
        }`)

		req, err := http.NewRequest(http.MethodPost, "/users/login", bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Login(verify, createToken, storeToken).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
	})
})
