package postgres

import (
	"fmt"
	"os"
	"time"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/retry"
)

const (
	CommentsTable = "comments"

	PageElements = 1
)

var (
	retryOpts = retry.Strategy{
		Attempts: 3,
		Delay:    3,
		Backoff:  2,
	}
)

type Postgres struct {
	db *dbpg.DB
}

func New(conn string) *Postgres {
	pg, err := dbpg.New(conn, nil, &dbpg.Options{
		MaxOpenConns:    100,
		MaxIdleConns:    20,
		ConnMaxLifetime: 2 * time.Hour,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = pg.Master.Ping()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return &Postgres{
		db: pg,
	}
}
