package authentication

import (
	"testing"
)

func TestCreateToken(t *testing.T) {

	var testTables = []struct {
		username string
		userId   int
	}{
		{"sandy", 4444},
	}

	for _, row := range testTables {
		token := CreateToken(row.username, row.userId)
		if len(token) == 0 {
			t.Errorf("TOken creation failed")
		}
	}
}
