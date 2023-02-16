package config

import (
	"encoding/json"
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

func (d *MySQLCreds) UnmarshalEnv(data string) error {
	return json.Unmarshal([]byte(data), d)
}
