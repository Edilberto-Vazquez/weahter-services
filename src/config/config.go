package config

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var (
	ENVS = map[string]string{
		"APP_ENV":  "",
		"GIN_MODE": "",
		"PORT":     "",
		"DB_URI":   "",
	}
)

func SetEnvironment() {

	log.Println("Load environment vars")

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	prefix := ""

	switch os.Getenv("GIN_MODE") {
	case "development":
		prefix = "DEV_"
	case "debug":
		prefix = "DB_"
	case "test":
		prefix = "TEST_"
	default:
		prefix = ""
		gin.SetMode(gin.ReleaseMode)
	}

	ENVS["APP_ENV"] = os.Getenv(prefix + "APP_ENV")
	ENVS["PORT"] = os.Getenv(prefix + "PORT")
	ENVS["DB_URI"] = os.Getenv(prefix + "DB_URI")

	fmt.Println(ENVS)
}
