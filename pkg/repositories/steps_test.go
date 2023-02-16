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

var _ = Describe("Steps Repository", func() {
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
		It("returns the list of steps for a recipe", func() {
			rows := sqlmock.NewRows([]string{"step_no", "instructions"}).
				AddRow(1, "Place ice cream in glass.").
				AddRow(2, "Top with Root Beer.")

			mock.ExpectQuery("^SELECT .+ FROM recipe_steps WHERE recipe_id=?").
				WithArgs(1).
				WillReturnRows(rows)

			repo := repositories.NewStepsRepository(db)
			recipes, err := repo.GetForRecipe(1)
			Expect(err).ToNot(HaveOccurred())

			Expect(recipes).To(Equal([]*repositories.Step{{
				StepNumber:   IntPointer(1),
				Instructions: StringPointer("Place ice cream in glass."),
			}, {
				StepNumber:   IntPointer(2),
				Instructions: StringPointer("Top with Root Beer."),
			}}))

			Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
		})

		It("returns an error if the query fails", func() {
			mock.ExpectQuery("^SELECT .+ FROM recipe_steps WHERE recipe_id=?").
				WithArgs(1).
				WillReturnError(errors.New("error"))

			repo := repositories.NewStepsRepository(db)
			_, err := repo.GetForRecipe(1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to fetch recipe steps"))
		})

		It("returns an error if the row cannot be scanned", func() {
			rows := sqlmock.NewRows([]string{"not", "expected", "columns"}).
				AddRow("bad", "values", "returned")

			mock.ExpectQuery("^SELECT .+ FROM recipe_steps WHERE recipe_id=?").
				WithArgs(1).
				WillReturnRows(rows)

			repo := repositories.NewStepsRepository(db)
			_, err := repo.GetForRecipe(1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to scan recipe steps"))
		})

		It("returns an error if the rows cannot all be scanned", func() {
			rows := sqlmock.NewRows([]string{"name", "ingredient_no", "amount", "measurement", "preparation"}).
				AddRow("Vanilla Ice Cream", 0, 1, "Scoop", nil).
				RowError(0, errors.New("some error"))

			mock.ExpectQuery("^SELECT .+ FROM recipe_steps WHERE recipe_id=?").
				WithArgs(1).
				WillReturnRows(rows)

			repo := repositories.NewStepsRepository(db)
			_, err := repo.GetForRecipe(1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to loop through recipe steps"))
		})
	})
})
