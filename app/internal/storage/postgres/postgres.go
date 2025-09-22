package postgres

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"CommentTree/pkg/errs"

	"github.com/lib/pq"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/retry"
)

const (
	// Tables names.
	CommentsTable = "comments"

	// Postgres errors.
	ViolatesForeignKey = "23503"

	// Params.
	ConnMaxLifeTime = 2
	MaxIdleConns    = 20
	MaxOpenConns    = 100

	// RetryStrategy.
	Attempts = 3
	Delay    = 3
	Backoff  = 2
)

var (
	retryOpts = retry.Strategy{
		Attempts: Attempts,
		Delay:    Delay,
		Backoff:  Backoff,
	}
)

type Postgres struct {
	db *dbpg.DB
}

func New(conn string) *Postgres {
	pg, err := dbpg.New(conn, nil, &dbpg.Options{
		MaxOpenConns:    MaxOpenConns,
		MaxIdleConns:    MaxIdleConns,
		ConnMaxLifetime: ConnMaxLifeTime * time.Hour,
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

func (p *Postgres) Shutdown() {
	_ = p.db.Master.Close()
	for i := 0; i < len(p.db.Slaves); i++ {
		_ = p.db.Slaves[i].Close()
	}
}

func (p *Postgres) unwrapInternalError(op string, err error) error {
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case ViolatesForeignKey:
			return errs.ErrDBViolatesForeignKey
		}
	}

	return fmt.Errorf("%s: %w", op, err)
}
