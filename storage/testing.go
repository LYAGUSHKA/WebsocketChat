package storage

import (
	"fmt"
	"strings"
	"testing"
)

func TestStore(t *testing.T, databaseURL string) (*Storage, func(...string)) {
	t.Helper()

	config := NewConfig()
	config.DatabaseURL = databaseURL
	s := New(config)
	if err := s.Open(); err != nil {
		t.Fatal(err)
	}

	return s, func(tables ...string) {
		if len(tables) > 0 {
			if _, err := s.Db.Exec(fmt.Sprintf("DELETE FROM %s; VACUUM", strings.Join(tables, ", "))); err != nil {
				t.Fatal(err)
			}
		}

		s.Close()
	}

}