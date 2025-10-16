package recipes_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/iplay88keys/my-recipe-library/pkg/api"
	"github.com/iplay88keys/my-recipe-library/pkg/api/recipes"
	. "github.com/iplay88keys/my-recipe-library/pkg/helpers"
	"github.com/iplay88keys/my-recipe-library/pkg/services"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetRecipe", func() {
	It("returns a recipe", func() {
		recipeDetail := &services.RecipeDetail{
			ID:          1,
			Name:        "Root Beer Float",
			Description: "Delicious",
			Creator:     "User1",
			Servings:    IntPointer(1),
			PrepTime:    StringPointer("5 m"),
			CookTime:    StringPointer("0 m"),
			CoolTime:    StringPointer("0 m"),
			TotalTime:   StringPointer("5 m"),
			Source:      StringPointer("Some Book"),
			Ingredients: []*services.IngredientDetail{{
				Name:     "Vanilla Ice Cream",
				Amount:   StringPointer("1"),
				Unit:     StringPointer("Scoop"),
				Notes:    StringPointer("Frozen"),
				OrderNum: 1,
			}, {
				Name:     "Root Beer",
				Amount:   nil,
				Unit:     nil,
				Notes:    nil,
				OrderNum: 2,
			}},
			Steps: []*services.StepDetail{{
				Instructions: "Place ice cream in glass.",
				OrderNum:     1,
				Notes:        nil,
			}, {
				Instructions: "Top with Root Beer.",
				OrderNum:     2,
				Notes:        nil,
			}},
		}

		fakeService := &mockRecipeFetcher{
			getRecipe: func(ctx context.Context, recipeID, userID int64) (*services.RecipeDetail, error) {
				return recipeDetail, nil
			},
		}

		req := httptest.NewRequest(http.MethodGet, "/recipes/1", nil)
		req.SetPathValue("id", "1")

		resp := recipes.GetRecipe(fakeService).Handle(&api.Request{
			Req:    req,
			UserID: 2,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		respBody, err := json.Marshal(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(respBody).To(MatchJSON(`{
            "id": 1,
            "name": "Root Beer Float",
            "description": "Delicious",
            "creator": "User1",
            "servings": 1,
            "prep_time": "5 m",
            "cook_time": "0 m",
            "cool_time": "0 m",
            "total_time": "5 m",
            "source": "Some Book",
            "ingredients": [{
                "ingredient": "Vanilla Ice Cream",
                "ingredient_number": 1,
                "amount": "1",
                "measurement": "Scoop",
                "preparation": "Frozen"
            }, {
                "ingredient": "Root Beer",
                "ingredient_number": 2,
                "amount": null,
                "measurement": null,
                "preparation": null
            }],
            "steps": [{
                "step_number": 1,
                "instructions": "Place ice cream in glass."
            }, {
                "step_number": 2,
                "instructions": "Top with Root Beer."
            }]
        }`))
	})

	It("sorts the recipe ingredients by ingredient number", func() {
		recipeDetail := &services.RecipeDetail{
			ID:          1,
			Name:        "Root Beer Float",
			Description: "Delicious",
			Creator:     "User1",
			Ingredients: []*services.IngredientDetail{{
				Name:     "Root Beer",
				Amount:   nil,
				Unit:     nil,
				Notes:    nil,
				OrderNum: 2,
			}, {
				Name:     "Vanilla Ice Cream",
				Amount:   StringPointer("1"),
				Unit:     StringPointer("Scoop"),
				Notes:    StringPointer("Frozen"),
				OrderNum: 1,
			}},
			Steps: []*services.StepDetail{},
		}

		fakeService := &mockRecipeFetcher{
			getRecipe: func(ctx context.Context, recipeID, userID int64) (*services.RecipeDetail, error) {
				return recipeDetail, nil
			},
		}

		req := httptest.NewRequest(http.MethodGet, "/recipes/1", nil)
		req.SetPathValue("id", "1")

		resp := recipes.GetRecipe(fakeService).Handle(&api.Request{
			Req:    req,
			UserID: 2,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		respBody, err := json.Marshal(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(respBody).To(MatchJSON(`{
            "id": 1,
            "name": "Root Beer Float",
            "description": "Delicious",
            "creator": "User1",
            "ingredients": [{
                "ingredient": "Vanilla Ice Cream",
                "ingredient_number": 1,
                "amount": "1",
                "measurement": "Scoop",
                "preparation": "Frozen"
            }, {
                "ingredient": "Root Beer",
                "ingredient_number": 2,
                "amount": null,
                "measurement": null,
                "preparation": null
            }],
            "steps": []
        }`))
	})

	It("sorts the recipe steps by step number", func() {
		recipeDetail := &services.RecipeDetail{
			ID:          1,
			Name:        "Root Beer Float",
			Description: "Delicious",
			Creator:     "User1",
			Ingredients: []*services.IngredientDetail{},
			Steps: []*services.StepDetail{{
				Instructions: "Top with Root Beer.",
				OrderNum:     2,
				Notes:        nil,
			}, {
				Instructions: "Place ice cream in glass.",
				OrderNum:     1,
				Notes:        nil,
			}},
		}

		fakeService := &mockRecipeFetcher{
			getRecipe: func(ctx context.Context, recipeID, userID int64) (*services.RecipeDetail, error) {
				return recipeDetail, nil
			},
		}

		req := httptest.NewRequest(http.MethodGet, "/recipes/1", nil)
		req.SetPathValue("id", "1")

		resp := recipes.GetRecipe(fakeService).Handle(&api.Request{
			Req:    req,
			UserID: 2,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		respBody, err := json.Marshal(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(respBody).To(MatchJSON(`{
            "id": 1,
            "name": "Root Beer Float",
            "description": "Delicious",
            "creator": "User1",
            "ingredients": [],
            "steps": [{
                "step_number": 1,
                "instructions": "Place ice cream in glass."
            }, {
                "step_number": 2,
                "instructions": "Top with Root Beer."
            }]
        }`))
	})

	It("returns an error if the recipe repository returns no rows", func() {
		fakeService := &mockRecipeFetcher{
			getRecipe: func(ctx context.Context, recipeID, userID int64) (*services.RecipeDetail, error) {
				return nil, sql.ErrNoRows
			},
		}

		req := httptest.NewRequest(http.MethodGet, "/recipes/1", nil)
		req.SetPathValue("id", "1")

		resp := recipes.GetRecipe(fakeService).Handle(&api.Request{
			Req:    req,
			UserID: 2,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
	})

	It("returns an error if the recipe repository call fails", func() {
		fakeService := &mockRecipeFetcher{
			getRecipe: func(ctx context.Context, recipeID, userID int64) (*services.RecipeDetail, error) {
				return nil, errors.New("some error")
			},
		}

		req := httptest.NewRequest(http.MethodGet, "/recipes/1", nil)
		req.SetPathValue("id", "1")

		resp := recipes.GetRecipe(fakeService).Handle(&api.Request{
			Req:    req,
			UserID: 2,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
	})

	It("returns an error if the ingredients repository call fails", func() {
		fakeService := &mockRecipeFetcher{
			getRecipe: func(ctx context.Context, recipeID, userID int64) (*services.RecipeDetail, error) {
				return nil, errors.New("some error")
			},
		}

		req := httptest.NewRequest(http.MethodGet, "/recipes/1", nil)
		req.SetPathValue("id", "1")

		resp := recipes.GetRecipe(fakeService).Handle(&api.Request{
			Req:    req,
			UserID: 2,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
	})

	It("returns an error if the steps repository call fails", func() {
		fakeService := &mockRecipeFetcher{
			getRecipe: func(ctx context.Context, recipeID, userID int64) (*services.RecipeDetail, error) {
				return nil, errors.New("some error")
			},
		}

		req := httptest.NewRequest(http.MethodGet, "/recipes/1", nil)
		req.SetPathValue("id", "1")

		resp := recipes.GetRecipe(fakeService).Handle(&api.Request{
			Req:    req,
			UserID: 2,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
	})

	It("returns an error if the provided route variable is not a number", func() {
		fakeService := &mockRecipeFetcher{
			getRecipe: func(ctx context.Context, recipeID, userID int64) (*services.RecipeDetail, error) {
				return nil, nil
			},
		}

		req := httptest.NewRequest(http.MethodGet, "/recipes/not-a-number", nil)
		req.SetPathValue("id", "not-a-number")

		resp := recipes.GetRecipe(fakeService).Handle(&api.Request{
			Req:    req,
			UserID: 2,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
	})
})

type mockRecipeFetcher struct {
	getRecipe func(ctx context.Context, recipeID, userID int64) (*services.RecipeDetail, error)
}

func (m *mockRecipeFetcher) GetRecipe(ctx context.Context, recipeID, userID int64) (*services.RecipeDetail, error) {
	return m.getRecipe(ctx, recipeID, userID)
}
