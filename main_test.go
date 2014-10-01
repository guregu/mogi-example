package main

import (
	"database/sql"

	"github.com/guregu/mogi"
)

func init() {
	db, _ = sql.Open("mogi", "")

	mogi.Verbose(true)
}
