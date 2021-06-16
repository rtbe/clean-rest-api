// Package config contains configuration for an application.
package config

import (
	"os"
	"sync"
)

var (
	once   sync.Once
	config Cfg
)

const (
	apiPort        = "API_PORT"
	dbPort         = "DB_PORT"
	dbHost         = "DB_HOST"
	dbUser         = "DB_USER"
	dbPassword     = "DB_PASSWORD"
	dbName         = "DB_NAME"
	authDbPort     = "AUTH_DB_PORT"
	authDbHost     = "AUTH_DB_HOST"
	authDbUser     = "AUTH_DB_USER"
	authDbPassword = "AUTH_DB_PASSWORD"
	authDbName     = "AUTH_DB_PASSWORD"
	jwtSalt        = "JWT_SALT"
)

// Cfg is an struct that holds environment variables.
type Cfg struct {
	APIPort        string
	DbPort         string
	AuthDbPort     string
	DbHost         string
	DbUser         string
	DbPassword     string
	DbName         string
	AuthDbHost     string
	AuthDbUser     string
	AuthDbPassword string
	AuthDBName     string
	JWTSalt        string
}

// New constructs an config from environment variables.
func New() *Cfg {
	once.Do(
		func() {
			config = Cfg{
				APIPort:        parseEnvString(apiPort, "8080"),
				DbPort:         parseEnvString(dbPort, "5432"),
				DbHost:         parseEnvString(dbHost, "db"),
				DbUser:         parseEnvString(dbUser, "admin"),
				DbPassword:     parseEnvString(dbPassword, "password"),
				DbName:         parseEnvString(dbUser, "admin"),
				AuthDbPort:     parseEnvString(authDbPassword, "27017"),
				AuthDbHost:     parseEnvString(authDbHost, "authDB"),
				AuthDbUser:     parseEnvString(authDbUser, "admin"),
				AuthDbPassword: parseEnvString(authDbPassword, "password"),
				AuthDBName:     parseEnvString(authDbName, "admin"),
				JWTSalt:        parseEnvString(jwtSalt, "secret123"),
			}
		},
	)
	return &config
}

// parseEnvString looks for environment variable value
// and if value is not found returns provided default value.
func parseEnvString(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = defaultValue
	}
	return value
}
