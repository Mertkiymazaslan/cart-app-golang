package testhelper

import (
	"database/sql"
	"gopkg.in/testfixtures.v2"
	"testing"
)

var (
	DefaultItemsFixturePath = "fixtures/defaultItems"
	DigitalItemFixturesPath = "fixtures/digitalItems"
	AddVasItemFixturesPath  = "fixtures/addVasItemFixtures"
	DefaultPath             = "fixtures"
)

func LoadFixtures(path string, t *testing.T, db *sql.DB) {
	fixtures, err := testfixtures.NewFolder(db, &testfixtures.PostgreSQL{}, path)
	if err != nil {
		t.Fatalf("error while getting fixtures: %v", err)
	}

	if err := fixtures.Load(); err != nil {
		t.Fatalf("error while loading fixtures: %v", err)
	}
}
