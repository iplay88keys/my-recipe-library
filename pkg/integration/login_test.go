package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/iplay88keys/my-recipe-library/pkg/repositories"

	"github.com/iplay88keys/my-recipe-library/pkg/api/users"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("login", func() {
	BeforeEach(func() {
		_, err := db.Exec("DELETE FROM users WHERE id IS NOT NULL")
		Expect(err).ToNot(HaveOccurred())
	})

	It("logs in with valid credentials", func() {
		username := "authorized_user"
		password := "Pa3$word123"

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), repositories.BCRYPT_COST)
		Expect(err).ToNot(HaveOccurred())

		_, err = db.Exec("INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)", username, "authorized_user@example.com", string(hashedPassword))
		Expect(err).ToNot(HaveOccurred())

		reqBody := []byte(fmt.Sprintf(`{
            "login": "%s",
            "password": "%s"
        }`, username, password))

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%s/api/v1/users/login", port), bytes.NewBuffer(reqBody))
		Expect(err).ToNot(HaveOccurred())

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		Expect(err).ToNot(HaveOccurred())

		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		Expect(err).ToNot(HaveOccurred())

		var output users.UserLoginResponse
		err = json.Unmarshal(body, &output)
		Expect(err).ToNot(HaveOccurred())

		Expect(output.AccessToken).ToNot(Equal(""))
		Expect(output.RefreshToken).ToNot(Equal(""))
	})

	It("returns unauthorized if the credentials are incorrect", func() {
		username := "unauthorized_user"
		_, err := db.Exec("INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)", username, "unauthorized_user@example.com", "Pa3$word123")
		Expect(err).ToNot(HaveOccurred())

		reqBody := []byte(fmt.Sprintf(`{
            "login": "%s",
            "password": "%s"
        }`, username, "bad-password"))

		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%s/api/v1/users/login", port), bytes.NewBuffer(reqBody))
		Expect(err).ToNot(HaveOccurred())

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		Expect(err).ToNot(HaveOccurred())

		Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))

		bytes, err := io.ReadAll(resp.Body)
		Expect(err).ToNot(HaveOccurred())

		Expect(string(bytes)).To(MatchJSON(` {
            "errors": {
                "alert": "Invalid login credentials"
            }
        }`))
	})
})
