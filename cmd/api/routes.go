package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/status", app.statusHandler)
	router.HandlerFunc(http.MethodPost, "/recipe", app.CreateRecipe)
	router.HandlerFunc(http.MethodGet, "/recipe", app.GetAllRecipes)

	return router
}
