package postgres

import (
	"CommentTree/pkg/errs"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/lib/pq"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/retry"
)

const (
	// Tables names.
	CommentsTable = "comments"

	// TODO: think about move it to comment entity.
	// Elements per page for pagination.
	PageElements = 1

	// Postgres errors.
	ViolatesForeignKey = "23503"
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

	err = pg.Master.PingContext(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	return &Postgres{
		db: pg,
	}
}

func (p *Postgres) UnwrapError(err error) error {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case ViolatesForeignKey:
			return errs.ErrDBViolatesForeignKey
		}
	}

	return err
}

func (p *Postgres) Shutdown() {
	_ = p.db.Master.Close()
	for i := 0; i < len(p.db.Slaves); i++ {
		_ = p.db.Slaves[i].Close()
	}
}
