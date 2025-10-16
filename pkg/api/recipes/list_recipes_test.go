package recipes_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/iplay88keys/my-recipe-library/pkg/api"
	"github.com/iplay88keys/my-recipe-library/pkg/api/recipes"
	"github.com/iplay88keys/my-recipe-library/pkg/services"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListRecipes", func() {
	It("returns the list of recipes", func() {
		recipeSummaries := []*services.RecipeSummary{{
			ID:          1,
			Name:        "First",
			Description: "One",
		}, {
			ID:          2,
			Name:        "Second",
			Description: "Two",
		}}

		fakeService := &mockRecipeLister{
			listRecipes: func(ctx context.Context, userID int64) ([]*services.RecipeSummary, error) {
				return recipeSummaries, nil
			},
		}

		req, err := http.NewRequest(http.MethodGet, "/recipes", nil)
		Expect(err).ToNot(HaveOccurred())

		resp := recipes.ListRecipes(fakeService).Handle(&api.Request{
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
                "description": "One"
            }, {
                "id": 2,
                "name": "Second",
                "description": "Two"
            }]
        }`))
	})

	It("returns no content if there are no recipes", func() {
		fakeService := &mockRecipeLister{
			listRecipes: func(ctx context.Context, userID int64) ([]*services.RecipeSummary, error) {
				return nil, sql.ErrNoRows
			},
		}

		req, err := http.NewRequest(http.MethodGet, "/recipes", nil)
		Expect(err).ToNot(HaveOccurred())

		resp := recipes.ListRecipes(fakeService).Handle(&api.Request{
			Req:    req,
			UserID: 2,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
	})

	It("returns an error if the repository call fails", func() {
		fakeService := &mockRecipeLister{
			listRecipes: func(ctx context.Context, userID int64) ([]*services.RecipeSummary, error) {
				return nil, errors.New("some error")
			},
		}

		req, err := http.NewRequest(http.MethodGet, "/recipes", nil)
		Expect(err).ToNot(HaveOccurred())

		resp := recipes.ListRecipes(fakeService).Handle(&api.Request{
			Req:    req,
			UserID: 2,
		})

		Expect(resp.StatusCode).To(Equal(http.StatusInternalServerError))
	})
})

type mockRecipeLister struct {
	listRecipes func(ctx context.Context, userID int64) ([]*services.RecipeSummary, error)
}

func (m *mockRecipeLister) ListRecipes(ctx context.Context, userID int64) ([]*services.RecipeSummary, error) {
	return m.listRecipes(ctx, userID)
}
