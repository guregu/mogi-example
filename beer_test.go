package main

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/guregu/mogi"
)

var beerFixture = Beer{
	ID:   42,
	Name: "Yona Yona Ale",
	Pct:  5.5,
}

func TestGetBeer(t *testing.T) {
	defer mogi.Reset()
	mogi.Select("id", "name", "pct").
		From("beer").
		Where("id", 42).
		StubCSV(`42,Yona Yona Ale,5.5`)

	beer, err := GetBeer(42)
	if err != nil {
		t.Fatal("err should be nil, but is:", err)
	}
	if !reflect.DeepEqual(beer, beerFixture) {
		t.Errorf("%#v ≠ %#v", beer, beerFixture)
	}
}

func TestGetBeerMissing(t *testing.T) {
	defer mogi.Reset()
	mogi.Select().
		From("beer").
		Where("id", 99).
		StubError(sql.ErrNoRows)

	_, err := GetBeer(99)
	if err != sql.ErrNoRows {
		t.Error("err should be sql.ErrNoRows, but is:", err)
	}
}
