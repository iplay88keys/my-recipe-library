package users

import (
	"net/http"

	"github.com/iplay88keys/my-recipe-library/pkg/token"

	"github.com/iplay88keys/my-recipe-library/pkg/api"
)

type validateToken func(r *http.Request) (*token.AccessDetails, error)
type deleteTokenDetails func(uuid string) (int64, error)

func Logout(validateToken validateToken, deleteTokenDetails deleteTokenDetails) *api.Endpoint {
	return &api.Endpoint{
		Path:   "users/logout",
		Method: http.MethodPost,
		Auth:   true,
		Handle: func(r *api.Request) *api.Response {
			details, err := validateToken(r.Req)
			if err != nil {
				return api.NewResponse(http.StatusUnauthorized, nil)
			}

			userID, err := deleteTokenDetails(details.AccessUuid)
			if err != nil || userID == 0 {
				return api.NewResponse(http.StatusUnauthorized, nil)
			}

			return api.NewResponse(http.StatusOK, nil)
		},
	}
}
