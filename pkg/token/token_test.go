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
			req, err := http.NewRequest(http.MethodPost, "example.com", nil)
			Expect(err).ToNot(HaveOccurred())

			// token from previous run of token.CreateToken(0)
			req.Header.Set(
				"Authorization",
				"bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJVc2VySUQiOjEwLCJBY2Nlc3NVVUlEIjoiNzc5ZjRlYzktYTA0My00ZjU1LWI2MDQtZGNlYmE3NmQwZTIyIiwiUmVmcmVzaFVVSUQiOiIiLCJleHAiOjE1OTQ3NDA3Mzh9.QzKb9sF-XRYD9gs8slrT7mlGObubQIsFkazgxv14b6U",
			)

			s := token.NewService("secret value", "refresh value")
			accessDetails, err := s.ValidateToken(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(accessDetails).To(Equal(&token.AccessDetails{
				AccessUuid: "779f4ec9-a043-4f55-b604-dceba76d0e22",
				UserId:     10,
			}))
		})

		It("returns false if the token is invalid", func() {
			req, err := http.NewRequest(http.MethodPost, "example.com", nil)
			Expect(err).ToNot(HaveOccurred())

			// token from previous run of token.CreateToken(0)
			req.Header.Set(
				"Authorization",
				"bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjowLCJleHAiOjE1OTQ2MTY0MzAsImlhdCI6MTU5NDYxMjgzMH0.9Bhbk7G25LG-fU9w_HOxDUw3u16cG0QULwIIUENXdJ4",
			)

			s := token.NewService("wrong value", "refresh value")
			accessDetails, err := s.ValidateToken(req)
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("signature is invalid"))
			Expect(accessDetails).To(BeNil())
		})
	})
})
