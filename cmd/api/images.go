package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/mnabil1718/blog.mnabil.dev/internal/data"
	"github.com/mnabil1718/blog.mnabil.dev/internal/storage"
	"github.com/mnabil1718/blog.mnabil.dev/internal/validator"
)

func (app *application) uploadImagesHandler(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	defer file.Close()

	v := validator.New()

	image, err := app.storage.Save(file, *fileHeader, true, app.config.Upload.TempPath, v)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrUnsupportedFormat):
			app.badRequestResponse(w, r, err)
		case errors.Is(err, storage.ErrFileRead):
			app.serverErrorResponse(w, r, err)
		case errors.Is(err, storage.ErrFileCreate):
			app.serverErrorResponse(w, r, err)
		case errors.Is(err, storage.ErrInvalidImage):
			app.badRequestResponse(w, r, err)
		case errors.Is(err, storage.ErrValidation):
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Images.Insert(image)
	if err != nil {
		if errors.Is(err, data.ErrDuplicateImageName) {
			v.AddError("name", "name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	image.URL = app.generateImageURL(image.Name)

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/images/%s", image.Name))

	err = app.writeJSON(w, http.StatusCreated, envelope{"image": image}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getImagesHandler(w http.ResponseWriter, r *http.Request) {

	name, err := app.getImageNameFromRequestContext(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	image, err := app.models.Images.GetByName(name)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	path, err := app.storage.GetFullPath(image)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	opts := &storage.ImageProcessingOption{}
	app.readProcessingOptions(w, r, opts)

	v := validator.New()

	if storage.ValidateImageProcessingOption(v, opts); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	img, err := storage.ProcessImage(path, opts)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = storage.EncodeImage(w, r, img, opts, image)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *application) getImagesMetadataHandler(w http.ResponseWriter, r *http.Request) {
	name, err := app.getImageNameFromRequestContext(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	image, err := app.models.Images.GetByName(name)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	image.URL = app.generateImageURL(image.Name)

	err = app.writeJSON(w, http.StatusCreated, envelope{"image": image}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
