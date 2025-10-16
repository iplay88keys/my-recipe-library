package recipes

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/iplay88keys/my-recipe-library/pkg/api"
	"github.com/iplay88keys/my-recipe-library/pkg/services"
)

type RecipeListResponse struct {
	Recipes []*RecipeSummaryResponse `json:"recipes"`
}

type RecipeSummaryResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type RecipeLister interface {
	ListRecipes(ctx context.Context, userID int64) ([]*services.RecipeSummary, error)
}

func ListRecipes(service RecipeLister) *api.Endpoint {
	return &api.Endpoint{
		Path:   "recipes",
		Method: http.MethodGet,
		Auth:   true,
		Handle: func(r *api.Request) *api.Response {
			recipeSummaries, err := service.ListRecipes(r.Req.Context(), r.UserID)
			if err != nil {
				if err == sql.ErrNoRows {
					return api.NewResponse(http.StatusNoContent, nil)
				}

				fmt.Printf("Error listing recipes: %s\n", err.Error())
				return api.NewResponse(http.StatusInternalServerError, nil)
			}

			recipes := make([]*RecipeSummaryResponse, len(recipeSummaries))
			for i, summary := range recipeSummaries {
				recipes[i] = &RecipeSummaryResponse{
					ID:          summary.ID,
					Name:        summary.Name,
					Description: summary.Description,
				}
			}

			resp := &RecipeListResponse{
				Recipes: recipes,
			}

			return api.NewResponse(http.StatusOK, resp)
		},
	}
}
