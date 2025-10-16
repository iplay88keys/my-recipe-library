package users_test

import (
	"bytes"
	"context"
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
		fakeService := &mockUserRegistrar{
			registerUser: func(ctx context.Context, username, email, password string) error {
				return nil
			},
		}

		body := []byte(`{
            "username": "username",
            "email":    "email@example.com",
            "password": "Pa3$12345"
        }`)

		req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Register(fakeService).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusOK))
	})

	It("returns validation info", func() {
		fakeService := &mockUserRegistrar{
			registerUser: func(ctx context.Context, username, email, password string) error {
				return nil
			},
		}

		body := []byte(`{
            "username": "",
            "email":    "",
            "password": ""
        }`)

		req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Register(fakeService).Handle(&api.Request{
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
		fakeService := &mockUserRegistrar{
			registerUser: func(ctx context.Context, username, email, password string) error {
				return errors.New("username already exists")
			},
		}

		body := []byte(`{
            "username": "username",
            "email":    "email@example.com",
            "password": "Pa3$12345"
        }`)

		req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Register(fakeService).Handle(&api.Request{
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
		fakeService := &mockUserRegistrar{
			registerUser: func(ctx context.Context, username, email, password string) error {
				return errors.New("email already exists")
			},
		}

		body := []byte(`{
            "username": "username",
            "email":    "email@example.com",
            "password": "Pa3$12345"
        }`)

		req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Register(fakeService).Handle(&api.Request{
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
		fakeService := &mockUserRegistrar{
			registerUser: func(ctx context.Context, username, email, password string) error {
				return errors.New("some error")
			},
		}

		req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer([]byte("")))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Register(fakeService).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
	})

	It("returns an error if the username check fails", func() {
		fakeService := &mockUserRegistrar{
			registerUser: func(ctx context.Context, username, email, password string) error {
				return errors.New("database error")
			},
		}

		body := []byte(`{
            "username": "username",
            "email":    "email@example.com",
            "password": "Pa3$12345"
        }`)

		req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Register(fakeService).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
	})

	It("returns an error if the email check fails", func() {
		fakeService := &mockUserRegistrar{
			registerUser: func(ctx context.Context, username, email, password string) error {
				return errors.New("database error")
			},
		}

		body := []byte(`{
            "username": "username",
            "email":    "email@example.com",
            "password": "Pa3$12345"
        }`)

		req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Register(fakeService).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
	})

	It("returns an error if the user insert fails", func() {
		fakeService := &mockUserRegistrar{
			registerUser: func(ctx context.Context, username, email, password string) error {
				return errors.New("insert failed")
			},
		}

		body := []byte(`{
            "username": "username",
            "email":    "email@example.com",
            "password": "Pa3$12345"
        }`)

		req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Register(fakeService).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
	})

	It("returns validation info for invalid password", func() {
		fakeService := &mockUserRegistrar{
			registerUser: func(ctx context.Context, username, email, password string) error {
				return nil
			},
		}

		body := []byte(`{
            "username": "username",
            "email":    "email@example.com",
            "password": "weak"
        }`)

		req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Register(fakeService).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

		respBody, err := json.Marshal(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(respBody).To(MatchJSON(`{
            "errors": {
                "password": "Uppercase letter missing, Numeric character missing, Special character missing, Must be between 6 and 64 characters long"
            }
        }`))
	})

	It("returns validation info for invalid email", func() {
		fakeService := &mockUserRegistrar{
			registerUser: func(ctx context.Context, username, email, password string) error {
				return nil
			},
		}

		body := []byte(`{
            "username": "username",
            "email":    "invalid-email",
            "password": "Pa3$12345"
        }`)

		req, err := http.NewRequest(http.MethodPost, "/users/register", bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		resp := users.Register(fakeService).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

		respBody, err := json.Marshal(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(respBody).To(MatchJSON(`{
            "errors": {
                "email": "Invalid email address"
            }
        }`))
	})
})

type mockUserRegistrar struct {
	registerUser func(ctx context.Context, username, email, password string) error
}

func (m *mockUserRegistrar) RegisterUser(ctx context.Context, username, email, password string) error {
	return m.registerUser(ctx, username, email, password)
}
