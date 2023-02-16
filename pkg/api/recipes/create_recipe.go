package recipes

import (
	"fmt"
	"net/http"

	. "github.com/iplay88keys/my-recipe-library/pkg/helpers"

	"github.com/iplay88keys/my-recipe-library/pkg/api"
	"github.com/iplay88keys/my-recipe-library/pkg/repositories"
)

type CreateRecipeResponse struct {
	RecipeID int64             `json:"recipe_id,omitempty"`
	Errors   map[string]string `json:"errors,omitempty"`
}

type createRecipe func(recipe *repositories.Recipe, userID int64) (int64, error)

func CreateRecipe(createRecipe createRecipe) *api.Endpoint {
	return &api.Endpoint{
		Path:   "recipes",
		Method: http.MethodPost,
		Auth:   true,
		Handle: func(r *api.Request) *api.Response {
			var recipe CreateRecipeRequest
			if err := r.Decode(&recipe); err != nil {
				fmt.Printf("Error decoding json body for add recipe: %s\n", err.Error())
				return api.NewResponse(http.StatusBadRequest, nil)
			}

			validationErrors := recipe.Validate()
			if len(validationErrors) > 0 {
				resp := &CreateRecipeResponse{
					Errors: validationErrors,
				}

				return api.NewResponse(http.StatusBadRequest, resp)
			}

			recipeID, err := createRecipe(&repositories.Recipe{
				Name:        StringPointer(recipe.Name),
				Description: StringPointer(recipe.Description),
				Servings:    IntPointer(recipe.Servings),
				PrepTime:    StringPointer(recipe.PrepTime),
				CookTime:    StringPointer(recipe.CookTime),
				CoolTime:    StringPointer(recipe.CoolTime),
				TotalTime:   StringPointer(recipe.TotalTime),
				Source:      StringPointer(recipe.Source),
			}, r.UserID)
			if err != nil {
				fmt.Printf("Error adding recipe: %s\n", err.Error())
				return api.NewResponse(http.StatusInternalServerError, nil)
			}

			resp := &CreateRecipeResponse{
				RecipeID: recipeID,
			}

			return api.NewResponse(http.StatusCreated, resp)
		},
	}
}
