package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"net/mail"

	"golang.org/x/crypto/bcrypt"
)

const BCRYPT_COST = 10

type Credentials struct {
	ID       int64
	Password string
}

type User struct {
	Username string
	Email    string
}

type UsersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (u *UsersRepository) ExistsByUsername(username string) (bool, error) {
	if username == "" {
		fmt.Println("Could not check for user: username required")
		return false, errors.New("could not check for user: username required")
	}

	userRow := u.db.QueryRow(existsByUsernameQuery, username)

	var user User
	err := userRow.Scan(&user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		fmt.Printf("Failed to scan row for user by username '%s': %s\n", username, err.Error())
		return false, errors.New("failed to query for user by username")
	}

	return true, nil
}

func (u *UsersRepository) ExistsByEmail(email string) (bool, error) {
	if email == "" {
		fmt.Println("Could not check for user: username required")
		return false, errors.New("could not check for user: email required")
	}

	userRow := u.db.QueryRow(existsByEmailQuery, email)

	var user User
	err := userRow.Scan(&user.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		fmt.Printf("Failed to scan row for user by email '%s': %s\n", email, err.Error())
		return false, errors.New("failed to query for user by email")
	}

	return true, nil
}

func (u *UsersRepository) Insert(username, email, password string) (int64, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), BCRYPT_COST)
	if err != nil {
		fmt.Printf("Could not hash the user password")
		return -1, errors.New("could not hash password")
	}

	res, err := u.db.Exec(insertUserQuery,
		username,
		email,
		string(hashedPassword),
	)

	if err != nil {
		fmt.Printf("User could not be saved: %s\n", err.Error())
		return -1, errors.New("user could not be added")
	}

	id, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("RecipeResponse was not saved correctly: %s\n", err.Error())
		return 0, errors.New("user was not saved correctly")
	}

	return id, nil
}

func (u *UsersRepository) Verify(login, password string) (bool, int64, error) {
	parser := mail.AddressParser{}
	_, err := parser.Parse(login)

	var result *sql.Row
	if err == nil {
		result = u.db.QueryRow(verifyByEmailQuery, login)
	} else {
		result = u.db.QueryRow(verifyByUsernameQuery, login)
	}

	storedCreds := &Credentials{}
	err = result.Scan(&storedCreds.ID, &storedCreds.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Did not find user by email or username for verification")
			return false, -1, nil
		}

		fmt.Println("Bad scan when verifying user credentials")
		return false, -1, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(password)); err != nil {
		fmt.Println("User credentials comparision failed")
		return false, -1, nil
	}

	return true, storedCreds.ID, nil
}

const existsByUsernameQuery = "SELECT username FROM users WHERE username=?"
const existsByEmailQuery = "SELECT email FROM users WHERE email=?"
const insertUserQuery = `INSERT INTO users
(username,
email,
password_hash)
VALUES (?, ?, ?)
`
const verifyByUsernameQuery = "select id, password_hash from users where username=?"
const verifyByEmailQuery = "select id, password_hash from users where email=?"
