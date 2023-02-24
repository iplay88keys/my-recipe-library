package integration_test

import (
    _ "github.com/go-sql-driver/mysql"

    . "github.com/onsi/ginkgo"
    //. "github.com/onsi/gomega"
)

var _ = Describe("GetCookbook", func() {
    //var (
    //  db         *sql.DB
    //  testDBName string
    //)
    //
    //BeforeEach(func() {
    //  testDBName = time.Now().Format(time.Stamp)
    //
    //  var err error
    //  db, err := sql.Open("mysql", strings.TrimSpace(databaseURL))
    //  if err != nil {
    //      panic(err)
    //  }
    //
    //  _, err = db.Exec("CREATE DATABASE ?", testDBName)
    //  Expect(err).ToNot(HaveOccurred())
    //
    //  userRepo := repositories.NewUsersRepository(db)
    //  user, err := userRepo.Create(&repositories.User{
    //
    //  })
    //
    //  recipeRepo := repositories.NewRecipesRepository(db)
    //  recipe, err := recipeRepo.Create(&repositories.RecipeResponse{
    //      Name:        StringPointer("RecipeResponse Name"),
    //      Description: StringPointer(""),
    //      Creator:     StringPointer(),
    //      Servings:    StringPointer(),
    //      PrepTime:    StringPointer(),
    //      CookTime:    StringPointer(),
    //      CoolTime:    StringPointer(),
    //      TotalTime:   StringPointer(),
    //      Source:      StringPointer(),
    //  })
    //  Expect(err).ToNot(HaveOccurred())
    //
    //  cookbookRepo := repositories.NewCookbookRepository(db)
    //  cookbook, err := recipeRepo.Create(&repositories.Cookbook{
    //
    //  })
    //  Expect(err).ToNot(HaveOccurred())
    //})
    //
    //AfterEach(func() {
    //  _, err := db.Exec("DROP DATABASE ?", testDBName)
    //  Expect(err).ToNot(HaveOccurred())
    //
    //  err = db.Close()
    //  Expect(err).ToNot(HaveOccurred())
    //})
    //
    //It("returns a list of recipes for a cookbook", func() {
    //  resp, err := client.Get(fmt.Sprintf("http://localhost:%s/api/v1/recipes/1", port))
    //  Expect(err).ToNot(HaveOccurred())
    //
    //  defer resp.Body.Close()
    //  bytes, err := ioutil.ReadAll(resp.Body)
    //  Expect(err).ToNot(HaveOccurred())
    //
    //  var recipeList recipes.RecipeResponse
    //  err = json.Unmarshal(bytes, &recipeList)
    //  Expect(err).ToNot(HaveOccurred())
    //
    //  Expect(recipeList).To(Equal(recipes.RecipeResponse{
    //      RecipeResponse: repositories.RecipeResponse{
    //          ID:          IntPointer(1),
    //          Name:        StringPointer("Root Beer Float"),
    //          Description: StringPointer("Delicious drink for a hot summer day."),
    //          Creator:     StringPointer("user"),
    //          Servings:    IntPointer(1),
    //          PrepTime:    StringPointer("5 m"),
    //          CookTime:    nil,
    //          CoolTime:    nil,
    //          TotalTime:   StringPointer("5 m"),
    //          Source:      nil,
    //      },
    //      Ingredients: []*repositories.Ingredient{{
    //          Ingredient:       StringPointer("Vanilla Ice Cream"),
    //          IngredientNumber: IntPointer(1),
    //          Amount:           StringPointer("1"),
    //          Measurement:      StringPointer("Scoop"),
    //          Preparation:      nil,
    //      }, {
    //          Ingredient:       StringPointer("Root Beer"),
    //          IngredientNumber: IntPointer(2),
    //          Amount:           nil,
    //          Measurement:      nil,
    //          Preparation:      nil,
    //      }},
    //      Steps: []*repositories.Step{{
    //          StepNumber:   IntPointer(1),
    //          Instructions: StringPointer("Place ice cream in glass."),
    //      }, {
    //          StepNumber:   IntPointer(2),
    //          Instructions: StringPointer("Top with Root Beer."),
    //      }},
    //  }))
    //})
})
