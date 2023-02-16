package recipes_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/iplay88keys/my-recipe-library/pkg/api"

	"github.com/iplay88keys/my-recipe-library/pkg/api/recipes"
	. "github.com/iplay88keys/my-recipe-library/pkg/helpers"
	"github.com/iplay88keys/my-recipe-library/pkg/repositories"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetRecipe", func() {
	It("returns a recipe", func() {
		getRecipe := func(recipeID, userID int64) (*repositories.Recipe, error) {
			return &repositories.Recipe{
				ID:          Int64Pointer(1),
				Name:        StringPointer("Root Beer Float"),
				Description: StringPointer("Delicious"),
				Creator:     StringPointer("User1"),
				Servings:    IntPointer(1),
				PrepTime:    StringPointer("5 m"),
				CookTime:    StringPointer("0 m"),
				CoolTime:    StringPointer("0 m"),
				TotalTime:   StringPointer("5 m"),
				Source:      StringPointer("Some Book"),
			}, nil
		}

		getIngredients := func(recipeID int64) ([]*repositories.Ingredient, error) {
			return []*repositories.Ingredient{{
				Ingredient:       StringPointer("Vanilla Ice Cream"),
				IngredientNumber: IntPointer(1),
				Amount:           StringPointer("1"),
				Measurement:      StringPointer("Scoop"),
				Preparation:      StringPointer("Frozen"),
			}, {
				Ingredient:       StringPointer("Root Beer"),
				IngredientNumber: IntPointer(2),
				Amount:           nil,
				Measurement:      nil,
				Preparation:      nil,
			}}, nil
		}

		getSteps := func(recipeID int64) ([]*repositories.Step, error) {
			return []*repositories.Step{{
				StepNumber:   IntPointer(1),
				Instructions: StringPointer("Place ice cream in glass."),
			}, {
				StepNumber:   IntPointer(2),
				Instructions: StringPointer("Top with Root Beer."),
			}}, nil
		}

		req, err := http.NewRequest(http.MethodGet, "/recipes/1", nil)
		Expect(err).ToNot(HaveOccurred())

		resp := recipes.GetRecipe(getRecipe, getIngredients, getSteps).Handle(&api.Request{
			Req:    req,
			UserID: 2,
			Vars:   map[string]string{"id": "1"},
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
		getRecipe := func(recipeID, userID int64) (*repositories.Recipe, error) {
			return &repositories.Recipe{
				ID:          Int64Pointer(1),
				Name:        StringPointer("Root Beer Float"),
				Description: StringPointer("Delicious"),
			}, nil
		}

		getIngredients := func(recipeID int64) ([]*repositories.Ingredient, error) {
			return []*repositories.Ingredient{{
				Ingredient:       StringPointer("Root Beer"),
				IngredientNumber: IntPointer(2),
				Amount:           nil,
				Measurement:      nil,
				Preparation:      nil,
			}, {
				Ingredient:       StringPointer("Vanilla Ice Cream"),
				IngredientNumber: IntPointer(1),
				Amount:           StringPointer("1"),
				Measurement:      StringPointer("Scoop"),
				Preparation:      StringPointer("Frozen"),
			}}, nil
		}

		getSteps := func(recipeID int64) ([]*repositories.Step, error) {
			return []*repositories.Step{}, nil
		}

		req, err := http.NewRequest(http.MethodGet, "/recipes/1", nil)
		Expect(err).ToNot(HaveOccurred())

		resp := recipes.GetRecipe(getRecipe, getIngredients, getSteps).Handle(&api.Request{
			Req:    req,
			UserID: 2,
			Vars:   map[string]string{"id": "1"},
		})

		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		respBody, err := json.Marshal(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(respBody).To(MatchJSON(`{
            "id": 1,
            "name": "Root Beer Float",
            "description": "Delicious",
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
		getRecipe := func(recipeID, userID int64) (*repositories.Recipe, error) {
			return &repositories.Recipe{
				ID:          Int64Pointer(1),
				Name:        StringPointer("Root Beer Float"),
				Description: StringPointer("Delicious"),
			}, nil
		}

		getIngredients := func(recipeID int64) ([]*repositories.Ingredient, error) {
			return []*repositories.Ingredient{}, nil
		}

		getSteps := func(recipeID int64) ([]*repositories.Step, error) {
			return []*repositories.Step{{
				StepNumber:   IntPointer(2),
				Instructions: StringPointer("Top with Root Beer."),
			}, {
				StepNumber:   IntPointer(1),
				Instructions: StringPointer("Place ice cream in glass."),
			}}, nil
		}

		req, err := http.NewRequest(http.MethodGet, "/recipes/1", nil)
		Expect(err).ToNot(HaveOccurred())

		resp := recipes.GetRecipe(getRecipe, getIngredients, getSteps).Handle(&api.Request{
			Req:    req,
			UserID: 2,
			Vars:   map[string]string{"id": "1"},
		})

		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		respBody, err := json.Marshal(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(respBody).To(MatchJSON(`{
            "id": 1,
            "name": "Root Beer Float",
            "description": "Delicious",
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
		getRecipe := func(recipeID, userID int64) (*repositories.Recipe, error) {
			return nil, sql.ErrNoRows
		}

		getIngredients := func(recipeID int64) ([]*repositories.Ingredient, error) {
			return []*repositories.Ingredient{}, nil
		}

		getSteps := func(recipeID int64) ([]*repositories.Step, error) {
			return []*repositories.Step{}, nil
		}

		req, err := http.NewRequest(http.MethodGet, "/recipes/1", nil)
		Expect(err).ToNot(HaveOccurred())

		resp := recipes.GetRecipe(getRecipe, getIngredients, getSteps).Handle(&api.Request{
			Req:    req,
			UserID: 2,
			Vars:   map[string]string{"id": "1"},
		})

		Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
	})

	It("returns an error if the recipe repository call fails", func() {
		getRecipe := func(recipeID, userID int64) (*repositories.Recipe, error) {
			return nil, errors.New("some error")
		}

		getIngredients := func(recipeID int64) ([]*repositories.Ingredient, error) {
			return []*repositories.Ingredient{}, nil
		}

		getSteps := func(recipeID int64) ([]*repositories.Step, error) {
			return []*repositories.Step{}, nil
		}

		req, err := http.NewRequest(http.MethodGet, "/recipes/1", nil)
		Expect(err).ToNot(HaveOccurred())

		resp := recipes.GetRecipe(getRecipe, getIngredients, getSteps).Handle(&api.Request{
			Req:    req,
			UserID: 2,
			Vars:   map[string]string{"id": "1"},
		})

		Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
	})

	It("returns an error if the ingredients repository call fails", func() {
		getRecipe := func(recipeID, userID int64) (*repositories.Recipe, error) {
			return &repositories.Recipe{}, nil
		}

		getIngredients := func(recipeID int64) ([]*repositories.Ingredient, error) {
			return nil, errors.New("some error")
		}

		getSteps := func(recipeID int64) ([]*repositories.Step, error) {
			return []*repositories.Step{}, nil
		}

		req, err := http.NewRequest(http.MethodGet, "/recipes/1", nil)
		Expect(err).ToNot(HaveOccurred())

		resp := recipes.GetRecipe(getRecipe, getIngredients, getSteps).Handle(&api.Request{
			Req:    req,
			UserID: 2,
			Vars:   map[string]string{"id": "1"},
		})

		Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
	})

	It("returns an error if the steps repository call fails", func() {
		getRecipe := func(recipeID, userID int64) (*repositories.Recipe, error) {
			return &repositories.Recipe{}, nil
		}

		getIngredients := func(recipeID int64) ([]*repositories.Ingredient, error) {
			return []*repositories.Ingredient{}, nil
		}

		getSteps := func(recipeID int64) ([]*repositories.Step, error) {
			return nil, errors.New("some error")
		}

		req, err := http.NewRequest(http.MethodGet, "/recipes/1", nil)
		Expect(err).ToNot(HaveOccurred())

		resp := recipes.GetRecipe(getRecipe, getIngredients, getSteps).Handle(&api.Request{
			Req:    req,
			UserID: 2,
			Vars:   map[string]string{"id": "1"},
		})

		Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
	})

	It("returns an error if the provided route variable is not a number", func() {
		getRecipe := func(recipeID, userID int64) (*repositories.Recipe, error) {
			return &repositories.Recipe{}, nil
		}

		getIngredients := func(recipeID int64) ([]*repositories.Ingredient, error) {
			return []*repositories.Ingredient{}, nil
		}

		getSteps := func(recipeID int64) ([]*repositories.Step, error) {
			return []*repositories.Step{}, nil
		}

		req, err := http.NewRequest(http.MethodGet, "/recipes/not-a-number", nil)
		Expect(err).ToNot(HaveOccurred())

		resp := recipes.GetRecipe(getRecipe, getIngredients, getSteps).Handle(&api.Request{
			Req:    req,
			UserID: 2,
			Vars:   map[string]string{"id": "not-a-number"},
		})

		Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
	})
})
