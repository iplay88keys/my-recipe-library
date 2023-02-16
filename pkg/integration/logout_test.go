package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/iplay88keys/my-recipe-library/pkg/repositories"

	"github.com/iplay88keys/my-recipe-library/pkg/api/users"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("logout", func() {
	var (
		username string
		password string
		token    string
	)

	BeforeEach(func() {
		_, err := db.Exec("DELETE FROM users WHERE id IS NOT NULL")
		Expect(err).ToNot(HaveOccurred())

		username = "logout_user"
		password = "Pa3$word123"

		usersRepo := repositories.NewUsersRepository(db)
		_, err = usersRepo.Insert(username, username+"@example.com", password)
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

		Expect(resp.StatusCode).To(Equal(200))

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		Expect(err).ToNot(HaveOccurred())

		var loginResponse users.UserLoginResponse
		err = json.Unmarshal(body, &loginResponse)
		Expect(err).ToNot(HaveOccurred())

		token = loginResponse.AccessToken
	})

	It("logs the user out", func() {
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%s/api/v1/users/logout", port), nil)
		Expect(err).ToNot(HaveOccurred())

		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", "bearer "+token)

		resp, err := client.Do(req)
		Expect(err).ToNot(HaveOccurred())

		Expect(resp.StatusCode).To(Equal(http.StatusOK))
	})

	It("returns unauthorized if the token is invalid", func() {
		req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://localhost:%s/api/v1/users/logout", port), nil)
		Expect(err).ToNot(HaveOccurred())

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		Expect(err).ToNot(HaveOccurred())

		Expect(resp.StatusCode).To(Equal(http.StatusUnauthorized))
	})
})
