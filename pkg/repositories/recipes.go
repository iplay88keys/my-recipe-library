package repositories

import (
	"database/sql"
	"errors"
	"fmt"
)

type Recipe struct {
	ID          *int64  `json:"id"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Creator     *string `json:"creator,omitempty"`
	Servings    *int    `json:"servings,omitempty"`
	PrepTime    *string `json:"prep_time,omitempty"`
	CookTime    *string `json:"cook_time,omitempty"`
	CoolTime    *string `json:"cool_time,omitempty"`
	TotalTime   *string `json:"total_time,omitempty"`
	Source      *string `json:"source,omitempty"`
}

type RecipesRepository struct {
	db *sql.DB
}

func NewRecipesRepository(db *sql.DB) *RecipesRepository {
	return &RecipesRepository{db: db}
}

func (r *RecipesRepository) List(userID int64) ([]*Recipe, error) {
	rows, err := r.db.Query(listRecipesQuery, userID)
	if err != nil {
		fmt.Printf("Failed to fetch recipes: %s\n", err.Error())
		return nil, errors.New("failed to fetch recipes")
	}
	defer rows.Close()

	var recipes []*Recipe
	for rows.Next() {
		r := &Recipe{}
		if err := rows.Scan(&r.ID, &r.Name, &r.Description); err != nil {
			fmt.Printf("Failed to scan recipes: %s\n", rows.Err())
			return nil, errors.New("failed to scan recipes")
		}
		recipes = append(recipes, r)
	}
	if rows.Err() != nil {
		fmt.Printf("Failed to loop through recipes: %s\n", rows.Err())
		return nil, errors.New("failed to retrieve recipes")
	}

	return recipes, nil
}

func (r *RecipesRepository) Get(id, userID int64) (*Recipe, error) {
	row := r.db.QueryRow(getRecipeQuery, id, userID)

	recipe := &Recipe{}
	if err := row.Scan(&recipe.ID,
		&recipe.Name,
		&recipe.Description,
		&recipe.Creator,
		&recipe.Servings,
		&recipe.PrepTime,
		&recipe.CookTime,
		&recipe.CoolTime,
		&recipe.TotalTime,
		&recipe.Source); err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}

		fmt.Printf("Failed to scan recipe '%d': %s\n", id, err.Error())
		return nil, errors.New("failed to retrieve recipe")
	}

	return recipe, nil
}

func (r *RecipesRepository) Insert(recipe *Recipe, userID int64) (int64, error) {
	res, err := r.db.Exec(insertRecipeQuery,
		userID,
		recipe.Name,
		recipe.Description,
		recipe.Servings,
		recipe.PrepTime,
		recipe.CookTime,
		recipe.CoolTime,
		recipe.TotalTime,
		recipe.Source,
	)

	if err != nil {
		fmt.Printf("RecipeResponse could not be saved: %s\n", err.Error())
		return 0, errors.New("recipe could not be saved")
	}

	id, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("RecipeResponse was not saved correctly: %s\n", err.Error())
		return 0, errors.New(fmt.Sprintf("recipe was not saved correctly: %s", err.Error()))
	}

	return id, nil
}

const listRecipesQuery = "SELECT id, name, description FROM recipes WHERE creator=?"
const getRecipeQuery = `SELECT
    r.id,
    r.name,
    r.description,
    u.username,
    r.servings,
    r.prep_time,
    r.cook_time,
    r.cool_time,
    r.total_time,
    r.source FROM recipes as r
LEFT JOIN users as u on r.creator=u.id
WHERE r.id=? AND r.creator=?
`
const insertRecipeQuery = `INSERT INTO recipes
    (creator,
    name,
    description,
    servings,
    prep_time,
    cook_time,
    cool_time,
    total_time,
    source)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
`
