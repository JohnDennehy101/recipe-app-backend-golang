package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/JohnDennehy101/recipe-app-backend-golang/models"
)

func (app *application) CreateRecipe(w http.ResponseWriter, r *http.Request) {

	var recipe models.Recipe

	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		app.logger.Fatal(err)
	}

	if err := json.Unmarshal(requestBody, &recipe); err != nil {
		app.logger.Fatal(err)
	}

	if err != nil {
		app.logger.Fatal(err)
	}

	err = app.models.DB.CreateRecipe(recipe)

	if err != nil {
		app.logger.Println(err)
	}
}

func (app *application) GetAllRecipes(w http.ResponseWriter, r *http.Request) {

	recipeList, err := app.models.DB.List()

	if err != nil {
		app.logger.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":     "success",
		"statusCode": 200,
		"data":       recipeList,
	})
}
