package users_test

import (
	"errors"
	"net/http"

	"github.com/iplay88keys/my-recipe-library/pkg/api"

	"github.com/iplay88keys/my-recipe-library/pkg/token"

	"github.com/iplay88keys/my-recipe-library/pkg/api/users"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("logout", func() {
	It("logs a user out", func() {
		validateToken := func(r *http.Request) (*token.AccessDetails, error) {
			return &token.AccessDetails{
				AccessUuid: "some-uuid",
			}, nil
		}

		deleteTokenDetails := func(uuid string) error {
			return nil
		}

		req, err := http.NewRequest(http.MethodPost, "/users/logout", nil)
		Expect(err).ToNot(HaveOccurred())

		resp := users.Logout(validateToken, deleteTokenDetails).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusOK))
	})

	It("returns unauthorized if the token cannot be validated", func() {
		validateToken := func(r *http.Request) (*token.AccessDetails, error) {
			return nil, errors.New("validation error")
		}

		deleteTokenDetails := func(uuid string) error {
			return nil
		}

		req, err := http.NewRequest(http.MethodPost, "/users/logout", nil)
		Expect(err).ToNot(HaveOccurred())

		resp := users.Logout(validateToken, deleteTokenDetails).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
	})

	It("returns unauthorized if the token cannot be deleted", func() {
		validateToken := func(r *http.Request) (*token.AccessDetails, error) {
			return &token.AccessDetails{
				AccessUuid: "some-uuid",
			}, nil
		}

		deleteTokenDetails := func(uuid string) error {
			return errors.New("token deletion failed")
		}

		req, err := http.NewRequest(http.MethodPost, "/users/logout", nil)
		Expect(err).ToNot(HaveOccurred())

		resp := users.Logout(validateToken, deleteTokenDetails).Handle(&api.Request{
			Req: req,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
	})
})
