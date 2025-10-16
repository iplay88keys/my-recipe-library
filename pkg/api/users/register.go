package users

import (
	"context"
	"fmt"
	"net/http"

	"github.com/iplay88keys/my-recipe-library/pkg/api"
)

type RegisterResponse struct {
	Errors map[string]string `json:"errors,omitempty"`
}

type UserRegistrar interface {
	RegisterUser(ctx context.Context, username, email, password string) error
}

func Register(service UserRegistrar) *api.Endpoint {
	return &api.Endpoint{
		Path:   "users/register",
		Method: http.MethodPost,
		Handle: func(r *api.Request) *api.Response {
			var user RegisterRequest
			if err := r.Decode(&user); err != nil {
				fmt.Println("Error decoding json body for registration")
				return api.NewResponse(http.StatusBadRequest, nil)
			}

			validationErrors := user.Validate(false, false)
			if len(validationErrors) > 0 {
				resp := &RegisterResponse{
					Errors: validationErrors,
				}
				return api.NewResponse(http.StatusBadRequest, resp)
			}

			err := service.RegisterUser(r.Req.Context(), user.Username, user.Email, user.Password)
			if err != nil {
				if err.Error() == "username already exists" {
					resp := &RegisterResponse{
						Errors: map[string]string{
							"username": "Username already in use",
						},
					}
					return api.NewResponse(http.StatusBadRequest, resp)
				}
				if err.Error() == "email already exists" {
					resp := &RegisterResponse{
						Errors: map[string]string{
							"email": "Email already in use",
						},
					}
					return api.NewResponse(http.StatusBadRequest, resp)
				}

				fmt.Println("Failed to register user:", err.Error())
				return api.NewResponse(http.StatusInternalServerError, nil)
			}

			return api.NewResponse(http.StatusOK, nil)
		},
	}
}
