package repositories

import (
	"database/sql"
	"errors"
	"fmt"
)

type Ingredient struct {
	Ingredient       *string `json:"ingredient"`
	IngredientNumber *int    `json:"ingredient_number"`
	Amount           *string `json:"amount"`
	Measurement      *string `json:"measurement"`
	Preparation      *string `json:"preparation"`
}

type IngredientsRepository struct {
	db *sql.DB
}

func NewIngredientsRepository(db *sql.DB) *IngredientsRepository {
	return &IngredientsRepository{db: db}
}

func (r *IngredientsRepository) GetForRecipe(recipeID int64) ([]*Ingredient, error) {
	rows, err := r.db.Query(getIngredientsForRecipeQuery, recipeID)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to fetch recipe ingredients: %s", err.Error()))
	}
	defer rows.Close()

	var recipeIngredients []*Ingredient
	for rows.Next() {
		r := &Ingredient{}
		if err := rows.Scan(&r.Ingredient, &r.IngredientNumber, &r.Amount, &r.Measurement, &r.Preparation); err != nil {
			return nil, errors.New(fmt.Sprintf("failed to scan recipe ingredients: %s", err.Error()))
		}
		recipeIngredients = append(recipeIngredients, r)
	}
	if rows.Err() != nil {
		return nil, errors.New(fmt.Sprintf("failed to loop through recipe ingredients: %s", rows.Err()))
	}

	return recipeIngredients, nil
}

const getIngredientsForRecipeQuery = `
  SELECT i.name,
         ri.ingredient_no,
         ri.amount,
         m.name,
         ri.preparation 
  FROM recipe_ingredients as ri
  LEFT JOIN ingredients as i on ri.ingredient_id=i.id
  LEFT JOIN measurements as m on ri.measurement_id=m.id
  WHERE ri.recipe_id=?
`
