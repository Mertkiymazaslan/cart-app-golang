package env

import (
	"fmt"
	"os"
)

var (
	DB_URL string
)

func Load() error {
	var ok bool
	errorMessage := "cannot read %s from environment variables"

	environment, ok := os.LookupEnv("ENVIRONMENT")
	if !ok {
		return fmt.Errorf(errorMessage, "ENVIRONMENT")
	}

	if environment == "PRODUCTION" {
		DB_URL, ok = os.LookupEnv("DB_URL")
		if !ok {
			return fmt.Errorf(errorMessage, "DB_URL")
		}
	} else {
		testDBUrl, ok := os.LookupEnv("TEST_DB_URL")
		if !ok {
			return fmt.Errorf(errorMessage, "TEST_DB_URL")
		}
		DB_URL = testDBUrl
	}

	return nil
}
