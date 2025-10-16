package repositories

import (
	"database/sql"
	"errors"
	"fmt"
)

type Step struct {
	StepNumber   *int
	Instructions *string
}

type StepsRepository struct {
	db *sql.DB
}

func NewStepsRepository(db *sql.DB) *StepsRepository {
	return &StepsRepository{db: db}
}

func (r *StepsRepository) GetForRecipe(recipeID int64) ([]*Step, error) {
	rows, err := r.db.Query(getStepsForRecipeQuery, recipeID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch recipe steps: %s", err.Error())
	}
	defer rows.Close()

	var recipeSteps []*Step
	for rows.Next() {
		r := &Step{}
		if err := rows.Scan(&r.StepNumber, &r.Instructions); err != nil {
			return nil, fmt.Errorf("failed to scan recipe steps: %s", err.Error())
		}
		recipeSteps = append(recipeSteps, r)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to loop through recipe steps: %s", rows.Err())
	}

	return recipeSteps, nil
}

func (r *StepsRepository) Insert(recipeID int64, step *Step) error {
	_, err := r.db.Exec(insertRecipeStepQuery,
		recipeID,
		step.StepNumber,
		step.Instructions,
	)

	if err != nil {
		fmt.Printf("Recipe step could not be saved: %s\n", err.Error())
		return errors.New("recipe step could not be saved")
	}

	return nil
}

const getStepsForRecipeQuery = `SELECT step_no, instructions FROM recipe_steps WHERE recipe_id=?`

const insertRecipeStepQuery = `
  INSERT INTO recipe_steps (recipe_id, step_no, instructions)
  VALUES (?, ?, ?)
`
