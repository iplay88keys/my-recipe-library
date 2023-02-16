package repositories_test

import (
	"database/sql"
	"errors"

	"github.com/DATA-DOG/go-sqlmock"

	. "github.com/iplay88keys/my-recipe-library/pkg/helpers"
	"github.com/iplay88keys/my-recipe-library/pkg/repositories"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Ingredients Repository", func() {
	var (
		db   *sql.DB
		mock sqlmock.Sqlmock
	)

	BeforeEach(func() {
		var err error
		db, mock, err = sqlmock.New()
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("GetForRecipe", func() {
		It("returns the ingredients for a recipe", func() {
			rows := sqlmock.NewRows([]string{"name", "ingredient_no", "amount", "measurement", "preparation"}).
				AddRow("Vanilla Ice Cream", 1, 1, "Scoop", nil).
				AddRow("Root Beer", 2, nil, nil, nil)

			mock.ExpectQuery(`^SELECT .* FROM recipe_ingredients .* WHERE .*=?`).
				WithArgs(1).
				WillReturnRows(rows)

			repo := repositories.NewIngredientsRepository(db)
			recipes, err := repo.GetForRecipe(1)
			Expect(err).ToNot(HaveOccurred())

			Expect(recipes).To(Equal([]*repositories.Ingredient{{
				Ingredient:       StringPointer("Vanilla Ice Cream"),
				IngredientNumber: IntPointer(1),
				Amount:           StringPointer("1"),
				Measurement:      StringPointer("Scoop"),
				Preparation:      nil,
			}, {
				Ingredient:       StringPointer("Root Beer"),
				IngredientNumber: IntPointer(2),
				Amount:           nil,
				Measurement:      nil,
				Preparation:      nil,
			}}))

			Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
		})

		It("returns an error if the query fails", func() {
			mock.ExpectQuery(`^SELECT .* FROM recipe_ingredients .* WHERE .*=?`).
				WithArgs(1).
				WillReturnError(errors.New("error"))

			repo := repositories.NewIngredientsRepository(db)
			_, err := repo.GetForRecipe(1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to fetch recipe ingredients"))
		})

		It("returns an error if the row cannot be scanned", func() {
			rows := sqlmock.NewRows([]string{"not", "expected", "columns"}).
				AddRow("bad", "values", "returned")

			mock.ExpectQuery(`^SELECT .* FROM recipe_ingredients .* WHERE .*=?`).
				WithArgs(1).
				WillReturnRows(rows)

			repo := repositories.NewIngredientsRepository(db)
			_, err := repo.GetForRecipe(1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to scan recipe ingredients"))
		})

		It("returns an error if the rows cannot all be scanned", func() {
			rows := sqlmock.NewRows([]string{"name", "ingredient_no", "amount", "measurement", "preparation"}).
				AddRow("Vanilla Ice Cream", 1, 1, "Scoop", nil).
				RowError(0, errors.New("some error"))

			mock.ExpectQuery(`^SELECT .* FROM recipe_ingredients .* WHERE .*=?`).
				WithArgs(1).
				WillReturnRows(rows)

			repo := repositories.NewIngredientsRepository(db)
			_, err := repo.GetForRecipe(1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to loop through recipe ingredients"))
		})
	})
})
