package recipes

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/iplay88keys/my-recipe-library/pkg/api"
	"github.com/iplay88keys/my-recipe-library/pkg/repositories"
)

type RecipeListResponse struct {
	Recipes []*repositories.Recipe `json:"recipes"`
}

type listRecipes func(userID int64) ([]*repositories.Recipe, error)

func ListRecipes(listRecipes listRecipes) *api.Endpoint {
	return &api.Endpoint{
		Path:   "recipes",
		Method: http.MethodGet,
		Auth:   true,
		Handle: func(r *api.Request) *api.Response {
			recipes, err := listRecipes(r.UserID)
			if err != nil {
				if err == sql.ErrNoRows {
					return api.NewResponse(http.StatusNoContent, nil)
				}

				fmt.Printf("Error listing recipes: %s\n", err.Error())
				return api.NewResponse(http.StatusInternalServerError, nil)
			}

			resp := &RecipeListResponse{
				Recipes: recipes,
			}

			return api.NewResponse(http.StatusOK, resp)
		},
	}
}
