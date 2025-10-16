package api_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/iplay88keys/my-recipe-library/pkg/token"

	"github.com/iplay88keys/my-recipe-library/pkg/api"
	"github.com/iplay88keys/my-recipe-library/pkg/helpers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("API", func() {
	var (
		server *api.API
		port   string
	)

	BeforeEach(func() {
		var err error

		port, err = helpers.GetRandomPort()
		Expect(err).ToNot(HaveOccurred())

		validateToken := func(r *http.Request) (*token.AccessDetails, error) {
			if r.Header.Get("Authorization") == "bearer token" {
				return &token.AccessDetails{
					AccessUuid: "some-uuid",
					UserId:     10,
				}, nil
			} else {
				return nil, errors.New("auth error")
			}
		}

		retrieveTokenDetails := func(details *token.AccessDetails) (int64, error) {
			return 10, nil
		}

		server = api.New(
			&mockTokenValidator{validateToken: validateToken},
			&mockAccessDetailsRetriever{retrieveTokenDetails: retrieveTokenDetails},
			&api.Config{
				Port:      port,
				StaticDir: "fixtures",
				Endpoints: []*api.Endpoint{{
					Path:   "test-unauthenticated-endpoint",
					Method: http.MethodGet,
					Handle: func(r *api.Request) *api.Response {
						return api.NewResponse(http.StatusOK, nil)
					},
				}, {
					Path:   "test-authenticated-endpoint",
					Method: http.MethodGet,
					Auth:   true,
					Handle: func(r *api.Request) *api.Response {
						return api.NewResponse(http.StatusOK, nil)
					},
				}},
			})
	})

	It("serves the index page for the react app", func() {
		stop := server.Start()
		defer stop()

		client := &http.Client{
			Timeout: 15 * time.Second,
		}

		resp, err := client.Get(fmt.Sprintf("http://localhost:%s", port))
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		body, err := io.ReadAll(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(body)).To(ContainSubstring("Test HTML"))
	})

	It("serves unauthenticated api pages", func() {
		stop := server.Start()
		defer stop()

		client := &http.Client{
			Timeout: 15 * time.Second,
		}

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%s/api/v1/test-unauthenticated-endpoint", port), nil)
		Expect(err).ToNot(HaveOccurred())

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
	})

	It("serves authenticated api pages if the user is authenticated", func() {
		stop := server.Start()
		defer stop()

		client := &http.Client{
			Timeout: 15 * time.Second,
		}

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%s/api/v1/test-authenticated-endpoint", port), nil)
		Expect(err).ToNot(HaveOccurred())

		req.Header.Set("Authorization", "bearer token")

		resp, err := client.Do(req)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))
	})

	It("serves the static files directly", func() {
		stop := server.Start()
		defer stop()

		client := &http.Client{
			Timeout: 15 * time.Second,
		}

		resp, err := client.Get(fmt.Sprintf("http://localhost:%s/static.html", port))
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		body, err := io.ReadAll(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(string(body)).To(ContainSubstring("Static HTML"))
	})

	It("returns 404 if the api page does not exist", func() {
		stop := server.Start()
		defer stop()

		client := &http.Client{
			Timeout: 15 * time.Second,
		}

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%s/api/v1/non-existent", port), nil)
		Expect(err).ToNot(HaveOccurred())

		req.Header.Set("Authorization", "bearer token")

		resp, err := client.Do(req)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
	})

	It("returns unauthorized if the user is not authenticated for an endpoint that requires auth", func() {
		stop := server.Start()
		defer stop()

		client := &http.Client{
			Timeout: 15 * time.Second,
		}

		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%s/api/v1/test-authenticated-endpoint", port), nil)
		Expect(err).ToNot(HaveOccurred())

		resp, err := client.Do(req)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
	})
})

type mockTokenValidator struct {
	validateToken func(r *http.Request) (*token.AccessDetails, error)
}

func (m *mockTokenValidator) ValidateToken(r *http.Request) (*token.AccessDetails, error) {
	return m.validateToken(r)
}

type mockAccessDetailsRetriever struct {
	retrieveTokenDetails func(details *token.AccessDetails) (int64, error)
}

func (m *mockAccessDetailsRetriever) RetrieveTokenDetails(details *token.AccessDetails) (int64, error) {
	return m.retrieveTokenDetails(details)
}
