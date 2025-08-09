package messages

import (
	"database/sql"
)

type Record struct {
	Prompt, Response string
}

func Create(db *sql.DB, role, text string) {
	_, _ = db.Exec("INSERT INTO messages(role, text) VALUES (?, ?)", role, text)
}

func FetchHistory(db *sql.DB) ([]Record, error) {
	rows, err := db.Query("SELECT role, text FROM messages ORDER BY id")
	if err != nil {
		return nil, err
	}
	var recs []Record
	for rows.Next() {
		var r Record
		if err := rows.Scan(&r.Prompt, &r.Response); err != nil {
			return nil, err
		}
		recs = append(recs, r)
	}
	return recs, nil
}
