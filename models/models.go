package models

import (
	"database/sql"
	"time"
)

type Models struct {
	DB DBModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

type Recipe struct {
	ID           int            `json:"id"`
	Title        string         `json:"title"`
	Description  string         `json:"description"`
	Ingredients  []*Ingredient  `json:"ingredients"`
	Instructions []*Instruction `json:"instruction"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type Ingredient struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Amount    float64   `json:"amount"`
	Unit      string    `json:"unit"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Instruction struct {
	ID        int       `json:"id"`
	Text      string    `json:"text"`
	Line      int       `json:"line"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
