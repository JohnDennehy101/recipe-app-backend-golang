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

func (m *DBModel) List() ([]Recipe, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	recipeListDBStatement := `
	SELECT id, title, description, created_at, updated_at FROM recipes 
	`

	instructionsListDBStatement := `SELECT id, text, line, created_at, updated_at FROM instructions where recipe_id = $1 order by line`

	ingredientsListDBStatement := `
	SELECT i.ingredient, r.amount, r.unit from recipe_ingredients r
	LEFT JOIN ingredients i on (r.ingredient_id = i.id)
	WHERE r.recipe_id = $1
	`

	rows, err := m.DB.QueryContext(ctx, recipeListDBStatement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var recipeList []Recipe
	var instructionList []*Instruction
	var ingredientList []*Ingredient

	for rows.Next() {

		var recipe Recipe

		err := rows.Scan(
			&recipe.ID,
			&recipe.Title,
			&recipe.Description,
			&recipe.CreatedAt,
			&recipe.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		instructionRows, err := m.DB.QueryContext(ctx, instructionsListDBStatement, recipe.ID)

		if err != nil {
			return nil, err
		}

		defer instructionRows.Close()

		for instructionRows.Next() {

			var instruction Instruction

			err = instructionRows.Scan(
				&instruction.ID,
				&instruction.Text,
				&instruction.Line,
				&instruction.CreatedAt,
				&instruction.UpdatedAt,
			)

			if err != nil {
				return nil, err
			}

			instructionList = append(instructionList, &instruction)

		}

		recipe.Instructions = instructionList

		ingredientRows, err := m.DB.QueryContext(ctx, ingredientsListDBStatement, recipe.ID)

		if err != nil {
			return nil, err
		}

		defer ingredientRows.Close()

		for ingredientRows.Next() {

			var ingredient Ingredient

			err = ingredientRows.Scan(
				&ingredient.Name,
				&ingredient.Amount,
				&ingredient.Unit,
			)

			if err != nil {
				return nil, err
			}

			ingredientList = append(ingredientList, &ingredient)

		}

		recipe.Ingredients = ingredientList

		recipeList = append(recipeList, recipe)

	}
	return recipeList, nil
}
