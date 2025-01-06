package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/images/:name", app.getImagesHandler)
	router.HandlerFunc(http.MethodGet, "/v1/images/:name/metadata", app.getImagesMetadataHandler)
	router.HandlerFunc(http.MethodPost, "/v1/images", app.uploadImagesHandler)

	return app.recoverPanic(app.enableCORS(app.rateLimit(router)))
}
