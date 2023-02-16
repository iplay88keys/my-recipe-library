package users_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/iplay88keys/my-recipe-library/pkg/api/users"
)

var _ = Describe("RegisterRequest", func() {
	Context("Validate", func() {
		It("verifies all fields exist", func() {
			req := users.RegisterRequest{
				Username: "",
				Email:    "",
				Password: "",
			}

			errors := req.Validate(false, false)
			Expect(errors).To(Equal(map[string]string{
				"username": "Required",
				"email":    "Required",
				"password": "Required",
			}))
		})

		Context("username", func() {
			It("returns no errors if the username is valid", func() {
				req := users.RegisterRequest{
					Username: "test_ing123",
					Email:    "email@example.com",
					Password: "1aA$123",
				}

				errors := req.Validate(false, false)
				Expect(errors).To(Equal(map[string]string{}))
			})

			It("returns an error if the username is already in use", func() {
				req := users.RegisterRequest{
					Username: "test_ing123",
					Email:    "email@example.com",
					Password: "1aA$123",
				}

				errors := req.Validate(true, false)
				Expect(errors).To(Equal(map[string]string{
					"username": "Username already in use",
				}))
			})

			It("returns an error if the username is shorter than 6 characters", func() {
				req := users.RegisterRequest{
					Username: "test",
					Email:    "email@example.com",
					Password: "1aA$123",
				}

				errors := req.Validate(false, false)
				Expect(errors).To(Equal(map[string]string{
					"username": "Must be between 6 and 30 characters long",
				}))
			})

			It("returns an error if the username is longer than 30 characters", func() {
				username := ""
				for i := 0; i < 31; i++ {
					username = username + "a"
				}

				req := users.RegisterRequest{
					Username: username,
					Email:    "email@example.com",
					Password: "1aA$123",
				}

				errors := req.Validate(false, false)
				Expect(errors).To(Equal(map[string]string{
					"username": "Must be between 6 and 30 characters long",
				}))
			})

			It("returns an error if the username starts with a number", func() {
				req := users.RegisterRequest{
					Username: "1test_ing",
					Email:    "email@example.com",
					Password: "1aA$123",
				}

				errors := req.Validate(false, false)
				Expect(errors).To(Equal(map[string]string{
					"username": "Cannot start with a number",
				}))
			})

			It("returns an error if the username starts with a underscore", func() {
				req := users.RegisterRequest{
					Username: "_testing123",
					Email:    "email@example.com",
					Password: "1aA$123",
				}

				errors := req.Validate(false, false)
				Expect(errors).To(Equal(map[string]string{
					"username": "Cannot start with an underscore",
				}))
			})

			It("returns an error if the username has any special character except for an underscore", func() {
				req := users.RegisterRequest{
					Username: "test#ing123",
					Email:    "email@example.com",
					Password: "1aA$123",
				}

				errors := req.Validate(false, false)
				Expect(errors).To(Equal(map[string]string{
					"username": "Only alphanumeric characters and underscores (_) allowed",
				}))
			})
		})

		Context("email", func() {
			It("returns no errors if the email address is valid", func() {
				req := users.RegisterRequest{
					Username: "testing",
					Email:    "good-email@example.com",
					Password: "Pa3$word123",
				}

				errors := req.Validate(false, false)
				Expect(errors).To(Equal(map[string]string{}))
			})

			It("returns an error if the email address is in use", func() {
				req := users.RegisterRequest{
					Username: "testing",
					Email:    "good-email@example.com",
					Password: "Pa3$word123",
				}

				errors := req.Validate(false, true)
				Expect(errors).To(Equal(map[string]string{
					"email": "Email already in use",
				}))
			})

			It("returns an error if the email address is invalid", func() {
				req := users.RegisterRequest{
					Username: "testing",
					Email:    "bad-email",
					Password: "Pa3$word123",
				}

				errors := req.Validate(false, false)
				Expect(errors).To(Equal(map[string]string{
					"email": "Invalid email address",
				}))
			})
		})

		Context("password", func() {
			It("returns no errors if the password is valid", func() {
				req := users.RegisterRequest{
					Username: "testing",
					Email:    "email@example.com",
					Password: "1aA$123",
				}

				errors := req.Validate(false, false)
				Expect(errors).To(Equal(map[string]string{}))
			})

			It("returns an error if the password is shorter than 6 characters", func() {
				req := users.RegisterRequest{
					Username: "testing",
					Email:    "email@example.com",
					Password: "1aA$",
				}

				errors := req.Validate(false, false)
				Expect(errors).To(Equal(map[string]string{
					"password": "Must be between 6 and 64 characters long",
				}))
			})

			It("returns an error if the password is longer than 64 characters", func() {
				password := "1aA$"
				for i := 0; i < 61; i++ {
					password = password + "a"
				}

				req := users.RegisterRequest{
					Username: "testing",
					Email:    "email@example.com",
					Password: password,
				}

				errors := req.Validate(false, false)
				Expect(errors).To(Equal(map[string]string{
					"password": "Must be between 6 and 64 characters long",
				}))
			})

			It("returns an error if the password is missing lower case letters", func() {
				req := users.RegisterRequest{
					Username: "testing",
					Email:    "email@example.com",
					Password: "123456A$",
				}

				errors := req.Validate(false, false)
				Expect(errors).To(Equal(map[string]string{
					"password": "Lowercase letter missing",
				}))
			})

			It("returns an error if the password is missing upper case letters", func() {
				req := users.RegisterRequest{
					Username: "testing",
					Email:    "email@example.com",
					Password: "123456a$",
				}

				errors := req.Validate(false, false)
				Expect(errors).To(Equal(map[string]string{
					"password": "Uppercase letter missing",
				}))
			})

			It("returns an error if the password is missing a special character", func() {
				req := users.RegisterRequest{
					Username: "testing",
					Email:    "email@example.com",
					Password: "123456aA",
				}

				errors := req.Validate(false, false)
				Expect(errors).To(Equal(map[string]string{
					"password": "Special character missing",
				}))
			})
		})
	})
})
