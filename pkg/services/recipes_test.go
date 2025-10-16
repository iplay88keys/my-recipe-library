package services_test

import (
	"context"
	"database/sql"
	"errors"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/iplay88keys/my-recipe-library/pkg/helpers"
	"github.com/iplay88keys/my-recipe-library/pkg/repositories"
	"github.com/iplay88keys/my-recipe-library/pkg/services"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RecipeService", func() {
	var (
		recipeService       *services.RecipeService
		mockRecipesRepo     *MockRecipesRepository
		mockIngredientsRepo *MockIngredientsRepository
		mockStepsRepo       *MockStepsRepository
		db                  *sql.DB
		mock                sqlmock.Sqlmock
		ctx                 context.Context
		userID              int64
		recipeID            int64
	)

	BeforeEach(func() {
		var err error
		db, mock, err = sqlmock.New()
		Expect(err).ToNot(HaveOccurred())

		mockRecipesRepo = &MockRecipesRepository{}
		mockIngredientsRepo = &MockIngredientsRepository{}
		mockStepsRepo = &MockStepsRepository{}
		recipeService = services.NewRecipeService(mockRecipesRepo, mockIngredientsRepo, mockStepsRepo, db)

		ctx = context.Background()
		userID = 1
		recipeID = 123
	})

	Describe("CreateRecipe", func() {
		var recipeInput *services.RecipeInput

		BeforeEach(func() {
			recipeInput = &services.RecipeInput{
				Name:        "Test Recipe",
				Description: "A test recipe",
				Servings:    helpers.IntPointer(4),
				PrepTime:    helpers.StringPointer("15 minutes"),
				CookTime:    helpers.StringPointer("30 minutes"),
				TotalTime:   helpers.StringPointer("45 minutes"),
				Source:      helpers.StringPointer("Test Cookbook"),
				Ingredients: []*services.IngredientInput{
					{
						Name:     "Flour",
						Amount:   "2",
						Unit:     "cups",
						Notes:    "sifted",
						OrderNum: 1,
					},
					{
						Name:     "Sugar",
						Amount:   "1",
						Unit:     "cup",
						Notes:    "",
						OrderNum: 2,
					},
				},
				Steps: []*services.StepInput{
					{
						Instructions: "Mix dry ingredients",
						OrderNum:     1,
						Notes:        helpers.StringPointer("Be gentle"),
					},
					{
						Instructions: "Add wet ingredients",
						OrderNum:     2,
						Notes:        nil,
					},
				},
			}
		})

		Context("when creating a recipe successfully", func() {
			It("creates a recipe with ingredients and steps", func() {
				mock.ExpectBegin()

				mockRecipesRepo.InsertFunc = func(recipe *repositories.Recipe, userID int64) (int64, error) {
					Expect(recipe.Name).To(Equal(helpers.StringPointer("Test Recipe")))
					Expect(recipe.Description).To(Equal(helpers.StringPointer("A test recipe")))
					Expect(recipe.Creator).To(Equal(helpers.StringPointer("User1")))
					Expect(*recipe.Servings).To(Equal(4))
					Expect(*recipe.PrepTime).To(Equal("15 minutes"))
					Expect(*recipe.CookTime).To(Equal("30 minutes"))
					Expect(*recipe.TotalTime).To(Equal("45 minutes"))
					Expect(*recipe.Source).To(Equal("Test Cookbook"))
					Expect(userID).To(Equal(int64(1)))
					return recipeID, nil
				}

				ingredientCallCount := 0
				mockIngredientsRepo.InsertFunc = func(recipeID int64, ingredient *repositories.Ingredient) error {
					ingredientCallCount++
					Expect(recipeID).To(Equal(recipeID))
					if ingredientCallCount == 1 {
						Expect(*ingredient.Ingredient).To(Equal("Flour"))
						Expect(*ingredient.Amount).To(Equal("2"))
						Expect(*ingredient.Measurement).To(Equal("cups"))
						Expect(*ingredient.Preparation).To(Equal("sifted"))
						Expect(*ingredient.IngredientNumber).To(Equal(1))
					} else if ingredientCallCount == 2 {
						Expect(*ingredient.Ingredient).To(Equal("Sugar"))
						Expect(*ingredient.Amount).To(Equal("1"))
						Expect(*ingredient.Measurement).To(Equal("cup"))
						Expect(*ingredient.Preparation).To(Equal(""))
						Expect(*ingredient.IngredientNumber).To(Equal(2))
					}
					return nil
				}

				stepCallCount := 0
				mockStepsRepo.InsertFunc = func(recipeID int64, step *repositories.Step) error {
					stepCallCount++
					Expect(recipeID).To(Equal(recipeID))
					if stepCallCount == 1 {
						Expect(*step.Instructions).To(Equal("Mix dry ingredients"))
						Expect(*step.StepNumber).To(Equal(1))
					} else if stepCallCount == 2 {
						Expect(*step.Instructions).To(Equal("Add wet ingredients"))
						Expect(*step.StepNumber).To(Equal(2))
					}
					return nil
				}

				mock.ExpectCommit()

				resultID, err := recipeService.CreateRecipe(ctx, userID, recipeInput)

				Expect(err).ToNot(HaveOccurred())
				Expect(resultID).To(Equal(recipeID))
				Expect(ingredientCallCount).To(Equal(2))
				Expect(stepCallCount).To(Equal(2))
			})
		})

		Context("when recipe insert fails", func() {
			It("returns an error", func() {
				mock.ExpectBegin()

				mockRecipesRepo.InsertFunc = func(recipe *repositories.Recipe, userID int64) (int64, error) {
					return 0, errors.New("recipe could not be saved")
				}

				mock.ExpectRollback()

				resultID, err := recipeService.CreateRecipe(ctx, userID, recipeInput)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("recipe could not be saved"))
				Expect(resultID).To(Equal(int64(0)))
			})
		})

		Context("when ingredient insert fails", func() {
			It("returns an error", func() {
				mock.ExpectBegin()

				mockRecipesRepo.InsertFunc = func(recipe *repositories.Recipe, userID int64) (int64, error) {
					return recipeID, nil
				}

				mockIngredientsRepo.InsertFunc = func(recipeID int64, ingredient *repositories.Ingredient) error {
					return errors.New("ingredient error")
				}

				mock.ExpectRollback()

				resultID, err := recipeService.CreateRecipe(ctx, userID, recipeInput)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("ingredient error"))
				Expect(resultID).To(Equal(int64(0)))
			})
		})

		Context("when step insert fails", func() {
			It("returns an error", func() {
				mock.ExpectBegin()

				mockRecipesRepo.InsertFunc = func(recipe *repositories.Recipe, userID int64) (int64, error) {
					return recipeID, nil
				}

				mockIngredientsRepo.InsertFunc = func(recipeID int64, ingredient *repositories.Ingredient) error {
					return nil
				}

				mockStepsRepo.InsertFunc = func(recipeID int64, step *repositories.Step) error {
					return errors.New("step error")
				}

				mock.ExpectRollback()

				resultID, err := recipeService.CreateRecipe(ctx, userID, recipeInput)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("step error"))
				Expect(resultID).To(Equal(int64(0)))
			})
		})
	})

	Describe("GetRecipe", func() {
		Context("when recipe exists", func() {
			It("returns the recipe with ingredients and steps", func() {
				mockRecipesRepo.GetFunc = func(id, userID int64) (*repositories.Recipe, error) {
					Expect(id).To(Equal(recipeID))
					Expect(userID).To(Equal(userID))
					return &repositories.Recipe{
						ID:          helpers.Int64Pointer(recipeID),
						Name:        helpers.StringPointer("Test Recipe"),
						Description: helpers.StringPointer("A test recipe"),
						Creator:     helpers.StringPointer("User1"),
						Servings:    helpers.IntPointer(4),
						PrepTime:    helpers.StringPointer("15 minutes"),
						CookTime:    helpers.StringPointer("30 minutes"),
						CoolTime:    helpers.StringPointer("45 minutes"),
						TotalTime:   helpers.StringPointer("45 minutes"),
						Source:      helpers.StringPointer("Test Cookbook"),
					}, nil
				}

				mockIngredientsRepo.GetForRecipeFunc = func(recipeID int64) ([]*repositories.Ingredient, error) {
					Expect(recipeID).To(Equal(recipeID))
					return []*repositories.Ingredient{
						{
							Ingredient:       helpers.StringPointer("Flour"),
							IngredientNumber: helpers.IntPointer(1),
							Amount:           helpers.StringPointer("2"),
							Measurement:      helpers.StringPointer("cups"),
							Preparation:      helpers.StringPointer("sifted"),
						},
						{
							Ingredient:       helpers.StringPointer("Sugar"),
							IngredientNumber: helpers.IntPointer(2),
							Amount:           helpers.StringPointer("1"),
							Measurement:      helpers.StringPointer("cup"),
							Preparation:      helpers.StringPointer(""),
						},
					}, nil
				}

				mockStepsRepo.GetForRecipeFunc = func(recipeID int64) ([]*repositories.Step, error) {
					Expect(recipeID).To(Equal(recipeID))
					return []*repositories.Step{
						{
							StepNumber:   helpers.IntPointer(1),
							Instructions: helpers.StringPointer("Mix dry ingredients"),
						},
						{
							StepNumber:   helpers.IntPointer(2),
							Instructions: helpers.StringPointer("Add wet ingredients"),
						},
					}, nil
				}

				result, err := recipeService.GetRecipe(ctx, recipeID, userID)

				Expect(err).ToNot(HaveOccurred())
				Expect(result).ToNot(BeNil())
				Expect(result.ID).To(Equal(recipeID))
				Expect(result.Name).To(Equal("Test Recipe"))
				Expect(result.Description).To(Equal("A test recipe"))
				Expect(result.Creator).To(Equal("User1"))
				Expect(*result.Servings).To(Equal(4))
				Expect(*result.PrepTime).To(Equal("15 minutes"))
				Expect(*result.CookTime).To(Equal("30 minutes"))
				Expect(*result.CoolTime).To(Equal("45 minutes"))
				Expect(*result.Source).To(Equal("Test Cookbook"))

				Expect(result.Ingredients).To(HaveLen(2))
				Expect(result.Ingredients[0].Name).To(Equal("Flour"))
				Expect(*result.Ingredients[0].Amount).To(Equal("2"))
				Expect(*result.Ingredients[0].Unit).To(Equal("cups"))
				Expect(*result.Ingredients[0].Notes).To(Equal("sifted"))
				Expect(result.Ingredients[0].OrderNum).To(Equal(1))

				Expect(result.Steps).To(HaveLen(2))
				Expect(result.Steps[0].Instructions).To(Equal("Mix dry ingredients"))
				Expect(result.Steps[0].OrderNum).To(Equal(1))
				Expect(result.Steps[0].Notes).To(BeNil())
			})
		})

		Context("when recipe does not exist", func() {
			It("returns an error", func() {
				mockRecipesRepo.GetFunc = func(id, userID int64) (*repositories.Recipe, error) {
					return nil, nil
				}

				result, err := recipeService.GetRecipe(ctx, recipeID, userID)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("recipe not found"))
				Expect(result).To(BeNil())
			})
		})

		Context("when ingredients query fails", func() {
			It("returns an error", func() {
				mockRecipesRepo.GetFunc = func(id, userID int64) (*repositories.Recipe, error) {
					return &repositories.Recipe{
						ID:          helpers.Int64Pointer(recipeID),
						Name:        helpers.StringPointer("Test Recipe"),
						Description: helpers.StringPointer("A test recipe"),
						Creator:     helpers.StringPointer("User1"),
					}, nil
				}

				mockIngredientsRepo.GetForRecipeFunc = func(recipeID int64) ([]*repositories.Ingredient, error) {
					return nil, errors.New("ingredients error")
				}

				result, err := recipeService.GetRecipe(ctx, recipeID, userID)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("ingredients error"))
				Expect(result).To(BeNil())
			})
		})

		Context("when steps query fails", func() {
			It("returns an error", func() {
				mockRecipesRepo.GetFunc = func(id, userID int64) (*repositories.Recipe, error) {
					return &repositories.Recipe{
						ID:          helpers.Int64Pointer(recipeID),
						Name:        helpers.StringPointer("Test Recipe"),
						Description: helpers.StringPointer("A test recipe"),
						Creator:     helpers.StringPointer("User1"),
					}, nil
				}

				mockIngredientsRepo.GetForRecipeFunc = func(recipeID int64) ([]*repositories.Ingredient, error) {
					return []*repositories.Ingredient{}, nil
				}

				mockStepsRepo.GetForRecipeFunc = func(recipeID int64) ([]*repositories.Step, error) {
					return nil, errors.New("steps error")
				}

				result, err := recipeService.GetRecipe(ctx, recipeID, userID)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("steps error"))
				Expect(result).To(BeNil())
			})
		})
	})

	Describe("ListRecipes", func() {
		Context("when recipes exist", func() {
			It("returns a list of recipe summaries", func() {
				mockRecipesRepo.ListFunc = func(userID int64) ([]*repositories.Recipe, error) {
					Expect(userID).To(Equal(userID))
					return []*repositories.Recipe{
						{
							ID:          helpers.Int64Pointer(1),
							Name:        helpers.StringPointer("Recipe 1"),
							Description: helpers.StringPointer("First recipe"),
						},
						{
							ID:          helpers.Int64Pointer(2),
							Name:        helpers.StringPointer("Recipe 2"),
							Description: helpers.StringPointer("Second recipe"),
						},
						{
							ID:          helpers.Int64Pointer(3),
							Name:        helpers.StringPointer("Recipe 3"),
							Description: helpers.StringPointer("Third recipe"),
						},
					}, nil
				}

				result, err := recipeService.ListRecipes(ctx, userID)

				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(HaveLen(3))

				Expect(result[0].ID).To(Equal(int64(1)))
				Expect(result[0].Name).To(Equal("Recipe 1"))
				Expect(result[0].Description).To(Equal("First recipe"))

				Expect(result[1].ID).To(Equal(int64(2)))
				Expect(result[1].Name).To(Equal("Recipe 2"))
				Expect(result[1].Description).To(Equal("Second recipe"))

				Expect(result[2].ID).To(Equal(int64(3)))
				Expect(result[2].Name).To(Equal("Recipe 3"))
				Expect(result[2].Description).To(Equal("Third recipe"))
			})
		})

		Context("when no recipes exist", func() {
			It("returns an empty list", func() {
				mockRecipesRepo.ListFunc = func(userID int64) ([]*repositories.Recipe, error) {
					return []*repositories.Recipe{}, nil
				}

				result, err := recipeService.ListRecipes(ctx, userID)

				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(HaveLen(0))
			})
		})

		Context("when database query fails", func() {
			It("returns an error", func() {
				mockRecipesRepo.ListFunc = func(userID int64) ([]*repositories.Recipe, error) {
					return nil, errors.New("database error")
				}

				result, err := recipeService.ListRecipes(ctx, userID)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("database error"))
				Expect(result).To(BeNil())
			})
		})
	})
})

