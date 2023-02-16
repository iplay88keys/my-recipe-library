package integration_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Register", func() {
	BeforeEach(func() {
		_, err := db.Exec("DELETE FROM users WHERE id IS NOT NULL")
		Expect(err).ToNot(HaveOccurred())
	})

	It("creates a new user", func() {
		username := "some_user"
		body := []byte(fmt.Sprintf(`{
            "username": "%s",
            "email": "someone@example.com",
            "password": "Pa3$word123"
        }`, username))
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%s/api/v1/users/register", port), bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		Expect(err).ToNot(HaveOccurred())

		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		var count int
		row := db.QueryRow("SELECT COUNT(*) FROM users WHERE username=?", username)
		err = row.Scan(&count)
		Expect(err).ToNot(HaveOccurred())
		Expect(count).To(Equal(1))
	})

	It("returns an error if the json data is invalid", func() {
		body := []byte(`{
            "username": "a",
            "email": "a",
            "password": "a"
        }`)
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%s/api/v1/users/register", port), bytes.NewBuffer(body))
		Expect(err).ToNot(HaveOccurred())

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		Expect(err).ToNot(HaveOccurred())

		Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))

		bytes, err := io.ReadAll(resp.Body)
		Expect(err).ToNot(HaveOccurred())

		Expect(string(bytes)).To(MatchJSON(` {
            "errors": {
                "email": "Invalid email address",
                "password": "Uppercase letter missing, Numeric character missing, Special character missing, Must be between 6 and 64 characters long",
                "username": "Must be between 6 and 30 characters long"
            }
        }`))
	})
})
