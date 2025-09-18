package postgres

import (
	"CommentTree/internal/entities/comment"
	"CommentTree/pkg/errs"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (p *Postgres) CreateComment(c comment.Comment) (int64, error) {
	const op = "internal.storage.postgres.comments.Create"

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
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (p *Postgres) generateAgrigatedOpts(parentID int64, opts *comment.GetterOpts) ([]any, string) {
	if opts != nil && opts.Empty() {
		opts = nil
	}

	args := []any{}
	argsStr := fmt.Sprintf(`
		select * from %s
	`, CommentsTable)

	if opts == nil {
		if parentID <= 0 {
			argsStr += " where parent_id is NULL"
			return args, argsStr
		} else {
			argsStr += " where parent_id = $1"
			args = append(args, parentID)
			return args, argsStr
		}
	}

	argIDX := 1

	if parentID > 0 {
		argsStr += fmt.Sprintf(" where parent_id = $%d", argIDX)
		args = append(args, parentID)
		argIDX++
	} else {
		argsStr += " where parent_id is NULL"
	}

	if opts.Substr != "" {
		argsStr += fmt.Sprintf(" and POSITION($%d IN message) > 0", argIDX)
		args = append(args, opts.Substr)
	}

	if opts.Page > 0 {
		// 1(PageElement) * 1(opts.Page is first) = 1, wrong offset for first page
		opts.Page--
		limit := PageElements
		offset := PageElements * opts.Page

		argsStr += fmt.Sprintf(" limit %d offset %d", limit, offset)
	}

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

	coms := []comment.Comment{}
	for rows.Next() {
		var tmp comment.Comment
		var parentID sql.NullInt64

		err := rows.Scan(&tmp.ID, &tmp.Message, &parentID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		if parentID.Valid {
			tmp.ParentID = parentID.Int64
		}

		coms = append(coms, tmp)
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
