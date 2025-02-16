package setup

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type env struct {
	DatabaseUrl string
	ServerPort  int
}

func GetEnv() env {
	out := env{}
	out.ServerPort = envAtoiOr("SERVER_PORT", 8080)
	out.DatabaseUrl = envBuildDatabaseUrl()
	return out
}

func envBuildDatabaseUrl() string {
	password := envOrPanic("DATABASE_PASSWORD")

	user := envOr("DATABASE_USER", "postgres")
	host := envOr("DATABASE_HOST", "localhost")
	port := envOr("DATABASE_PORT", "5432")
	dbName := envOr("DATABASE_NAME", user)
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, password, host, port, dbName)
}

func envOrPanic(key string) string {
	value := envOr(key, "")
	if len(value) == 0 {
		log.Panic(fmt.Errorf("env variable not set: %s", key))
	}
	return value
}

func envAtoiOr(key string, defaultValue int) int {
	value := envOr(key, strconv.Itoa(defaultValue))

	out, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return out
}

func envOr(key string, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok || len(value) == 0 {
		return defaultValue
	}
	return value
}
