package main

import (
	"github.com/julienschmidt/httprouter"
)

func (a *application) routes() *httprouter.Router {
	// create the router
	router := httprouter.New()

	return router
}
