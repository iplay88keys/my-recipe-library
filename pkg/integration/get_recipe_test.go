package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	. "github.com/iplay88keys/my-recipe-library/pkg/helpers"

	"github.com/iplay88keys/my-recipe-library/pkg/api/recipes"
	"github.com/iplay88keys/my-recipe-library/pkg/api/users"
	"github.com/iplay88keys/my-recipe-library/pkg/repositories"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetRecipe", func() {
	var (
		username        string
		password        string
		recipeID        int64
		anotherRecipeID int64
		token           string
	)

	BeforeEach(func() {
		_, err := db.Exec("DELETE FROM users WHERE id IS NOT NULL")
		Expect(err).ToNot(HaveOccurred())

		username = "get_recipe_user"
		password = "Pa3$word123"

		usersRepo := repositories.NewUsersRepository(db)
		userID, err := usersRepo.Insert(username, username+"@example.com", password)
		Expect(err).ToNot(HaveOccurred())

		anotherUserID, err := usersRepo.Insert("get_recipe_different_user", "get_recipe_different_user@example.com", password)
		Expect(err).ToNot(HaveOccurred())

		res, err := db.Exec(`INSERT INTO recipes (
            creator, name, description, servings, prep_time, total_time
            ) VALUES (?, ?, ?, ?, ?, ?)`,
			userID, "Ice Cream", "Yum.", 1, "1 m", "1 m")
		Expect(err).ToNot(HaveOccurred())

		recipeID, err = res.LastInsertId()
		Expect(err).ToNot(HaveOccurred())

		res, err = db.Exec(`INSERT INTO recipes (
            creator, name, description, servings, prep_time, total_time
            ) VALUES (?, ?, ?, ?, ?, ?)`,
			anotherUserID, "Should not have access", "Hidden.", 1, "1 m", "1 m")
		Expect(err).ToNot(HaveOccurred())

		anotherRecipeID, err = res.LastInsertId()
		Expect(err).ToNot(HaveOccurred())

		res, err = db.Exec("INSERT INTO ingredients (name) VALUES (?)",
			"Vanilla Ice Cream")
		Expect(err).ToNot(HaveOccurred())

		ingredientID, err := res.LastInsertId()
		Expect(err).ToNot(HaveOccurred())

		res, err = db.Exec("INSERT INTO measurements (name) VALUES (?)",
			"Scoop")
		Expect(err).ToNot(HaveOccurred())

		measurementID, err := res.LastInsertId()
		Expect(err).ToNot(HaveOccurred())

		_, err = db.Exec(`INSERT INTO recipe_ingredients (
            recipe_id, ingredient_id, ingredient_no, amount, measurement_id
            ) VALUES (?, ?, ?, ?, ?)`,
			recipeID, ingredientID, 1, 1, measurementID)
		Expect(err).ToNot(HaveOccurred())

		_, err = db.Exec(`INSERT INTO recipe_steps (
            recipe_id, step_no, instructions
            ) VALUES (?, ?, ?)`,
			recipeID, 1, "Place ice cream in bowl.")
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
		It("returns a single recipe if owned by that user", func() {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%s/api/v1/recipes/%d", port, recipeID), nil)
			Expect(err).ToNot(HaveOccurred())

			req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))

			resp, err := client.Do(req)
			Expect(err).ToNot(HaveOccurred())

			defer resp.Body.Close()
			bytes, err := io.ReadAll(resp.Body)
			Expect(err).ToNot(HaveOccurred())

			var recipeList recipes.RecipeResponse
			err = json.Unmarshal(bytes, &recipeList)
			Expect(err).ToNot(HaveOccurred())

			Expect(recipeList).To(Equal(recipes.RecipeResponse{
				Recipe: repositories.Recipe{
					ID:          Int64Pointer(recipeID),
					Name:        StringPointer("Ice Cream"),
					Description: StringPointer("Yum."),
					Creator:     StringPointer(username),
					Servings:    IntPointer(1),
					PrepTime:    StringPointer("1 m"),
					CookTime:    nil,
					CoolTime:    nil,
					TotalTime:   StringPointer("1 m"),
					Source:      nil,
				},
				Ingredients: []*repositories.Ingredient{{
					Ingredient:       StringPointer("Vanilla Ice Cream"),
					IngredientNumber: IntPointer(1),
					Amount:           StringPointer("1"),
					Measurement:      StringPointer("Scoop"),
					Preparation:      nil,
				}},
				Steps: []*repositories.Step{{
					StepNumber:   IntPointer(1),
					Instructions: StringPointer("Place ice cream in bowl."),
				}},
			}))
		})

		It("returns not found the recipe is not owned by that user", func() {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%s/api/v1/recipes/%d", port, anotherRecipeID), nil)
			Expect(err).ToNot(HaveOccurred())

			req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))

			resp, err := client.Do(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("returns not found when the recipe is not found", func() {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%s/api/v1/recipes/9999", port), nil)
			Expect(err).ToNot(HaveOccurred())

			req.Header.Set("Authorization", fmt.Sprintf("bearer %s", token))

			resp, err := client.Do(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
		})
	})

	It("returns unauthorized when not authenticated", func() {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://localhost:%s/api/v1/recipes/%d", port, recipeID), nil)
		Expect(err).ToNot(HaveOccurred())

		resp, err := client.Do(req)
		Expect(err).ToNot(HaveOccurred())
		Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
	})
})
