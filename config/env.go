package config

import (
	"os"

	"github.com/joho/godotenv"
)

type AppEnvConfig struct {
	MongoDBName     string
	MongoDBUrl      string
	MongoDBTestName string
	ServerPort      string
	JwtSecret       string
}

func (c *AppEnvConfig) SetMongoDbName(name string) {
	c.MongoDBName = name
}

var EnvConfig AppEnvConfig

func init() {

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	EnvConfig.MongoDBName = os.Getenv("MONGO_DB_NAME")
	EnvConfig.MongoDBUrl = os.Getenv("MONGO_DB_URL")
	EnvConfig.MongoDBTestName = os.Getenv("MONGO_DB_TEST_NAME")
	EnvConfig.ServerPort = os.Getenv("PORT")
	EnvConfig.JwtSecret = os.Getenv("JWT_SECRET")
}
