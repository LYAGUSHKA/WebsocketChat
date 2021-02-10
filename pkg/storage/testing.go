package sqlstorage

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbURL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func TestStore(t *testing.T, databaseURL string) (*Storage, func(...string)) {
	t.Helper()

	db, err := newDB(databaseURL)
	if err != nil {
		t.Fatal(err)
	}
	s := New(db)

	return s, func(tables ...string) {
		if len(tables) > 0 {
			if _, err := s.Db.Exec(fmt.Sprintf("DELETE FROM %s; VACUUM", strings.Join(tables, ", "))); err != nil {
				t.Fatal(err)
			}
		}

		db.Close()

	}
}
