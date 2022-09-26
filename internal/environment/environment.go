package environment

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

var lock = &sync.Mutex{}

// Single is the singleton instance of the environment
type Single struct {
	ENVIRONMENT string // nolint: golint
	APP_VERSION string // nolint: golint
	APP_PORT    string // nolint: golint
	LOG_LEVEL   string // nolint: golint
}

func init() {
	envVar := os.Getenv("ENVIRONMENT")

	err := godotenv.Load(".env.%s", envVar)
	if err != nil {
		log.Println("Error loading .env.local file")
	}

	env := GetInstance()
	env.Setup()
}

// Setup sets up the environment
func (e *Single) Setup() {
	e.ENVIRONMENT = os.Getenv("ENVIRONMENT")
	e.APP_VERSION = os.Getenv("APP_VERSION")
	e.APP_PORT = getenv("APPLICATION_PORT", "7070")
	e.LOG_LEVEL = getenv("LOG_LEVEL", "debug")
}

// IsDevelopment returns true if the environment is development
func (e *Single) IsDevelopment() bool {
	return e.ENVIRONMENT == "local"
}

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

var singleInstance *Single

// GetInstance returns the singleton instance of the environment
func GetInstance() *Single {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			fmt.Println("Creating single instance now.")
			singleInstance = &Single{}
			singleInstance.Setup()
		} else {
			fmt.Println("Single instance already created.")
		}
	}

	return singleInstance
}
