package repositories_test

import (
	"database/sql"
	"errors"

	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/crypto/bcrypt"

	"github.com/iplay88keys/my-recipe-library/pkg/repositories"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Users Repository", func() {
	var (
		db   *sql.DB
		mock sqlmock.Sqlmock
	)

	BeforeEach(func() {
		var err error
		db, mock, err = sqlmock.New()
		Expect(err).ToNot(HaveOccurred())
	})

	Describe("ExistsByUsername", func() {
		It("returns true if a user exists exists with a specified username", func() {
			usernameRow := sqlmock.NewRows([]string{"username"}).
				AddRow("recipeGuru")

			mock.ExpectQuery("^SELECT .+ FROM users WHERE username=?").
				WithArgs("recipeGuru").
				WillReturnRows(usernameRow)

			repo := repositories.NewUsersRepository(db)
			userExists, err := repo.ExistsByUsername("recipeGuru")
			Expect(err).ToNot(HaveOccurred())
			Expect(userExists).To(BeTrue())
		})

		It("returns false if no user is found", func() {
			mock.ExpectQuery("^SELECT .+ FROM users WHERE username=?").
				WithArgs("missing").
				WillReturnError(sql.ErrNoRows)

			repo := repositories.NewUsersRepository(db)

			userExists, err := repo.ExistsByUsername("missing")
			Expect(err).ToNot(HaveOccurred())
			Expect(userExists).To(BeFalse())
		})

		It("returns false and an error if the username is empty", func() {
			repo := repositories.NewUsersRepository(db)

			userExists, err := repo.ExistsByUsername("")
			Expect(err).To(HaveOccurred())
			Expect(userExists).To(BeFalse())
			Expect(err.Error()).To(Equal("could not check for user: username required"))
		})

		It("returns false and an error if an error occurs when querying", func() {
			mock.ExpectQuery("^SELECT .+ FROM users WHERE username=?").
				WithArgs("missing").
				WillReturnError(errors.New("blah"))

			repo := repositories.NewUsersRepository(db)

			userExists, err := repo.ExistsByUsername("missing")
			Expect(err).To(HaveOccurred())
			Expect(userExists).To(BeFalse())
			Expect(err.Error()).To(Equal("failed to query for user by username"))
		})
	})

	Describe("ExistsByEmail", func() {
		It("returns true if a user exists with a specified email", func() {
			emailRow := sqlmock.NewRows([]string{"email"}).
				AddRow("busyCook@example.com")

			mock.ExpectQuery("^SELECT .+ FROM users WHERE email=?").
				WithArgs("busyCook@example.com").
				WillReturnRows(emailRow)

			repo := repositories.NewUsersRepository(db)

			userExists, err := repo.ExistsByEmail("busyCook@example.com")
			Expect(err).ToNot(HaveOccurred())
			Expect(userExists).To(BeTrue())
		})

		It("returns false if no user is found", func() {
			mock.ExpectQuery("^SELECT .+ FROM users WHERE email=?").
				WithArgs("missing@example.com").
				WillReturnError(sql.ErrNoRows)

			repo := repositories.NewUsersRepository(db)

			userExists, err := repo.ExistsByEmail("missing@example.com")
			Expect(err).ToNot(HaveOccurred())
			Expect(userExists).To(BeFalse())
		})

		It("returns false and an error if the email is empty", func() {
			repo := repositories.NewUsersRepository(db)

			userExists, err := repo.ExistsByEmail("")
			Expect(err).To(HaveOccurred())
			Expect(userExists).To(BeFalse())
			Expect(err.Error()).To(Equal("could not check for user: email required"))
		})

		It("returns false and an error if an error occurs when querying", func() {
			mock.ExpectQuery("^SELECT .+ FROM users WHERE email=?").
				WithArgs("missing@example.com").
				WillReturnError(errors.New("error"))

			repo := repositories.NewUsersRepository(db)

			userExists, err := repo.ExistsByEmail("missing@example.com")
			Expect(err).To(HaveOccurred())
			Expect(userExists).To(BeFalse())
			Expect(err.Error()).To(Equal("failed to query for user by email"))
		})
	})

	Describe("Insert", func() {
		It("inserts a user", func() {
			res := sqlmock.NewResult(0, 1)

			mock.ExpectExec("^INSERT INTO users").
				WithArgs(
					"some username",
					"some email",
					sqlmock.AnyArg(),
				).WillReturnResult(res)

			repo := repositories.NewUsersRepository(db)
			id, err := repo.Insert("some username", "some email", "some-password")
			Expect(err).ToNot(HaveOccurred())

			Expect(id).To(BeEquivalentTo(0))

			Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
		})

		It("returns an error if a constraint fails", func() {
			mock.ExpectExec("^INSERT INTO users").
				WillReturnError(errors.New("constraint fails"))

			repo := repositories.NewUsersRepository(db)
			_, err := repo.Insert("some username", "some email", "some-password")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("user could not be added"))
		})

		It("returns an error if the result's LastInsertId fails", func() {
			res := sqlmock.NewErrorResult(errors.New("some error"))

			mock.ExpectExec("^INSERT INTO users").
				WillReturnResult(res)

			repo := repositories.NewUsersRepository(db)
			_, err := repo.Insert("some username", "some email", "some-password")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("user was not saved correctly"))
		})
	})

	Describe("Login", func() {
		It("verifies a user's credentials by username", func() {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte("some-password"), repositories.BCRYPT_COST)
			Expect(err).ToNot(HaveOccurred())

			passwordRow := sqlmock.NewRows([]string{"id", "password"}).
				AddRow("10", hashedPassword)

			mock.ExpectQuery("^select id, password_hash from users where username=?").
				WithArgs("some_username").
				WillReturnRows(passwordRow)

			repo := repositories.NewUsersRepository(db)

			valid, userID, err := repo.Verify("some_username", "some-password")
			Expect(err).ToNot(HaveOccurred())
			Expect(valid).To(BeTrue())
			Expect(userID).To(BeEquivalentTo(10))

			Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
		})

		It("verifies a user's credentials by email", func() {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte("some-password"), repositories.BCRYPT_COST)
			Expect(err).ToNot(HaveOccurred())

			passwordRow := sqlmock.NewRows([]string{"id", "password"}).
				AddRow("20", hashedPassword)

			mock.ExpectQuery("^select id, password_hash from users where email=?").
				WithArgs("someone@example.com").
				WillReturnRows(passwordRow)

			repo := repositories.NewUsersRepository(db)

			valid, userID, err := repo.Verify("someone@example.com", "some-password")
			Expect(err).ToNot(HaveOccurred())
			Expect(valid).To(BeTrue())
			Expect(userID).To(BeEquivalentTo(20))

			Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
		})

		It("returns false if the password is incorrect", func() {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte("some-password"), repositories.BCRYPT_COST)
			Expect(err).ToNot(HaveOccurred())

			passwordRow := sqlmock.NewRows([]string{"id", "password"}).
				AddRow("10", hashedPassword)

			mock.ExpectQuery("^select id, password_hash from users where email=?").
				WithArgs("someone@example.com").
				WillReturnRows(passwordRow)

			repo := repositories.NewUsersRepository(db)

			valid, userID, err := repo.Verify("someone@example.com", "invalid-password")
			Expect(err).ToNot(HaveOccurred())
			Expect(valid).To(BeFalse())
			Expect(userID).To(BeEquivalentTo(-1))

			Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
		})

		It("returns false if the user cannot be found", func() {
			mock.ExpectQuery("^select id, password_hash from users where email=?").
				WithArgs("missing@example.com").
				WillReturnError(sql.ErrNoRows)

			repo := repositories.NewUsersRepository(db)

			valid, userID, err := repo.Verify("missing@example.com", "some-password")
			Expect(err).ToNot(HaveOccurred())
			Expect(valid).To(BeFalse())
			Expect(userID).To(BeEquivalentTo(-1))

			Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
		})

		It("returns an error if there is a problem finding the user cannot be found by username or email", func() {
			mock.ExpectQuery("^select id, password_hash from users where email=?").
				WithArgs("someone@example.com").
				WillReturnError(errors.New("some error"))

			repo := repositories.NewUsersRepository(db)

			valid, userID, err := repo.Verify("someone@example.com", "some-password")
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("some error"))
			Expect(valid).To(BeFalse())
			Expect(userID).To(BeEquivalentTo(-1))

			Expect(mock.ExpectationsWereMet()).ToNot(HaveOccurred())
		})
	})
})
