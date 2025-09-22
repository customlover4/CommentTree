package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"CommentTree/internal/entities/comment"
	"CommentTree/pkg/errs"

	"github.com/wb-go/wbf/zlog"
)

func (p *Postgres) issetComment(id int64) bool {
	const op = "internal.storage.postgres.comments.issetComment"

	// Invisible root comment for others with id 0.
	if id == 0 {
		return true
	}

	q := fmt.Sprintf("select * from %s where id = $1", CommentsTable)

	var tmp comment.Comment
	var nullPID sql.NullInt64
	err := p.db.Master.QueryRowContext(context.Background(), q, id).
		Scan(&tmp.ID, &tmp.Message, &nullPID)
	if errors.Is(err, sql.ErrNoRows) {
		return false
	} else if err != nil {
		zlog.Logger.Error().Err(fmt.Errorf("%s: %w", op, err)).Send()
		return false
	}

	return true
}

func (p *Postgres) CreateComment(c comment.Comment) (int64, error) {
	const op = "internal.storage.postgres.comments.Create"

	if !p.issetComment(c.ParentID) {
		return 0, errs.ErrDBViolatesForeignKey
	}

	q := fmt.Sprintf(
		"insert into %s (message, parent_id) values ($1, $2) returning id;",
		CommentsTable,
	)

	// if we don't have parent id, we will insert NULL to db
	haveParentID := (c.ParentID > 0)
	parentIDArg := sql.NullInt64{
		Int64: c.ParentID, Valid: haveParentID,
	}

	var id int64
	err := p.db.Master.QueryRowContext(
		context.Background(), q, c.Message, parentIDArg,
	).Scan(&id)
	if err != nil {
		return 0, p.unwrapInternalError(op, err)
	}

	return id, nil
}

func (p *Postgres) generateAgrigatedOpts(parentID int64, opts *comment.GetterOpts) ([]any, string) {
	if opts == nil {
		opts = &comment.GetterOpts{}
	}

	// Initialization values.
	args := []any{}
	argsStr := fmt.Sprintf(`
		select * from %s
	`, CommentsTable)

	argIDX := 1
	var addArg bool
	if parentID > 0 {
		argsStr += fmt.Sprintf(" where parent_id = $%d", argIDX)
		args = append(args, parentID)
		argIDX++
		addArg = true
	} else if !opts.SearchGlobal {
		/*
			If we want to search for substr in all comments,
			not just in roots, we must add a new flag,
			because without this flag, we will always set the parent_id
			value to NULL when ParentID == 0.
			Now, when ParentID == 0 and the SearchGlobal flag is set,
			it will not add filters for id.
		*/
		argsStr += " where parent_id is NULL"
		addArg = true
	}

	if opts.Substr != "" {
		// If we didn't set id, we wouldn't add 'and'.
		if addArg {
			argsStr += " and"
		} else {
			argsStr += " where"
		}

		argsStr += fmt.Sprintf(" POSITION($%d IN message) > 0", argIDX)
		args = append(args, opts.Substr)
	}

	if opts.Page > 0 {
		// 1(PageElement) * 1(opts.Page is first) = 1, wrong offset for first page.
		opts.Page--
	}

	limit := comment.PageElements
	offset := comment.PageElements * opts.Page
	argsStr += fmt.Sprintf(" limit %d offset %d", limit, offset)

	return args, argsStr
}

func (p *Postgres) Parent(id int64) (comment.Comment, error) {
	const op = "internal.storage.postgres.comments.parent"

	qParent := fmt.Sprintf(`select * from %s where id = $1`, CommentsTable)

	var parent comment.Comment
	var nullID sql.NullInt64

	err := p.db.Master.QueryRowContext(context.Background(), qParent, id).
		Scan(&parent.ID, &parent.Message, &nullID)

	if errors.Is(err, sql.ErrNoRows) {
		return comment.Comment{}, errs.ErrDBNotFound
	} else if err != nil {
		return comment.Comment{}, fmt.Errorf("%s: %w", op, err)
	}

	if nullID.Valid {
		parent.ParentID = nullID.Int64
	}

	return parent, nil
}

func (p *Postgres) lastID(parentID int64) int64 {
	const op = "internal.storage.postgres.lastID"

	args := make([]any, 0, 1)

	q := fmt.Sprintf(`
		select id from %s where parent_id`,
		CommentsTable,
	)

	// Because we can load more from root comments and from child comments.
	if parentID == 0 {
		q += " is null"
	} else {
		q += "=$1"
		args = append(args, parentID)
	}

	q += ` order by id desc limit 1;`

	var id int64
	err := p.db.Master.QueryRowContext(
		context.Background(), q, args...,
	).Scan(&id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		zlog.Logger.Error().Err(fmt.Errorf("%s: %w", op, err)).Send()
		return 0
	}

	return id
}

func (p *Postgres) Childs(parentID int64, opts *comment.GetterOpts) ([]comment.Comment, error) {
	const op = "internal.storage.postgres.comments.childs"

	args, q := p.generateAgrigatedOpts(parentID, opts)
	rows, err := p.db.QueryWithRetry(
		context.Background(), retryOpts,
		q, args...,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	defer func() {
		_ = rows.Close()
	}()

	if errors.Is(rows.Err(), sql.ErrNoRows) {
		return nil, errs.ErrDBNotFound
	} else if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	lastID := p.lastID(parentID)
	coms := []comment.Comment{}
	for i := 0; rows.Next(); i++ {
		var tmp comment.Comment
		var parentID sql.NullInt64

		err := rows.Scan(&tmp.ID, &tmp.Message, &parentID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if parentID.Valid {
			tmp.ParentID = parentID.Int64
		}
		if tmp.ID == lastID {
			tmp.HaveNext = false
		} else {
			tmp.HaveNext = true
		}

		coms = append(coms, tmp)
	}

	// If we have load more button after last comment.
	if lastID == 0 && len(coms) > 0 {
		coms[len(coms)-1].HaveNext = false
	}

	return coms, nil
}

func (p *Postgres) DeleteComment(id int64) error {
	const op = "internal.storage.postgres.comments.Delete"

	q := fmt.Sprintf("delete from %s where id = $1", CommentsTable)

	res, err := p.db.ExecContext(context.Background(), q, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if affected, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	} else if affected == 0 {
		return errs.ErrDBNotAffected
	}

	return nil
}
