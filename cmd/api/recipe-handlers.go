package main

import (
	"net/http"
	"time"

	"github.com/JohnDennehy101/recipe-app-backend-golang/models"
)

func (app *application) CreateRecipe(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.logger.Fatal(err)
	}

	title := r.FormValue("title")
	description := r.FormValue("description")

	recipe := models.Recipe{
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = app.models.DB.CreateRecipe(recipe)

	if err != nil {
		app.logger.Println(err)
	}
}
