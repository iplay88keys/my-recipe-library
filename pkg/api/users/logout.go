package users

import (
	"net/http"

	"github.com/iplay88keys/my-recipe-library/pkg/api"
	"github.com/iplay88keys/my-recipe-library/pkg/token"
)

type LogoutService interface {
	ValidateToken(r *http.Request) (*token.AccessDetails, error)
	DeleteTokenDetails(uuid string) error
}

func Logout(service LogoutService) *api.Endpoint {
	return &api.Endpoint{
		Path:   "users/logout",
		Method: http.MethodPost,
		Auth:   true,
		Handle: func(r *api.Request) *api.Response {
			details, err := service.ValidateToken(r.Req)
			if err != nil {
				return api.NewResponse(http.StatusUnauthorized, nil)
			}

			err = service.DeleteTokenDetails(details.AccessUuid)
			if err != nil {
				return api.NewResponse(http.StatusUnauthorized, nil)
			}

			return api.NewResponse(http.StatusOK, nil)
		},
	}
}
