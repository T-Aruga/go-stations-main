package service

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	stmt, err := s.db.PrepareContext(ctx, insert)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	res, err := stmt.ExecContext(ctx, subject, description)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	t := &model.TODO{ID: id}

	stmt, err = s.db.PrepareContext(ctx, confirm)
	if err != nil {
		return nil, err
	}
	err = stmt.QueryRowContext(ctx, t.ID).Scan(&t.Subject, &t.Description, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return t, nil

}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)
	var rows *sql.Rows
	if prevID == 0 {
		rs, err := s.readQueryCtx(ctx, read, prevID, size)
		if err != nil {
			return nil, err
		}
		rows = rs
	} else {
		rs, err := s.readQueryCtx(ctx, readWithID, prevID, size)
		if err != nil {
			return nil, err
		}
		rows = rs
	}
	defer rows.Close()

	todos := make([]*model.TODO, 0, size)
	for rows.Next() {
		todo := &model.TODO{}

		err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
		if err != nil {
			return nil, err
		}

		todos = append(todos, todo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	stmt, err := s.db.PrepareContext(ctx, update)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, subject, description, id)
	if err != nil {
		return nil, err
	}
	updatedRows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if updatedRows == 0 {
		return nil, &model.ErrNotFound{}
	}

	t := &model.TODO{ID: id}

	stmt, err = s.db.PrepareContext(ctx, confirm)
	if err != nil {
		return nil, err
	}
	err = stmt.QueryRowContext(ctx, id).Scan(&t.Subject, &t.Description, &t.CreatedAt, &t.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return t, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`
	if len(ids) == 0 {
		return nil
	}
	delete := fmt.Sprintf(deleteFmt, strings.Repeat(", ?", len(ids)-1))

	stmt, err := s.db.PrepareContext(ctx, delete)
	if err != nil {
		return err
	}
	defer stmt.Close()

	args := make([]interface{}, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}

	res, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return err
	}

	deletedRows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if deletedRows == 0 {
		return &model.ErrNotFound{}
	}

	return nil
}

func (s *TODOService) readQueryCtx(ctx context.Context, readSql string, prevID, size int64) (*sql.Rows, error) {
	stmt, err := s.db.PrepareContext(ctx, readSql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var rows *sql.Rows
	if prevID == 0 {
		rows, err = stmt.QueryContext(ctx, size)
		if err != nil {
			return nil, err
		}
		return rows, nil
	} else {
		rows, err = stmt.QueryContext(ctx, prevID, size)
		if err != nil {
			return nil, err
		}
	}
	return rows, nil
}
