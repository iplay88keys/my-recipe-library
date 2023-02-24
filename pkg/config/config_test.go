package config_test

import (
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    "os"

    "github.com/iplay88keys/my-recipe-library/pkg/config"
)

const sqlCreds = "{\"url\": \"sql-url\",\"gcloud_instance_name\": \"instance\",\"gcloud_db_name\": \"db\",\"gcloud_user\": \"user\",\"gcloud_password\": \"password\"}"

var _ = Describe("Config", func() {
    It("Loads the config from env vars", func() {
        env := make(map[string]string)
        env["MYSQL_CREDS"] = sqlCreds
        env["REDIS_URL"] = "redis-url"
        env["ACCESS_SECRET"] = "access"
        env["REFRESH_SECRET"] = "refresh"
        env["PORT"] = "1234"
        env["STATIC_DIR"] = "static"

        for key, val := range env {
            err := os.Setenv(key, val)
            Expect(err).ToNot(HaveOccurred())
        }

        defer unsetEnvs(env)

        cfg, err := config.Load()
        Expect(err).ToNot(HaveOccurred())
        Expect(cfg).To(Equal(&config.Config{
            MySQLCreds: config.MySQLCreds{
                URL:          "sql-url",
                InstanceName: "instance",
                DBName:       "db",
                User:         "user",
                Password:     "password",
            },
            RedisURL:      "redis-url",
            AccessSecret:  "access",
            RefreshSecret: "refresh",
            Port:          "1234",
            Static:        "static",
        }))
    })

    It("Loads the config when only the required vars are set", func() {
        env := make(map[string]string)
        env["MYSQL_CREDS"] = sqlCreds
        env["REDIS_URL"] = "redis-url"
        env["ACCESS_SECRET"] = "access"
        env["REFRESH_SECRET"] = "refresh"

        for key, val := range env {
            err := os.Setenv(key, val)
            Expect(err).ToNot(HaveOccurred())
        }

        defer unsetEnvs(env)

        _, err := config.Load()
        Expect(err).ToNot(HaveOccurred())
    })

    It("Throws an error if the required fields are not set", func() {
        env := make(map[string]string)
        env["PORT"] = "1234"
        env["STATIC_DIR"] = "static"

        for key, val := range env {
            err := os.Setenv(key, val)
            Expect(err).ToNot(HaveOccurred())
        }

        defer unsetEnvs(env)

        _, err := config.Load()
        Expect(err).To(HaveOccurred())
    })
})

func unsetEnvs(env map[string]string) {
    for key := range env {
        _ = os.Unsetenv(key)
    }
}
