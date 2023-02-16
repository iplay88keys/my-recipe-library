package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/iplay88keys/my-recipe-library/pkg/api/users"

	"github.com/iplay88keys/my-recipe-library/pkg/api/recipes"
	. "github.com/iplay88keys/my-recipe-library/pkg/helpers"
	"github.com/iplay88keys/my-recipe-library/pkg/repositories"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ListRecipes", func() {
	var (
		username      string
		password      string
		firstRecipeID int64
		token         string
	)

	BeforeEach(func() {
		_, err := db.Exec("DELETE FROM users WHERE id IS NOT NULL")
		Expect(err).ToNot(HaveOccurred())

		username = "list_recipes_user"
		password = "Pa3$word123"

		usersRepo := repositories.NewUsersRepository(db)
		userID, err := usersRepo.Insert(username, username+"@example.com", password)
		Expect(err).ToNot(HaveOccurred())

		res, err := db.Exec(`INSERT INTO recipes (
            creator, name, description, servings, prep_time, total_time
            ) VALUES (?, ?, ?, ?, ?, ?)`,
			userID, "Root Beer Float", "Delicious drink for a hot summer day.", 1, "5 m", "5 m")
		Expect(err).ToNot(HaveOccurred())

		firstRecipeID, err = res.LastInsertId()
		Expect(err).ToNot(HaveOccurred())

		userID, err = usersRepo.Insert("another_"+username, "another_"+username+"@example.com", password)
		Expect(err).ToNot(HaveOccurred())

		_, err = db.Exec(`INSERT INTO recipes (
            creator, name, description, servings, prep_time, cook_time, total_time
            ) VALUES (?, ?, ?, ?, ?, ?, ?)`,
			userID, "Nana's Beans", "Spruced up baked beans.", 8, "10 m", "1-2 hrs", "1-2 hrs")
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
		It("returns a list of recipes for the user", func() {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%s/api/v1/recipes", port), nil)
			Expect(err).ToNot(HaveOccurred())

			req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))

			resp, err := client.Do(req)
			Expect(err).ToNot(HaveOccurred())

			defer resp.Body.Close()
			bytes, err := io.ReadAll(resp.Body)
			Expect(err).ToNot(HaveOccurred())

			var recipeList recipes.RecipeListResponse
			err = json.Unmarshal(bytes, &recipeList)
			Expect(err).ToNot(HaveOccurred())

			Expect(recipeList).To(Equal(recipes.RecipeListResponse{
				Recipes: []*repositories.Recipe{{
					ID:          Int64Pointer(firstRecipeID),
					Name:        StringPointer("Root Beer Float"),
					Description: StringPointer("Delicious drink for a hot summer day."),
					Creator:     nil,
					PrepTime:    nil,
					Source:      nil,
				}},
			}))
		})
	})

	It("returns unauthorized when not authenticated", func() {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%s/api/v1/recipes", port), nil)
		Expect(err).ToNot(HaveOccurred())

		resp, err := client.Do(req)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
	})
})
