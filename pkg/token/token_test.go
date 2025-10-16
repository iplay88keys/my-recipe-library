package token_test

import (
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/iplay88keys/my-recipe-library/pkg/token"
)

var _ = Describe("token", func() {
	Context("CreateToken", func() {
		It("creates a token", func() {
			s := token.NewService("secret value", "refresh value")
			details, err := s.CreateToken(10)
			Expect(err).ToNot(HaveOccurred())

			Expect(len(details.AccessToken)).To(BeNumerically(">", 0))
			Expect(len(details.RefreshToken)).To(BeNumerically(">", 0))

			Expect(details.AccessExpires).To(BeNumerically(">", time.Now().Unix()))
			Expect(details.RefreshExpires).To(BeNumerically(">", time.Now().Unix()))

			Expect(details.AccessUuid).ToNot(Equal(""))
			Expect(details.RefreshUuid).ToNot(Equal(""))
		})
	})

	Context("ValidateToken", func() {
		It("returns user info if the token is valid", func() {
			s := token.NewService("secret value", "refresh value")

			userID := int64(10)
			tokenDetails, err := s.CreateToken(userID)
			Expect(err).ToNot(HaveOccurred())

			req, err := http.NewRequest(http.MethodPost, "example.com", nil)
			Expect(err).ToNot(HaveOccurred())

			req.Header.Set("Authorization", "bearer "+tokenDetails.AccessToken)

			accessDetails, err := s.ValidateToken(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(accessDetails).To(Equal(&token.AccessDetails{
				AccessUuid: tokenDetails.AccessUuid,
				UserId:     userID,
			}))
		})

		It("returns false if the token is invalid", func() {
			// Create a token with one secret
			s1 := token.NewService("secret value", "refresh value")
			tokenDetails, err := s1.CreateToken(0)
			Expect(err).ToNot(HaveOccurred())

			req, err := http.NewRequest(http.MethodPost, "example.com", nil)
			Expect(err).ToNot(HaveOccurred())

			req.Header.Set("Authorization", "bearer "+tokenDetails.AccessToken)

			// Try to validate with a different secret
			s2 := token.NewService("wrong value", "refresh value")
			accessDetails, err := s2.ValidateToken(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("signature is invalid"))
			Expect(accessDetails).To(BeNil())
		})

		It("returns error if the token is malformed", func() {
			req, err := http.NewRequest(http.MethodPost, "example.com", nil)
			Expect(err).ToNot(HaveOccurred())

			req.Header.Set("Authorization", "bearer invalid.token.here")

			s := token.NewService("secret value", "refresh value")
			accessDetails, err := s.ValidateToken(req)
			Expect(err).To(HaveOccurred())
			Expect(accessDetails).To(BeNil())
		})

		It("returns error if authorization header is missing", func() {
			req, err := http.NewRequest(http.MethodPost, "example.com", nil)
			Expect(err).ToNot(HaveOccurred())

			s := token.NewService("secret value", "refresh value")
			accessDetails, err := s.ValidateToken(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("authorization header missing token"))
			Expect(accessDetails).To(BeNil())
		})
	})
})
