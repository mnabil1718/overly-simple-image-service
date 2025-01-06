package data

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/mnabil1718/blog.mnabil.dev/internal/validator"
)

var (
	ErrDuplicateImageName = errors.New("duplicate image name")
)

type Image struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Alt       string    `json:"alt"`
	FileName  string    `json:"file_name,omitempty"`
	Size      int32     `json:"size,omitempty"`
	Width     int32     `json:"width,omitempty"`
	Height    int32     `json:"height,omitempty"`
	MIMEType  string    `json:"mime_type,omitempty"`
	URL       string    `json:"url,omitempty"` // will always be empty from DB, remember to set in handlers
	IsTemp    bool      `json:"-"`
	UpdatedAt time.Time `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	Version   int32     `json:"-"`
}

func ValidateImageName(v *validator.Validator, name string) {
	v.Check(name != "", "name", "must be provided")
	v.Check(len(name) <= 750, "name", "must not be more than 750 bytes long")
	v.Check(validator.Matches(name, validator.ImageNameRX), "name", "must be a valid image name")
}

func ValidateImageFileName(v *validator.Validator, fileName string) {
	v.Check(fileName != "", "file_name", "must be provided")
	v.Check(validator.Matches(fileName, validator.ImageFileNameRX), "file_name", "must be a valid image file name")
}

func ValidateImage(v *validator.Validator, image *Image) {
	ValidateImageName(v, image.Name)
	ValidateImageFileName(v, image.FileName)
	v.Check(image.Alt != "", "alt", "must be provided")
	v.Check(len(image.Alt) <= 750, "alt", "must be less than 750 bytes long")
	v.Check(image.Size > 0, "size", "must be more than zero")
	v.Check(image.Size < 10*1024*1024, "size", "must be less than 10 MB")
	v.Check(image.Height > 0, "height", "must be more than zero")
	v.Check(image.Width > 0, "width", "must be more than zero")
	v.Check(v.In(image.MIMEType, "image/jpeg", "image/png", "image/webp", "image/gif"), "mime_type", "must either be .jpeg, .png, .webp, or .gif")
}

type ImageModel struct {
	DB *sql.DB
}

func (model ImageModel) Insert(image *Image) error {
	SQL := `INSERT INTO images (name, alt, file_name, size, width, height, mime_type)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			RETURNING id, created_at, updated_at, version`

	args := []interface{}{image.Name, image.Alt, image.FileName, image.Size, image.Width, image.Height, image.MIMEType}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := model.DB.QueryRowContext(ctx, SQL, args...).Scan(&image.ID, &image.CreatedAt, &image.UpdatedAt, &image.Version)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), `violates unique constraint "images_name_key"`):
			return ErrDuplicateImageName
		default:
			return err
		}
	}

	return nil
}

func (model ImageModel) GetByName(name string) (*Image, error) {

	SQL := `SELECT id, name, alt, file_name, size, width, height, mime_type, created_at, updated_at, version, is_temp
			FROM images WHERE
			name=$1`

	image := &Image{}

	args := []interface{}{name}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := model.DB.QueryRowContext(ctx, SQL, args...).Scan(&image.ID, &image.Name, &image.Alt, &image.FileName, &image.Size, &image.Width, &image.Height, &image.MIMEType, &image.CreatedAt, &image.UpdatedAt, &image.Version, &image.IsTemp)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound

		default:
			return nil, err
		}
	}

	return image, nil
}

func (model ImageModel) Update(image *Image) error {
	SQL := `UPDATE images
	 				SET alt=$1, is_temp=$2, updated_at=$3, version=version + 1
					WHERE id=$4 AND version=$5 
					RETURNING version`

	args := []interface{}{image.Alt, image.IsTemp, image.UpdatedAt, image.ID, image.Version}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := model.DB.QueryRowContext(ctx, SQL, args...).Scan(&image.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}
