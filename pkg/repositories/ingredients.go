package repositories

import (
	"database/sql"
	"errors"
	"fmt"
)

type Ingredient struct {
	ID               *int64
	Ingredient       *string
	IngredientNumber *int
	Amount           *string
	Measurement      *string
	Preparation      *string
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
		return nil, fmt.Errorf("failed to fetch recipe ingredients: %s", err.Error())
	}
	defer rows.Close()

	var recipeIngredients []*Ingredient
	for rows.Next() {
		ingredient := &Ingredient{}
		if err := rows.Scan(&ingredient.Ingredient, &ingredient.IngredientNumber, &ingredient.Amount, &ingredient.Measurement, &ingredient.Preparation); err != nil {
			return nil, fmt.Errorf("failed to scan recipe ingredients: %s", err.Error())
		}
		recipeIngredients = append(recipeIngredients, ingredient)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to loop through recipe ingredients: %s", rows.Err())
	}

	return recipeIngredients, nil
}

func (r *IngredientsRepository) Insert(recipeID int64, ingredient *Ingredient) error {
	ingredientID, err := r.getOrCreateIngredient(*ingredient.Ingredient)
	if err != nil {
		return err
	}

	var measurementID *int64
	if ingredient.Measurement != nil && *ingredient.Measurement != "" {
		measurementID, err = r.getOrCreateMeasurement(*ingredient.Measurement)
		if err != nil {
			return err
		}
	}

	_, err = r.db.Exec(insertRecipeIngredientQuery,
		recipeID,
		ingredientID,
		ingredient.IngredientNumber,
		ingredient.Amount,
		measurementID,
		ingredient.Preparation,
	)

	if err != nil {
		fmt.Printf("Recipe ingredient could not be saved: %s\n", err.Error())
		return errors.New("recipe ingredient could not be saved")
	}

	return nil
}

func (r *IngredientsRepository) getOrCreateIngredient(name string) (int64, error) {
	var id int64
	err := r.db.QueryRow("SELECT id FROM ingredients WHERE name = ?", name).Scan(&id)
	if err == nil {
		return id, nil
	}

	res, err := r.db.Exec("INSERT INTO ingredients (name) VALUES (?)", name)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (r *IngredientsRepository) getOrCreateMeasurement(name string) (*int64, error) {
	var id int64
	err := r.db.QueryRow("SELECT id FROM measurements WHERE name = ?", name).Scan(&id)
	if err == nil {
		return &id, nil
	}

	res, err := r.db.Exec("INSERT INTO measurements (name) VALUES (?)", name)
	if err != nil {
		return nil, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &id, nil
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
const insertRecipeIngredientQuery = `
  INSERT INTO recipe_ingredients (recipe_id, ingredient_id, ingredient_no, amount, measurement_id, preparation)
  VALUES (?, ?, ?, ?, ?, ?)
`
