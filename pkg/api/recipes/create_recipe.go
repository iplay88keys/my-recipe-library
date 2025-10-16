package recipes

import (
	"context"
	"fmt"
	"net/http"

	. "github.com/iplay88keys/my-recipe-library/pkg/helpers"

	"github.com/iplay88keys/my-recipe-library/pkg/api"
	"github.com/iplay88keys/my-recipe-library/pkg/services"
)

type CreateRecipeResponse struct {
	RecipeID int64             `json:"recipe_id,omitempty"`
	Errors   map[string]string `json:"errors,omitempty"`
}

type RecipeCreator interface {
	CreateRecipe(ctx context.Context, userID int64, recipe *services.RecipeInput) (int64, error)
}

func CreateRecipe(service RecipeCreator) *api.Endpoint {
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

			ingredients := make([]*services.IngredientInput, len(recipe.Ingredients))
			for i, ingredient := range recipe.Ingredients {
				ingredients[i] = &services.IngredientInput{
					Name:     ingredient.Name,
					Amount:   ingredient.Amount,
					Unit:     ingredient.Unit,
					Notes:    ingredient.Notes,
					OrderNum: ingredient.OrderNum,
				}
			}

			steps := make([]*services.StepInput, len(recipe.Steps))
			for i, step := range recipe.Steps {
				steps[i] = &services.StepInput{
					Instructions: step.Instructions,
					OrderNum:     step.OrderNum,
					Notes:        &step.Notes,
				}
			}

			recipeInput := &services.RecipeInput{
				Name:        recipe.Name,
				Description: recipe.Description,
				Servings:    IntPointer(recipe.Servings),
				PrepTime:    StringPointer(recipe.PrepTime),
				CookTime:    StringPointer(recipe.CookTime),
				CoolTime:    StringPointer(recipe.CoolTime),
				TotalTime:   StringPointer(recipe.TotalTime),
				Source:      StringPointer(recipe.Source),
				Ingredients: ingredients,
				Steps:       steps,
			}

			recipeID, err := service.CreateRecipe(r.Req.Context(), r.UserID, recipeInput)
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
