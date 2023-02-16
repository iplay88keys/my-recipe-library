package users

import (
	"fmt"
	"net/http"

	"github.com/iplay88keys/my-recipe-library/pkg/api"
)

type RegisterResponse struct {
	Errors map[string]string `json:"errors,omitempty"`
}

type existsByUsername func(username string) (bool, error)
type existsByEmail func(email string) (bool, error)
type insertUser func(username, email, password string) (int64, error)

func Register(existsByUsername existsByUsername, existsByEmail existsByEmail, insertUser insertUser) *api.Endpoint {
	return &api.Endpoint{
		Path:   "users/register",
		Method: http.MethodPost,
		Handle: func(r *api.Request) *api.Response {
			var user RegisterRequest
			if err := r.Decode(&user); err != nil {
				fmt.Println("Error decoding json body for registration")
				return api.NewResponse(http.StatusBadRequest, nil)
			}

			usernameExists, err := existsByUsername(user.Username)
			if err != nil {
				fmt.Println("Error checking if user exists by username for registration")
				return api.NewResponse(http.StatusInternalServerError, nil)
			}

			emailExists, err := existsByEmail(user.Email)
			if err != nil {
				fmt.Println("Error checking if user exists by email for registration")
				return api.NewResponse(http.StatusInternalServerError, nil)
			}

			validationErrors := user.Validate(usernameExists, emailExists)
			if len(validationErrors) > 0 {
				resp := &RegisterResponse{
					Errors: validationErrors,
				}

				return api.NewResponse(http.StatusBadRequest, resp)
			}

			_, err = insertUser(user.Username, user.Email, user.Password)
			if err != nil {
				fmt.Println("Failed to register user")
				return api.NewResponse(http.StatusInternalServerError, nil)
			}

			return api.NewResponse(http.StatusOK, nil)
		},
	}
}
