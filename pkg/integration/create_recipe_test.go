package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/iplay88keys/my-recipe-library/pkg/api/recipes"
	"github.com/iplay88keys/my-recipe-library/pkg/api/users"
	"github.com/iplay88keys/my-recipe-library/pkg/repositories"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CreateRecipe", func() {
	var (
		userID int64
		token  string
	)

	BeforeEach(func() {
		_, err := db.Exec("DELETE FROM users WHERE id IS NOT NULL")
		Expect(err).ToNot(HaveOccurred())

		username := "create_recipe_user"
		password := "Pa3$word123"

		usersRepo := repositories.NewUsersRepository(db)
		userID, err = usersRepo.Insert(username, username+"@example.com", password)
		Expect(err).ToNot(HaveOccurred())

		reqBody := []byte(fmt.Sprintf(`{
            "login": "%s",
            "password": "%s"
        }`, username, password))

		resp, err := http.Post(fmt.Sprintf("http://localhost:%s/api/v1/users/login", port), "application/json", bytes.NewBuffer(reqBody))
		Expect(err).ToNot(HaveOccurred())

		Expect(resp.StatusCode).To(Equal(200))

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		Expect(err).ToNot(HaveOccurred())

		var loginResponse users.UserLoginResponse
		err = json.Unmarshal(body, &loginResponse)
		Expect(err).ToNot(HaveOccurred())

		token = loginResponse.AccessToken
	})

	Context("authenticated", func() {
		It("creates a recipe", func() {
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

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%s/api/v1/recipes", port), bytes.NewBuffer(body))
			Expect(err).ToNot(HaveOccurred())

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))

			resp, err := client.Do(req)
			Expect(err).ToNot(HaveOccurred())

			defer resp.Body.Close()
			bytes, err := io.ReadAll(resp.Body)
			Expect(err).ToNot(HaveOccurred())

			var response recipes.CreateRecipeResponse
			err = json.Unmarshal(bytes, &response)
			Expect(err).ToNot(HaveOccurred())

			row := db.QueryRow("SELECT id FROM recipes WHERE creator=? LIMIT 1", userID)

			var recipeID int64
			err = row.Scan(&recipeID)
			Expect(err).ToNot(HaveOccurred())

			Expect(response).To(Equal(recipes.CreateRecipeResponse{
				RecipeID: recipeID,
			}))
		})
	})

	It("returns unauthorized when not authenticated", func() {
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

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%s/api/v1/recipes", port), bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
	})
})
