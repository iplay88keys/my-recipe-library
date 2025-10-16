package users

import (
	"fmt"
	"net/http"

	"github.com/iplay88keys/my-recipe-library/pkg/api"
	"github.com/iplay88keys/my-recipe-library/pkg/token"
)

type UserLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	AccessToken  string            `json:"access_token,omitempty"`
	RefreshToken string            `json:"refresh_token,omitempty"`
	Errors       map[string]string `json:"errors,omitempty"`
}

type LoginService interface {
	Verify(login, password string) (bool, int64, error)
	CreateToken(userID int64) (*token.Details, error)
	StoreTokenDetails(userID int64, details *token.Details) error
}

func Login(service LoginService) *api.Endpoint {
	return &api.Endpoint{
		Path:   "users/login",
		Method: http.MethodPost,
		Handle: func(r *api.Request) *api.Response {
			var user UserLoginRequest
			if err := r.Decode(&user); err != nil {
				fmt.Println("Error decoding json body for login")
				return api.NewResponse(http.StatusBadRequest, nil)
			}

			valid, userID, err := service.Verify(user.Login, user.Password)
			if err != nil {
				fmt.Println("Error logging user in")
				return api.NewResponse(http.StatusInternalServerError, nil)
			}

			if !valid {
				errors := make(map[string]string)
				errors["alert"] = "Invalid login credentials"

				resp := &UserLoginResponse{
					Errors: errors,
				}

				fmt.Println("Error validating user for login")
				return api.NewResponse(http.StatusUnauthorized, resp)
			}

			tokenDetails, err := service.CreateToken(userID)
			if err != nil {
				fmt.Printf("Error creating token for user login: %s\n", err.Error())
				return api.NewResponse(http.StatusInternalServerError, nil)
			}

			err = service.StoreTokenDetails(userID, tokenDetails)
			if err != nil {
				fmt.Printf("Error saving token for user login: %s\n", err.Error())
				return api.NewResponse(http.StatusInternalServerError, nil)
			}

			resp := &UserLoginResponse{
				AccessToken:  tokenDetails.AccessToken,
				RefreshToken: tokenDetails.RefreshToken,
			}

			return api.NewResponse(http.StatusOK, resp)
		},
	}
}
