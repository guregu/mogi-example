package main

type Beer struct {
	ID   int64
	Name string
	Pct  float32
}

func GetBeer(id int64) (beer Beer, err error) {
	query := `SELECT id, name, pct FROM beer WHERE id = ?`
	err = db.QueryRow(query, id).Scan(&beer.ID, &beer.Name, &beer.Pct)
	return
}
