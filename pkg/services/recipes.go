package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/iplay88keys/my-recipe-library/pkg/repositories"
)

func stringPtr(s string) *string {
	return &s
}

type RecipesRepositoryInterface interface {
	Insert(recipe *repositories.Recipe, userID int64) (int64, error)
	Get(id, userID int64) (*repositories.Recipe, error)
	List(userID int64) ([]*repositories.Recipe, error)
}

type IngredientsRepositoryInterface interface {
	Insert(recipeID int64, ingredient *repositories.Ingredient) error
	GetForRecipe(recipeID int64) ([]*repositories.Ingredient, error)
}

type StepsRepositoryInterface interface {
	Insert(recipeID int64, step *repositories.Step) error
	GetForRecipe(recipeID int64) ([]*repositories.Step, error)
}

type RecipeService struct {
	recipesRepo     RecipesRepositoryInterface
	ingredientsRepo IngredientsRepositoryInterface
	stepsRepo       StepsRepositoryInterface
	db              *sql.DB
}

func NewRecipeService(
	recipesRepo RecipesRepositoryInterface,
	ingredientsRepo IngredientsRepositoryInterface,
	stepsRepo StepsRepositoryInterface,
	db *sql.DB,
) *RecipeService {
	return &RecipeService{
		recipesRepo:     recipesRepo,
		ingredientsRepo: ingredientsRepo,
		stepsRepo:       stepsRepo,
		db:              db,
	}
}

type RecipeInput struct {
	Name        string
	Description string
	Servings    *int
	PrepTime    *string
	CookTime    *string
	CoolTime    *string
	TotalTime   *string
	Source      *string
	Ingredients []*IngredientInput
	Steps       []*StepInput
}

type IngredientInput struct {
	Name     string
	Amount   string
	Unit     string
	Notes    string
	OrderNum int
}

type StepInput struct {
	Instructions string
	OrderNum     int
	Notes        *string
}

type RecipeDetail struct {
	ID          int64
	Name        string
	Description string
	Creator     string
	Servings    *int
	PrepTime    *string
	CookTime    *string
	CoolTime    *string
	TotalTime   *string
	Source      *string
	Ingredients []*IngredientDetail
	Steps       []*StepDetail
}

type IngredientDetail struct {
	Name     string
	Amount   *string
	Unit     *string
	Notes    *string
	OrderNum int
}

type StepDetail struct {
	Instructions string
	OrderNum     int
	Notes        *string
}

type RecipeSummary struct {
	ID          int64
	Name        string
	Description string
}

func (s *RecipeService) CreateRecipe(ctx context.Context, userID int64, recipe *RecipeInput) (int64, error) {
	return s.runInTransaction(ctx, func(tx *sql.Tx) (int64, error) {
		recipeID, err := s.recipesRepo.Insert(&repositories.Recipe{
			Name:        &recipe.Name,
			Description: &recipe.Description,
			Creator:     stringPtr(fmt.Sprintf("User%d", userID)),
			Servings:    recipe.Servings,
			PrepTime:    recipe.PrepTime,
			CookTime:    recipe.CookTime,
			CoolTime:    recipe.CoolTime,
			TotalTime:   recipe.TotalTime,
			Source:      recipe.Source,
		}, userID)
		if err != nil {
			return 0, err
		}

		for _, ingredient := range recipe.Ingredients {
			err = s.ingredientsRepo.Insert(recipeID, &repositories.Ingredient{
				Ingredient:       &ingredient.Name,
				IngredientNumber: &ingredient.OrderNum,
				Amount:           &ingredient.Amount,
				Measurement:      &ingredient.Unit,
				Preparation:      &ingredient.Notes,
			})
			if err != nil {
				return 0, err
			}
		}

		for _, step := range recipe.Steps {
			err = s.stepsRepo.Insert(recipeID, &repositories.Step{
				StepNumber:   &step.OrderNum,
				Instructions: &step.Instructions,
			})
			if err != nil {
				return 0, err
			}
		}

		return recipeID, nil
	})
}

func (s *RecipeService) GetRecipe(ctx context.Context, recipeID, userID int64) (*RecipeDetail, error) {
	recipe, err := s.recipesRepo.Get(recipeID, userID)
	if err != nil {
		return nil, err
	}
	if recipe == nil {
		return nil, errors.New("recipe not found")
	}

	ingredients, err := s.ingredientsRepo.GetForRecipe(recipeID)
	if err != nil {
		return nil, err
	}

	steps, err := s.stepsRepo.GetForRecipe(recipeID)
	if err != nil {
		return nil, err
	}

	recipeDetail := &RecipeDetail{
		ID:          *recipe.ID,
		Name:        *recipe.Name,
		Description: *recipe.Description,
		Creator:     *recipe.Creator,
		Servings:    recipe.Servings,
		PrepTime:    recipe.PrepTime,
		CookTime:    recipe.CookTime,
		CoolTime:    recipe.CoolTime,
		TotalTime:   recipe.TotalTime,
		Source:      recipe.Source,
		Ingredients: make([]*IngredientDetail, len(ingredients)),
		Steps:       make([]*StepDetail, len(steps)),
	}

	for i, ingredient := range ingredients {
		recipeDetail.Ingredients[i] = &IngredientDetail{
			Name:     *ingredient.Ingredient,
			Amount:   ingredient.Amount,
			Unit:     ingredient.Measurement,
			Notes:    ingredient.Preparation,
			OrderNum: *ingredient.IngredientNumber,
		}
	}

	for i, step := range steps {
		recipeDetail.Steps[i] = &StepDetail{
			Instructions: *step.Instructions,
			OrderNum:     *step.StepNumber,
			Notes:        nil,
		}
	}

	return recipeDetail, nil
}

func (s *RecipeService) ListRecipes(ctx context.Context, userID int64) ([]*RecipeSummary, error) {
	recipes, err := s.recipesRepo.List(userID)
	if err != nil {
		return nil, err
	}

	summaries := make([]*RecipeSummary, len(recipes))
	for i, recipe := range recipes {
		summaries[i] = &RecipeSummary{
			ID:          *recipe.ID,
			Name:        *recipe.Name,
			Description: *recipe.Description,
		}
	}

	return summaries, nil
}

func (s *RecipeService) runInTransaction(ctx context.Context, fn func(*sql.Tx) (int64, error)) (int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	return fn(tx)
}
