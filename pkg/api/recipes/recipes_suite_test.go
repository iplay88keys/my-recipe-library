package recipes_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestRecipes(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Recipes Suite")
}

var (
	osStdout *os.File
	osStderr *os.File
)

var _ = BeforeSuite(func() {
	osStdout = os.Stdout
	osStderr = os.Stderr

	os.Stdout = nil
	os.Stderr = nil
})

var _ = AfterSuite(func() {
	os.Stdout = osStdout
	os.Stderr = osStderr
})
