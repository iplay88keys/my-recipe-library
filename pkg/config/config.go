package config

import (
    "encoding/json"
    "fmt"
    "os"
    "reflect"
    "strings"
)

type Config struct {
    MySQLCreds    MySQLCreds `env:"MYSQL_CREDS,    required"`
    RedisURL      string     `env:"REDIS_URL,      required"`
    AccessSecret  string     `env:"ACCESS_SECRET,  required"`
    RefreshSecret string     `env:"REFRESH_SECRET, required"`
    Port          string     `env:"PORT"`
    Static        string     `env:"STATIC_DIR"`
}

type MySQLCreds struct {
    URL          string `json:"url"`
    InstanceName string `json:"gcloud_instance_name"`
    DBName       string `json:"gcloud_db_name"`
    User         string `json:"gcloud_user"`
    Password     string `json:"gcloud_password"`
}

type Unmarshaller interface {
    UnmarshalEnv(data string) error
}

func (d *MySQLCreds) UnmarshalEnv(data string) error {
    return json.Unmarshal([]byte(data), d)
}

func Load() (*Config, error) {
    cfg := &Config{}

    v := reflect.ValueOf(cfg).Elem()

    for i := 0; i < v.NumField(); i++ {
        tag := v.Type().Field(i).Tag

        env, ok := tag.Lookup("env")
        if !ok {
            return nil, fmt.Errorf("expected %s to contain env tag", v.Field(i))
        }

        envSplit := strings.Split(env, ",")
        envVal := os.Getenv(envSplit[0])

        if envVal == "" {
            if containsRequired(envSplit) {
                return nil, fmt.Errorf("expected %s to have env var set", v.Field(i))
            }

            continue
        }

        switch v.Field(i).Kind() {
        case reflect.String:
            v.Field(i).SetString(envVal)
        case reflect.Struct:
            var unmarshaller Unmarshaller
            if v.Field(i).CanAddr() {
                unmarshaller = v.Field(i).Addr().Interface().(Unmarshaller)
            } else {
                unmarshaller = v.Field(i).Interface().(Unmarshaller)
            }

            err := unmarshaller.UnmarshalEnv(envVal)
            if err != nil {
                return nil, fmt.Errorf("failed to unmarshal struct: %s", err.Error())
            }
        default:
            return nil, fmt.Errorf("unsupported type: %s", v.Field(i).Kind())
        }
    }

    return cfg, nil
}

func containsRequired(tags []string) bool {
    for _, tag := range tags {
        if strings.Contains(tag, "required") {
            return true
        }
    }

    return false
}
