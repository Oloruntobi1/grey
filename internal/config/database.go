package config

import (
	"fmt"
)

func GetDatabaseConfig() string {
	user := getEnv("POSTGRES_USER", "db_user")
	password := getEnv("POSTGRES_PASSWORD", "db_pass")
	host := getEnv("POSTGRES_HOST", "localhost")
	port := getEnv("POSTGRES_PORT", "5444")
	name := getEnv("POSTGRES_DB_NAME", "grey-app-db")
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		user, password, host, port, name)
}
