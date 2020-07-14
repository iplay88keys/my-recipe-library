package recipes_test

import (
    "database/sql"
    "errors"
    "net/http"
    "net/http/httptest"

    "github.com/gorilla/mux"
    "golang.org/x/net/context"

    "github.com/iplay88keys/recipe-box/pkg/api/auth"

    "github.com/iplay88keys/recipe-box/pkg/api/recipes"
    . "github.com/iplay88keys/recipe-box/pkg/helpers"
    "github.com/iplay88keys/recipe-box/pkg/repositories"

    . "github.com/onsi/ginkgo"
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

        req := httptest.NewRequest("GET", "/recipes/1", nil)
        req = mux.SetURLVars(req, map[string]string{"id": "1"})
        req = req.WithContext(context.WithValue(req.Context(), auth.ContextUserKey, int64(2)))

        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(recipes.GetRecipe(getRecipe, getIngredients, getSteps).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusOK))
        Expect(rr.Body.String()).To(MatchJSON(`{
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

        req := httptest.NewRequest("GET", "/recipes/1", nil)
        req = mux.SetURLVars(req, map[string]string{"id": "1"})
        req = req.WithContext(context.WithValue(req.Context(), auth.ContextUserKey, int64(2)))

        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(recipes.GetRecipe(getRecipe, getIngredients, getSteps).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusOK))
        Expect(rr.Body.String()).To(MatchJSON(`{
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

        req := httptest.NewRequest("GET", "/recipes/1", nil)
        req = mux.SetURLVars(req, map[string]string{"id": "1"})
        req = req.WithContext(context.WithValue(req.Context(), auth.ContextUserKey, int64(2)))

        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(recipes.GetRecipe(getRecipe, getIngredients, getSteps).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusOK))
        Expect(rr.Body.String()).To(MatchJSON(`{
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

    It("returns an error if the userID does not exist on the context", func() {
        getRecipe := func(recipeID, userID int64) (*repositories.Recipe, error) {
            return nil, sql.ErrNoRows
        }

        getIngredients := func(recipeID int64) ([]*repositories.Ingredient, error) {
            return []*repositories.Ingredient{}, nil
        }

        getSteps := func(recipeID int64) ([]*repositories.Step, error) {
            return []*repositories.Step{}, nil
        }

        req := httptest.NewRequest("GET", "/recipes/1", nil)
        req = mux.SetURLVars(req, map[string]string{"id": "1"})

        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(recipes.GetRecipe(getRecipe, getIngredients, getSteps).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusInternalServerError))
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

        req := httptest.NewRequest("GET", "/recipes/1", nil)
        req = mux.SetURLVars(req, map[string]string{"id": "1"})
        req = req.WithContext(context.WithValue(req.Context(), auth.ContextUserKey, int64(2)))

        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(recipes.GetRecipe(getRecipe, getIngredients, getSteps).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusNotFound))
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

        req := httptest.NewRequest("GET", "/recipes/1", nil)
        req = mux.SetURLVars(req, map[string]string{"id": "1"})
        req = req.WithContext(context.WithValue(req.Context(), auth.ContextUserKey, int64(2)))

        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(recipes.GetRecipe(getRecipe, getIngredients, getSteps).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusInternalServerError))
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

        req := httptest.NewRequest("GET", "/recipes/1", nil)
        req = mux.SetURLVars(req, map[string]string{"id": "1"})
        req = req.WithContext(context.WithValue(req.Context(), auth.ContextUserKey, int64(2)))

        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(recipes.GetRecipe(getRecipe, getIngredients, getSteps).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusInternalServerError))
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

        req := httptest.NewRequest("GET", "/recipes/1", nil)
        req = mux.SetURLVars(req, map[string]string{"id": "1"})
        req = req.WithContext(context.WithValue(req.Context(), auth.ContextUserKey, int64(2)))

        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(recipes.GetRecipe(getRecipe, getIngredients, getSteps).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusInternalServerError))
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

        req := httptest.NewRequest("GET", "/recipes/not-a-number", nil)
        req = mux.SetURLVars(req, map[string]string{"id": "not-a-number"})
        req = req.WithContext(context.WithValue(req.Context(), auth.ContextUserKey, int64(2)))

        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(recipes.GetRecipe(getRecipe, getIngredients, getSteps).Handler)

        handler.ServeHTTP(rr, req)
        Expect(rr.Code).To(Equal(http.StatusInternalServerError))
    })
})
