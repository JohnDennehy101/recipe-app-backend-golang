package models

import (
	"context"
	"database/sql"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

func (m *DBModel) Get(id int) (*Recipe, error) {
	return nil, nil
}

func (m *DBModel) CreateRecipe(res Recipe) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	dbStatement := `INSERT INTO recipes (title, description, created_at, updated_at) values($1, $2, $3, $4)`

	_, err := m.DB.ExecContext(ctx, dbStatement, res.Title, res.Description, res.CreatedAt, res.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}
