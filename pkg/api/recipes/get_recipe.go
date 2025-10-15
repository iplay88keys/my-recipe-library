package recipes

import (
	"database/sql"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/iplay88keys/my-recipe-library/pkg/api"
	"github.com/iplay88keys/my-recipe-library/pkg/repositories"
)

type RecipeResponse struct {
	repositories.Recipe
	Ingredients []*repositories.Ingredient `json:"ingredients"`
	Steps       []*repositories.Step       `json:"steps"`
}

type ByIngredientNumber []*repositories.Ingredient

func (a ByIngredientNumber) Len() int      { return len(a) }
func (a ByIngredientNumber) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByIngredientNumber) Less(i, j int) bool {
	return *a[i].IngredientNumber < *a[j].IngredientNumber
}

type ByStepNumber []*repositories.Step

func (a ByStepNumber) Len() int           { return len(a) }
func (a ByStepNumber) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByStepNumber) Less(i, j int) bool { return *a[i].StepNumber < *a[j].StepNumber }

type getRecipe func(recipeID, userID int64) (*repositories.Recipe, error)
type getIngredientsForRecipe func(recipeID int64) ([]*repositories.Ingredient, error)
type getStepsForRecipe func(recipeID int64) ([]*repositories.Step, error)

func GetRecipe(getRecipe getRecipe, getIngredientsForRecipe getIngredientsForRecipe, getStepsForRecipe getStepsForRecipe) *api.Endpoint {
	return &api.Endpoint{
		Path:   "recipes/{id}",
		Method: http.MethodGet,
		Auth:   true,
		Handle: func(r *api.Request) *api.Response {
			idStr := r.Req.PathValue("id")
			if idStr == "" {
				fmt.Printf("Recipe endpoint missing id\n")
				return api.NewResponse(http.StatusBadRequest, nil)
			}

			recipeID, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				fmt.Printf("Recipe endpoint invalid id: %s\n", err.Error())
				return api.NewResponse(http.StatusBadRequest, nil)
			}

			recipe, err := getRecipe(recipeID, r.UserID)
			if err != nil {
				if err == sql.ErrNoRows {
					return api.NewResponse(http.StatusNotFound, nil)
				}

				fmt.Printf("Error getting recipe: %s\n", err.Error())
				return api.NewResponse(http.StatusInternalServerError, nil)
			}

			recipeIngredients, err := getIngredientsForRecipe(recipeID)
			if err != nil {
				fmt.Printf("Error getting ingredients for recipe: %s\n", err.Error())
				return api.NewResponse(http.StatusInternalServerError, nil)
			}

			recipeSteps, err := getStepsForRecipe(recipeID)
			if err != nil {
				fmt.Printf("Error getting steps for recipe: %s\n", err.Error())
				return api.NewResponse(http.StatusInternalServerError, nil)
			}

			sort.Sort(ByIngredientNumber(recipeIngredients))
			sort.Sort(ByStepNumber(recipeSteps))

			resp := &RecipeResponse{
				Recipe:      *recipe,
				Ingredients: recipeIngredients,
				Steps:       recipeSteps,
			}

			return api.NewResponse(http.StatusOK, resp)
		},
	}
}
