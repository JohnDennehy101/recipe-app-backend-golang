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

	recipeTableDBStatement := `INSERT INTO recipes (title, description, created_at, updated_at) values($1, $2, $3, $4) RETURNING id`

	ingredientTableDBStatement := `INSERT INTO ingredients (ingredient, created_at, updated_at) values($1,$2,$3) RETURNING id`

	recipeIngredientTableDBStatement := `INSERT INTO recipe_ingredients (recipe_id, ingredient_id, amount, unit, created_at, updated_at) values ($1,$2,$3,$4,$5,$6)`

	recipeInstructionTableDBStatement := `INSERT INTO instructions (text, line, recipe_id, created_at, updated_at) values ($1, $2,$3,$4,$5)`

	var recipeID int
	err := m.DB.QueryRow(
		recipeTableDBStatement,
		res.Title, res.Description, time.Now(), time.Now()).Scan(&recipeID)
	if err != nil {
		return err
	}

	for _, ingredient := range res.Ingredients {
		var ingredientID int
		err := m.DB.QueryRow(ingredientTableDBStatement, ingredient.Name, time.Now(), time.Now()).Scan(&ingredientID)

		if err != nil {
			return err
		}

		_, err = m.DB.ExecContext(ctx, recipeIngredientTableDBStatement, recipeID, ingredientID, ingredient.Amount, ingredient.Unit, time.Now(), time.Now())

		if err != nil {
			return err
		}

	}

	for _, instruction := range res.Instructions {
		_, err := m.DB.ExecContext(ctx, recipeInstructionTableDBStatement, instruction.Text, instruction.Line, recipeID, time.Now(), time.Now())

		if err != nil {
			return err
		}

	}

	return nil
}
