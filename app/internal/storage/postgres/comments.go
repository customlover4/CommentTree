package postgres

import (
	"CommentTree/internal/entities/comment"
	"CommentTree/pkg/errs"
	"context"
	"database/sql"
	"fmt"
	"strconv"
)

func (p *Postgres) Create(c comment.Comment) (int64, error) {
	const op = "internal.storage.postgres.comments.Create"

	q := fmt.Sprintf(
		"insert into %s (message, parent_id) values ($1, $2) returning id;",
		CommentsTable,
	)

	// if we don't have parent id, we will insert NULL to db
	parentID := false
	if c.ParentID > 0 {
		parentID = true
	}
	parentIDArg := sql.NullInt64{Int64: c.ParentID, Valid: parentID}

	var id int64
	err := p.db.Master.QueryRow(q, c.Message, parentIDArg).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (p *Postgres) generateAgrigatedOpts(id int64, opts *comment.GetterOpts) ([]any, string) {
	args := make([]any, 0, 10)
	argsStr := ""

	if id <= 0 && opts == nil {
		return nil, ""
	} else if id > 0 && opts == nil {
		argsStr += " where id = $1 or parent_id = $2"
		args = append(args, id, id)
		return args, argsStr
	} else if opts != nil {
		tmp := " where"
		i := 1

		if id > 0 {
			tmp = fmt.Sprintf(" (id = $%d", i)
			args = append(args, id)
			i++

			if opts.Substr != "" {
				tmp += fmt.Sprintf(" and POSITION($%d IN message) > 0", i)
				args = append(args, opts.Substr)
				i++
			}

			tmp += ") or parent_id = $" + strconv.Itoa(i)
			args = append(args, id)
			i++
		} else {
			if opts.Substr != "" {
				tmp += fmt.Sprintf(" POSITION($%d IN message) > 0", i)
				args = append(args, opts.Substr)
				i++
			}
		}

		argsStr += tmp
		if opts.Page > 0 {
			// 1(PageElement) * 1(opts.Page is first) = 1, wrong offset for first page
			opts.Page--
			limit := PageElements
			offset := PageElements * opts.Page

			argsStr += fmt.Sprintf(" limit %d offset %d", limit, offset)
		}
	}

	return args, argsStr
}

func (p *Postgres) Comments(id int64, opts *comment.GetterOpts) ([]comment.Comment, error) {
	const op = "internal.storage.postgres.comments.Comments"

	q := fmt.Sprintf(`
		select * from %s
	`, CommentsTable)

	var args []any

	argsAppend, qAppend := p.generateAgrigatedOpts(id, opts)
	args = argsAppend
	q += qAppend

	rows, err := p.db.QueryWithRetry(
		context.Background(), retryOpts,
		q, args...,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	coms := make([]comment.Comment, 0, 10)
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
	if len(coms) == 0 {
		return nil, errs.ErrDBNotFound
	}

	return coms, nil
}

func (p *Postgres) Delete(id int64) error {
	const op = "internal.storage.postgres.comments.Delete"

	q := fmt.Sprintf("delete from %s where id = $1", CommentsTable)

	res, err := p.db.ExecContext(context.Background(), q, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if affected, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	} else if affected == 0 {
		// with this specific err we can show upper layer,
		// what we wont log this error and can handle this error
		// like specific case in handlers (503 or 404)
		return errs.ErrDBNotAffected
	}

	return nil
}
