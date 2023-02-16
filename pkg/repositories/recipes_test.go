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

var _ = Describe("RecipeResponse Repository", func() {
	var (
		db   *sql.DB
		mock sqlmock.Sqlmock
	)

	BeforeEach(func() {
		var err error
		db, mock, err = sqlmock.New()
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("List", func() {
		It("returns the list of all recipes", func() {
			rows := sqlmock.NewRows([]string{"id", "name", "description"}).
				AddRow(0, "First RecipeResponse", "The First").
				AddRow(1, "Second RecipeResponse", "The Second")

			mock.ExpectQuery("^SELECT .+ FROM recipes WHERE .+=?$").
				WithArgs(10).
				WillReturnRows(rows)

			repo := repositories.NewRecipesRepository(db)
			recipes, err := repo.List(10)
			Expect(err).ToNot(HaveOccurred())

			Expect(recipes).To(Equal([]*repositories.Recipe{{
				ID:          Int64Pointer(0),
				Name:        StringPointer("First RecipeResponse"),
				Description: StringPointer("The First"),
			}, {
				ID:          Int64Pointer(1),
				Name:        StringPointer("Second RecipeResponse"),
				Description: StringPointer("The Second"),
			}}))

			Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
		})

		It("returns an error if no recipes are found", func() {
			mock.ExpectQuery("^SELECT .+ FROM recipes WHERE .+=?$").
				WithArgs(20).
				WillReturnError(sql.ErrNoRows)

			repo := repositories.NewRecipesRepository(db)
			_, err := repo.List(20)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to fetch recipes"))
		})

		It("returns an error if the row cannot be scanned", func() {
			rows := sqlmock.NewRows([]string{"not", "expected", "columns"}).
				AddRow("bad", "values", "returned")

			mock.ExpectQuery("^SELECT .+ FROM recipes WHERE .+=?$").
				WithArgs(30).
				WillReturnRows(rows)

			repo := repositories.NewRecipesRepository(db)
			_, err := repo.List(30)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to scan recipes"))
		})

		It("returns an error if the rows cannot all be scanned", func() {
			rows := sqlmock.NewRows([]string{"id", "name", "description"}).
				AddRow(0, "First RecipeResponse", "The First").
				RowError(0, errors.New("some error"))

			mock.ExpectQuery("^SELECT .+ FROM recipes WHERE .+=?$").
				WithArgs(40).
				WillReturnRows(rows)

			repo := repositories.NewRecipesRepository(db)
			_, err := repo.List(40)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to retrieve recipes"))
		})
	})

	Describe("Get", func() {
		It("returns a recipe by its id", func() {
			rows := sqlmock.NewRows([]string{
				"id",
				"name",
				"description",
				"creator",
				"servings",
				"prep_time",
				"cook_time",
				"cool_time",
				"total_time",
				"source",
			}).AddRow(
				1,
				"RecipeResponse Name",
				"RecipeResponse Description",
				"Some Creator",
				3,
				"10 m",
				"30 m",
				"5 m",
				"45 m",
				"Some Book",
			)

			mock.ExpectQuery("^SELECT .+ FROM recipes .+ WHERE .+=? AND .+=?$").
				WithArgs(1, 2).
				WillReturnRows(rows)

			repo := repositories.NewRecipesRepository(db)
			recipe, err := repo.Get(1, 2)
			Expect(err).ToNot(HaveOccurred())

			Expect(recipe).To(Equal(&repositories.Recipe{
				ID:          Int64Pointer(1),
				Name:        StringPointer("RecipeResponse Name"),
				Description: StringPointer("RecipeResponse Description"),
				Creator:     StringPointer("Some Creator"),
				Servings:    IntPointer(3),
				PrepTime:    StringPointer("10 m"),
				CookTime:    StringPointer("30 m"),
				CoolTime:    StringPointer("5 m"),
				TotalTime:   StringPointer("45 m"),
				Source:      StringPointer("Some Book"),
			}))

			Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
		})

		It("returns an error if the recipe cannot be found", func() {
			mock.ExpectQuery("^SELECT .+ FROM recipes .+ WHERE .+=? AND .+=?$").
				WithArgs(0, 0).
				WillReturnError(sql.ErrNoRows)

			repo := repositories.NewRecipesRepository(db)
			_, err := repo.Get(0, 0)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(sql.ErrNoRows))
		})

		It("returns an error if the row cannot be scanned", func() {
			rows := sqlmock.NewRows([]string{"not", "expected", "columns"}).
				AddRow("bad", "values", "returned")

			mock.ExpectQuery("^SELECT .+ FROM recipes .+ WHERE .+=? AND .+=?$").
				WithArgs(0, 0).
				WillReturnRows(rows)

			repo := repositories.NewRecipesRepository(db)
			_, err := repo.Get(0, 0)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("failed to retrieve recipe"))
		})
	})

	Describe("Insert", func() {
		It("inserts a recipe", func() {
			res := sqlmock.NewResult(0, 1)

			mock.ExpectExec("^INSERT INTO recipes").
				WithArgs(
					1,
					"RecipeResponse Name",
					"RecipeResponse Description",
					3,
					"1 hr",
					"2 m",
					"3 m",
					"1 hr 5 m",
					"some website",
				).WillReturnResult(res)

			repo := repositories.NewRecipesRepository(db)
			id, err := repo.Insert(&repositories.Recipe{
				Name:        StringPointer("RecipeResponse Name"),
				Description: StringPointer("RecipeResponse Description"),
				Servings:    IntPointer(3),
				PrepTime:    StringPointer("1 hr"),
				CookTime:    StringPointer("2 m"),
				CoolTime:    StringPointer("3 m"),
				TotalTime:   StringPointer("1 hr 5 m"),
				Source:      StringPointer("some website"),
			}, 1)
			Expect(err).ToNot(HaveOccurred())

			Expect(id).To(BeEquivalentTo(0))

			Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
		})

		It("inserts a recipe even if some values are missing", func() {
			res := sqlmock.NewResult(0, 1)

			mock.ExpectExec("^INSERT INTO recipes").
				WithArgs(
					1,
					"RecipeResponse Name",
					"RecipeResponse Description",
					3,
					nil,
					nil,
					nil,
					nil,
					nil,
				).WillReturnResult(res)

			repo := repositories.NewRecipesRepository(db)
			id, err := repo.Insert(&repositories.Recipe{
				Name:        StringPointer("RecipeResponse Name"),
				Description: StringPointer("RecipeResponse Description"),
				Servings:    IntPointer(3),
			}, 1)
			Expect(err).ToNot(HaveOccurred())

			Expect(id).To(BeEquivalentTo(0))

			Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
		})

		It("returns an error if the recipe cannot be inserted into the database", func() {
			mock.ExpectExec("^INSERT INTO recipes").
				WillReturnError(errors.New("constraint fails"))

			repo := repositories.NewRecipesRepository(db)
			_, err := repo.Insert(&repositories.Recipe{}, 1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("recipe could not be saved"))
		})

		It("returns an error if the result's LastInsertId fails", func() {
			res := sqlmock.NewErrorResult(errors.New("some error"))

			mock.ExpectExec("^INSERT INTO recipes").
				WillReturnResult(res)

			repo := repositories.NewRecipesRepository(db)
			_, err := repo.Insert(&repositories.Recipe{}, 1)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("recipe was not saved correctly"))
		})
	})
})
