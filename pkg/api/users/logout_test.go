package users_test

import (
	"errors"
	"net/http"

	"github.com/iplay88keys/my-recipe-library/pkg/api"
	"github.com/iplay88keys/my-recipe-library/pkg/api/users"
	"github.com/iplay88keys/my-recipe-library/pkg/token"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("logout", func() {
	It("logs a user out", func() {
		fakeLogoutService := &mockLogoutService{
			validateToken: func(r *http.Request) (*token.AccessDetails, error) {
				return &token.AccessDetails{
					AccessUuid: "some-uuid",
				}, nil
			},
			deleteTokenDetails: func(uuid string) error {
				return nil
			},
		}

		req, err := http.NewRequest(http.MethodPost, "/users/logout", nil)
		Expect(err).ToNot(HaveOccurred())

		resp := users.Logout(fakeLogoutService).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusOK))
	})

	It("returns unauthorized if the token cannot be validated", func() {
		fakeLogoutService := &mockLogoutService{
			validateToken: func(r *http.Request) (*token.AccessDetails, error) {
				return nil, errors.New("validation error")
			},
			deleteTokenDetails: func(uuid string) error {
				return nil
			},
		}

		req, err := http.NewRequest(http.MethodPost, "/users/logout", nil)
		Expect(err).ToNot(HaveOccurred())

		resp := users.Logout(fakeLogoutService).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
	})

	It("returns unauthorized if the token cannot be deleted", func() {
		fakeLogoutService := &mockLogoutService{
			validateToken: func(r *http.Request) (*token.AccessDetails, error) {
				return &token.AccessDetails{
					AccessUuid: "some-uuid",
				}, nil
			},
			deleteTokenDetails: func(uuid string) error {
				return errors.New("token deletion failed")
			},
		}

		req, err := http.NewRequest(http.MethodPost, "/users/logout", nil)
		Expect(err).ToNot(HaveOccurred())

		resp := users.Logout(fakeLogoutService).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
	})
})

type mockLogoutService struct {
	validateToken      func(r *http.Request) (*token.AccessDetails, error)
	deleteTokenDetails func(uuid string) error
}

func (m *mockLogoutService) ValidateToken(r *http.Request) (*token.AccessDetails, error) {
	return m.validateToken(r)
}

func (m *mockLogoutService) DeleteTokenDetails(uuid string) error {
	return m.deleteTokenDetails(uuid)
}
