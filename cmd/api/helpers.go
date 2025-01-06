package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/julienschmidt/httprouter"
	"github.com/mnabil1718/blog.mnabil.dev/internal/config"
	"github.com/mnabil1718/blog.mnabil.dev/internal/storage"
	"github.com/mnabil1718/blog.mnabil.dev/internal/utils"
	"github.com/mnabil1718/blog.mnabil.dev/internal/validator"
)

func openDB(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.DB.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)
	duration, err := time.ParseDuration(cfg.DB.MaxIdleTime) // "15m" or "5s"
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}

func (app *application) readProcessingOptions(w http.ResponseWriter, r *http.Request, opts *storage.ImageProcessingOption) {
	queryString := r.URL.Query()
	v := validator.New()

	opts.Crop = app.readBool(queryString, "crop", v)
	opts.Width = app.readInt(queryString, "w", 0, v)
	opts.Height = app.readInt(queryString, "h", 0, v)
	opts.Quality = app.readInt(queryString, "q", 100, v)
	opts.BlurSigma = app.readFloat(queryString, "blur", 0, v)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
}

func (app *application) getImageNameFromRequestContext(request *http.Request) (string, error) {
	params := httprouter.ParamsFromContext(request.Context())
	name := params.ByName("name")

	err := utils.ValidateImageName(name)
	if err != nil {
		return "", err
	}

	return name, nil
}

type envelope map[string]interface{}

func (app *application) writeJSON(writer http.ResponseWriter, code int, data envelope, headers http.Header) error {
	resp, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp = append(resp, '\n')

	for key, value := range headers {
		writer.Header()[key] = value
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	writer.Write(resp)

	return nil
}

func (app *application) readBool(queryString url.Values, key string, v *validator.Validator) bool {

	value := queryString.Get(key)

	if value == "" {
		return false
	}

	res, err := strconv.ParseBool(value)
	if err != nil {
		v.AddError(key, fmt.Sprintf("%s must be a boolean value.", key))
	}

	return res
}

func (app *application) readInt(queryString url.Values, key string, defaultValue int, v *validator.Validator) int {
	value := queryString.Get(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		v.AddError(key, fmt.Sprintf("%s must be an integer value.", key))
	}

	return intValue
}

func (app *application) readFloat(queryString url.Values, key string, defaultValue float64, v *validator.Validator) float64 {
	value := queryString.Get(key)
	if value == "" {
		return defaultValue
	}

	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		v.AddError(key, fmt.Sprintf("%s must be a float value.", key))
	}

	return floatValue
}

func (application *application) generateImageURL(name string) string {
	return fmt.Sprintf("http://%s:%d/v1/images/%s", application.config.Host, application.config.Port, name)
}
