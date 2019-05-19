// +build postgresql

package pshgo

import (
	"math/rand"
	"strings"

	"github.com/lib/pq"
	"github.com/pkg/errors"
)

func (r Relationships) Postgresql(name string) (string, error) {
	rels, ok := r[name]
	if !ok {
		return "", errors.New("missing relationship")
	}

	if len(rels) == 0 {
		return "", errors.New("empty relationship")
	}

	for len(rels) > 0 {
		rand.Shuffle(len(rels), func(i, j int) {
			rels[i], rels[j] = rels[j], rels[i]
		})

		dbURL := rels[0].URL(true, false)
		dbURL = strings.Replace(dbURL, "pgsql://", "postgresql://", 1)
		dbOpen, err := pq.ParseURL(dbURL)
		if err == nil {
			dbOpen += " sslmode=disable"
			return dbOpen, nil
		}

		rels = rels[1:]
	}

	return "", errors.New("error parsing postgres url")
}
