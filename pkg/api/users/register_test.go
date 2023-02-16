package users_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/iplay88keys/my-recipe-library/pkg/api"

	"github.com/iplay88keys/my-recipe-library/pkg/api/users"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("register", func() {
	It("creates a user", func() {
		existsByUsername := func(username string) (bool, error) {
			return false, nil
		}

		existsByEmail := func(email string) (bool, error) {
			return false, nil
		}

		insertUser := func(username, email, password string) (int64, error) {
			return 1, nil
		}

		body := []byte(`{
            "username": "username",
            "email":    "email@example.com",
            "password": "Pa3$12345"
        }`)

		req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Register(existsByUsername, existsByEmail, insertUser).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusOK))
	})

	It("returns validation info", func() {
		existsByUsername := func(username string) (bool, error) {
			return false, nil
		}

		existsByEmail := func(email string) (bool, error) {
			return false, nil
		}

		insertUser := func(username, email, password string) (int64, error) {
			return 1, nil
		}

		body := []byte(`{
            "username": "",
            "email":    "",
            "password": ""
        }`)

		req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Register(existsByUsername, existsByEmail, insertUser).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

		respBody, err := json.Marshal(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(respBody).To(MatchJSON(`{
            "errors": {
                "email": "Required",
                "password": "Required",
                "username": "Required"
            }
        }`))
	})

	It("returns info if the username already exists", func() {
		existsByUsername := func(username string) (bool, error) {
			return true, nil
		}

		existsByEmail := func(email string) (bool, error) {
			return false, nil
		}

		insertUser := func(username, email, password string) (int64, error) {
			return 1, nil
		}

		body := []byte(`{
            "username": "username",
            "email":    "email@example.com",
            "password": "Pa3$12345"
        }`)

		req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Register(existsByUsername, existsByEmail, insertUser).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

		respBody, err := json.Marshal(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(respBody).To(MatchJSON(`{
            "errors": {
                "username": "Username already in use"
            }
        }`))
	})

	It("returns info if the email already exists", func() {
		existsByUsername := func(username string) (bool, error) {
			return false, nil
		}

		existsByEmail := func(email string) (bool, error) {
			return true, nil
		}

		insertUser := func(username, email, password string) (int64, error) {
			return 1, nil
		}

		body := []byte(`{
            "username": "username",
            "email":    "email@example.com",
            "password": "Pa3$12345"
        }`)

		req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Register(existsByUsername, existsByEmail, insertUser).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

		respBody, err := json.Marshal(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(respBody).To(MatchJSON(`{
            "errors": {
                "email": "Email already in use"
            }
        }`))
	})

	It("returns bad request if the body is empty", func() {
		existsByUsername := func(username string) (bool, error) {
			return false, nil
		}

		existsByEmail := func(email string) (bool, error) {
			return false, nil
		}

		insertUser := func(username, email, password string) (int64, error) {
			return -1, errors.New("some error")
		}

		req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer([]byte("")))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Register(existsByUsername, existsByEmail, insertUser).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
	})

	It("returns an error if the username check fails", func() {
		existsByUsername := func(username string) (bool, error) {
			return false, errors.New("error username")
		}

		existsByEmail := func(email string) (bool, error) {
			return false, nil
		}

		insertUser := func(username, email, password string) (int64, error) {
			return 1, nil
		}

		body := []byte(`{
            "username": "username",
            "email":    "email@example.com",
            "password": "Pa3$12345"
        }`)

		req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Register(existsByUsername, existsByEmail, insertUser).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
	})

	It("returns an error if the email check fails", func() {
		existsByUsername := func(username string) (bool, error) {
			return false, nil
		}

		existsByEmail := func(email string) (bool, error) {
			return false, errors.New("error email")
		}

		insertUser := func(username, email, password string) (int64, error) {
			return 1, nil
		}

		body := []byte(`{
            "username": "username",
            "email":    "email@example.com",
            "password": "Pa3$12345"
        }`)

		req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Register(existsByUsername, existsByEmail, insertUser).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
	})

	It("returns an error if the user insert fails", func() {
		existsByUsername := func(username string) (bool, error) {
			return false, nil
		}

		existsByEmail := func(email string) (bool, error) {
			return false, nil
		}

		insertUser := func(username, email, password string) (int64, error) {
			return -1, errors.New("some error")
		}

		body := []byte(`{
            "username": "username",
            "email":    "email@example.com",
            "password": "Pa3$12345"
        }`)

		req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Register(existsByUsername, existsByEmail, insertUser).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
	})
})
