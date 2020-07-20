package api_test

import (
    "os"
    "testing"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestApi(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "API Suite")
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
