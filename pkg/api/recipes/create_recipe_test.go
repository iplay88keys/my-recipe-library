package recipes_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/iplay88keys/my-recipe-library/pkg/api"

	"github.com/iplay88keys/my-recipe-library/pkg/api/recipes"
	"github.com/iplay88keys/my-recipe-library/pkg/repositories"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateRecipe", func() {
	It("returns any validation", func() {
		addRecipe := func(recipe *repositories.Recipe, userID int64) (int64, error) {
			return 1, nil
		}

		body := []byte(`{
            "name": "Root Beer Float",
            "description": "Delicious",
            "servings": 1,
            "prep_time": "5 m",
            "cook_time": "0 m",
            "cool_time": "0 m",
            "total_time": "5 m",
            "source": "Some Book"
        }`)

		req, err := http.NewRequest(http.MethodPost, "/recipes", bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		resp := recipes.CreateRecipe(addRecipe).Handle(&api.Request{
			Req:    req,
			UserID: 2,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusCreated))

		respBody, err := json.Marshal(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(respBody).To(MatchJSON(`{
            "recipe_id": 1
        }`))
	})

	It("returns any validation errors", func() {
		addRecipe := func(recipe *repositories.Recipe, userID int64) (int64, error) {
			return 1, nil
		}

		req, err := http.NewRequest(http.MethodPost, "/recipes", bytes.NewBuffer([]byte("{}")))
		Expect(err).ToNot(HaveOccurred())

		resp := recipes.CreateRecipe(addRecipe).Handle(&api.Request{
			Req:    req,
			UserID: 2,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

		respBody, err := json.Marshal(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(respBody).To(MatchJSON(`{
            "errors": {
                "name": "Required",
                "description": "Required",
                "servings": "Required"
            }
        }`))
	})

	It("returns an error if the recipe repository call fails", func() {
		addRecipe := func(recipe *repositories.Recipe, userID int64) (int64, error) {
			return 1, errors.New("some error")
		}

		body := []byte(`{
            "name": "Root Beer Float",
            "description": "Delicious",
            "servings": 1,
            "prep_time": "5 m",
            "cook_time": "0 m",
            "cool_time": "0 m",
            "total_time": "5 m",
            "source": "Some Book"
        }`)

		req, err := http.NewRequest(http.MethodPost, "/recipes", bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		resp := recipes.CreateRecipe(addRecipe).Handle(&api.Request{
			Req:    req,
			UserID: 2,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
	})
})
