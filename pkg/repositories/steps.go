package repositories

import (
	"database/sql"
	"errors"
	"fmt"
)

type Step struct {
	StepNumber   *int    `json:"step_number"`
	Instructions *string `json:"instructions"`
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
		return nil, errors.New(fmt.Sprintf("failed to fetch recipe steps: %s", err.Error()))
	}
	defer rows.Close()

	var recipeSteps []*Step
	for rows.Next() {
		r := &Step{}
		if err := rows.Scan(&r.StepNumber, &r.Instructions); err != nil {
			return nil, errors.New(fmt.Sprintf("failed to scan recipe steps: %s", err.Error()))
		}
		recipeSteps = append(recipeSteps, r)
	}
	if rows.Err() != nil {
		return nil, errors.New(fmt.Sprintf("failed to loop through recipe steps: %s", rows.Err()))
	}

	return recipeSteps, nil
}

const getStepsForRecipeQuery = `SELECT step_no, instructions FROM recipe_steps WHERE recipe_id=?`
