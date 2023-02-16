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

var _ = Describe("ListRecipes", func() {
	It("returns the list of recipes", func() {
		listRecipes := func(userID int64) ([]*repositories.Recipe, error) {
			return []*repositories.Recipe{{
				ID:          Int64Pointer(1),
				Name:        StringPointer("First"),
				Description: StringPointer("One"),
				Creator:     StringPointer("Some Creator"),
				Source:      StringPointer("Some Website"),
			}, {
				ID:          Int64Pointer(2),
				Name:        StringPointer("Second"),
				Description: StringPointer("Two"),
				Creator:     StringPointer("Some Creator"),
				Source:      nil,
			}}, nil
		}

		req, err := http.NewRequest(http.MethodGet, "/recipes", nil)
		Expect(err).ToNot(HaveOccurred())

		resp := recipes.ListRecipes(listRecipes).Handle(&api.Request{
			Req:    req,
			UserID: 2,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		respBody, err := json.Marshal(resp.Body)
		Expect(err).ToNot(HaveOccurred())
		Expect(respBody).To(MatchJSON(`{
            "recipes": [{
                "id": 1,
                "name": "First",
                "description": "One",
                "creator": "Some Creator",
                "source": "Some Website"
            }, {
                "id": 2,
                "name": "Second",
                "description": "Two",
                "creator": "Some Creator"
            }]
        }`))
	})

	It("returns no content if there are no recipes", func() {
		listRecipes := func(userID int64) ([]*repositories.Recipe, error) {
			return nil, sql.ErrNoRows
		}

		req, err := http.NewRequest(http.MethodGet, "/recipes", nil)
		Expect(err).ToNot(HaveOccurred())

		resp := recipes.ListRecipes(listRecipes).Handle(&api.Request{
			Req:    req,
			UserID: 2,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
	})

	It("returns an error if the repository call fails", func() {
		listRecipes := func(userID int64) ([]*repositories.Recipe, error) {
			return nil, errors.New("some error")
		}

		req, err := http.NewRequest(http.MethodGet, "/recipes", nil)
		Expect(err).ToNot(HaveOccurred())

		resp := recipes.ListRecipes(listRecipes).Handle(&api.Request{
			Req:    req,
			UserID: 2,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
	})
})
