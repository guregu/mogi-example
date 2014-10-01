package main

import (
	"database/sql"

	_ "github.com/guregu/mogi"
)

func init() {
	db, _ = sql.Open("mogi", "")
}
