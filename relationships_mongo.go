// +build mongo

package pshgo

import (
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/globalsign/mgo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func (r Relationships) MongoDB(name string) (*mgo.Database, error) {
	rels, ok := r[name]
	if !ok {
		return nil, errors.New("missing relationship")
	}

	if len(rels) == 0 {
		return nil, errors.New("empty relationship")
	}

	hosts := make([]string, len(rels))
	for idx, rel := range rels {
		hosts[idx] = net.JoinHostPort(rel.Host, strconv.Itoa(rel.Port))
	}

	var b strings.Builder
	b.WriteString("mongodb://")
	b.WriteString(strings.Join(hosts, ","))
	b.WriteString("/")
	b.WriteString(rels[0].Path)

	url := b.String()

	mgo.SetLogger(log.New(os.Stderr, "mongo ", log.LstdFlags))

	for count := 0; count < 10; count++ {
		sess, err := mgo.Dial(url)
		if err == nil {
			db := sess.DB(rels[0].Path)
			if err := db.Login(rels[0].Username, rels[0].Password); err != nil {
				logrus.WithError(err).Warn("error logging into mongo database")
			}

			return db, nil
		}
		logrus.WithError(err).WithField("attempt", count+1).Warn("failed to connect to mongo server")
	}

	return nil, errors.New("failed to connect to mongo server")
}
