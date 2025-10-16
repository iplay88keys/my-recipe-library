package recipes

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/iplay88keys/my-recipe-library/pkg/api"
	"github.com/iplay88keys/my-recipe-library/pkg/services"
)

type RecipeResponse struct {
	ID          int64                 `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Creator     string                `json:"creator"`
	Servings    *int                  `json:"servings,omitempty"`
	PrepTime    *string               `json:"prep_time,omitempty"`
	CookTime    *string               `json:"cook_time,omitempty"`
	CoolTime    *string               `json:"cool_time,omitempty"`
	TotalTime   *string               `json:"total_time,omitempty"`
	Source      *string               `json:"source,omitempty"`
	Ingredients []*IngredientResponse `json:"ingredients"`
	Steps       []*StepResponse       `json:"steps"`
}

type IngredientResponse struct {
	Ingredient       string  `json:"ingredient"`
	IngredientNumber int     `json:"ingredient_number"`
	Amount           *string `json:"amount"`
	Measurement      *string `json:"measurement"`
	Preparation      *string `json:"preparation"`
}

type StepResponse struct {
	StepNumber   int    `json:"step_number"`
	Instructions string `json:"instructions"`
}

type ByIngredientNumber []*services.IngredientDetail

func (a ByIngredientNumber) Len() int      { return len(a) }
func (a ByIngredientNumber) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByIngredientNumber) Less(i, j int) bool {
	return a[i].OrderNum < a[j].OrderNum
}

type ByStepNumber []*services.StepDetail

func (a ByStepNumber) Len() int           { return len(a) }
func (a ByStepNumber) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByStepNumber) Less(i, j int) bool { return a[i].OrderNum < a[j].OrderNum }

type RecipeFetcher interface {
	GetRecipe(ctx context.Context, recipeID, userID int64) (*services.RecipeDetail, error)
}

func GetRecipe(service RecipeFetcher) *api.Endpoint {
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

			recipeDetail, err := service.GetRecipe(r.Req.Context(), recipeID, r.UserID)
			if err != nil {
				if err == sql.ErrNoRows {
					return api.NewResponse(http.StatusNotFound, nil)
				}

				fmt.Printf("Error getting recipe: %s\n", err.Error())
				return api.NewResponse(http.StatusInternalServerError, nil)
			}

			sort.Sort(ByIngredientNumber(recipeDetail.Ingredients))
			sort.Sort(ByStepNumber(recipeDetail.Steps))

			ingredients := make([]*IngredientResponse, len(recipeDetail.Ingredients))
			for i, ingredient := range recipeDetail.Ingredients {
				ingredients[i] = &IngredientResponse{
					Ingredient:       ingredient.Name,
					IngredientNumber: ingredient.OrderNum,
					Amount:           ingredient.Amount,
					Measurement:      ingredient.Unit,
					Preparation:      ingredient.Notes,
				}
			}

			steps := make([]*StepResponse, len(recipeDetail.Steps))
			for i, step := range recipeDetail.Steps {
				steps[i] = &StepResponse{
					StepNumber:   step.OrderNum,
					Instructions: step.Instructions,
				}
			}

			resp := &RecipeResponse{
				ID:          recipeDetail.ID,
				Name:        recipeDetail.Name,
				Description: recipeDetail.Description,
				Creator:     recipeDetail.Creator,
				Servings:    recipeDetail.Servings,
				PrepTime:    recipeDetail.PrepTime,
				CookTime:    recipeDetail.CookTime,
				CoolTime:    recipeDetail.CoolTime,
				TotalTime:   recipeDetail.TotalTime,
				Source:      recipeDetail.Source,
				Ingredients: ingredients,
				Steps:       steps,
			}

			return api.NewResponse(http.StatusOK, resp)
		},
	}
}