type MockRecipesRepository struct {
	InsertFunc func(recipe *repositories.Recipe, userID int64) (int64, error)
	GetFunc    func(id, userID int64) (*repositories.Recipe, error)
	ListFunc   func(userID int64) ([]*repositories.Recipe, error)
}

func (m *MockRecipesRepository) Insert(recipe *repositories.Recipe, userID int64) (int64, error) {
	if m.InsertFunc != nil {
		return m.InsertFunc(recipe, userID)
	}
	return 0, nil
}

func (m *MockRecipesRepository) Get(id, userID int64) (*repositories.Recipe, error) {
	if m.GetFunc != nil {
		return m.GetFunc(id, userID)
	}
	return nil, nil
}

func (m *MockRecipesRepository) List(userID int64) ([]*repositories.Recipe, error) {
	if m.ListFunc != nil {
		return m.ListFunc(userID)
	}
	return nil, nil
}

type MockIngredientsRepository struct {
	InsertFunc       func(recipeID int64, ingredient *repositories.Ingredient) error
	GetForRecipeFunc func(recipeID int64) ([]*repositories.Ingredient, error)
}

func (m *MockIngredientsRepository) Insert(recipeID int64, ingredient *repositories.Ingredient) error {
	if m.InsertFunc != nil {
		return m.InsertFunc(recipeID, ingredient)
	}
	return nil
}

func (m *MockIngredientsRepository) GetForRecipe(recipeID int64) ([]*repositories.Ingredient, error) {
	if m.GetForRecipeFunc != nil {
		return m.GetForRecipeFunc(recipeID)
	}
	return nil, nil
}

type MockStepsRepository struct {
	InsertFunc       func(recipeID int64, step *repositories.Step) error
	GetForRecipeFunc func(recipeID int64) ([]*repositories.Step, error)
}

func (m *MockStepsRepository) Insert(recipeID int64, step *repositories.Step) error {
	if m.InsertFunc != nil {
		return m.InsertFunc(recipeID, step)
	}
	return nil
}

func (m *MockStepsRepository) GetForRecipe(recipeID int64) ([]*repositories.Step, error) {
	if m.GetForRecipeFunc != nil {
		return m.GetForRecipeFunc(recipeID)
	}
	return nil, nil
}
