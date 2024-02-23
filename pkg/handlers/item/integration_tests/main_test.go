package integration_tests

import (
	"checkoutProject/pkg/bootstrap"
	"checkoutProject/pkg/common/database"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
)

var TestDB *gorm.DB

func TestMain(m *testing.M) {
	err := bootstrap.Initialize()
	if err != nil {
		log.Fatal(err.Error())
	}

	TestDB = database.GetInstance()
	os.Exit(m.Run())
}
