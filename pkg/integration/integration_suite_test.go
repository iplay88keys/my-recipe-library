package integration_test

import (
    "database/sql"
    "fmt"
    "net/http"
    "os"
    "os/exec"
    "testing"
    "time"

    "github.com/onsi/gomega/gexec"

    _ "github.com/go-sql-driver/mysql"

    . "github.com/iplay88keys/my-recipe-library/pkg/helpers"

    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

func TestIntegration(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Integration Suite")
}

var (
    pathToExecutable      string
    databaseURL           string
    databaseVarsAvailable bool
    db                    *sql.DB
    port                  string
    client                *http.Client
    session               *gexec.Session
)

var _ = BeforeSuite(func() {
    user := os.Getenv("DATABASE_USERNAME")
    password := os.Getenv("DATABASE_PASSWORD")
    host := os.Getenv("DATABASE_HOST")
    port := os.Getenv("DATABASE_PORT")
    databaseName := os.Getenv("DATABASE_NAME")

    if user != "" && password != "" && host != "" && port != "" && databaseName != "" {
        databaseURL = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, databaseName)
        databaseVarsAvailable = true

        var err error
        db, err = sql.Open("mysql", databaseURL)
        if err != nil {
            panic(err)
        }

        pathToExecutable, err = gexec.Build("github.com/iplay88keys/my-recipe-library")
        Expect(err).ToNot(HaveOccurred())
    }

    client = &http.Client{
        Timeout: 10 * time.Second,
    }
})

var _ = AfterSuite(func() {
    gexec.CleanupBuildArtifacts()
})

var _ = BeforeEach(func() {
    if !databaseVarsAvailable {
        Skip("Missing some database information. Run the tests from 'scripts/test.sh' to start up the database.")
    }

    var err error
    port, err = GetRandomPort()
    Expect(err).ToNot(HaveOccurred())

    err = os.Setenv("PORT", port)
    Expect(err).ToNot(HaveOccurred())

    err = os.Setenv("DATABASE_URL", databaseURL)
    Expect(err).ToNot(HaveOccurred())

    cmd := exec.Command(pathToExecutable)

    session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
    Expect(err).ToNot(HaveOccurred())
    time.Sleep(1 * time.Second)
})

var _ = AfterEach(func() {
    session.Kill()
    Eventually(session).Should(gexec.Exit())
})
