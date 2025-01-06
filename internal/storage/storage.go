package storage

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/mnabil1718/blog.mnabil.dev/internal/data"
	"github.com/mnabil1718/blog.mnabil.dev/internal/utils"
	"github.com/mnabil1718/blog.mnabil.dev/internal/validator"
)

type ImageStorage struct {
	path     string
	tempPath string
}

func New(path, tempPath string) (*ImageStorage, error) {

	for _, item := range []string{path, tempPath} {
		if err := os.MkdirAll(item, 0755); err != nil { // Set directory permissions to 0755
			return nil, fmt.Errorf("failed to create directory %s: %w", item, err)
		}
	}

	return &ImageStorage{path, tempPath}, nil
}

func (s *ImageStorage) GetFullPath(image *data.Image) (string, error) {
	var basePath string

	if image.IsTemp {
		basePath = s.tempPath
	} else {
		basePath = s.path
	}

	if basePath == "" || image.FileName == "" {
		return "", errors.New("invalid path or filename")
	}

	return filepath.Join(basePath, image.FileName), nil
}

func (s *ImageStorage) Save(file multipart.File, fileHeader multipart.FileHeader, isTemp bool, path string, v *validator.Validator) (*data.Image, error) {
	// Step 1: Determine MIME type and validate support
	mimeType, err := detectMimeType(file)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrUnsupportedFormat, err)
	}

	// Step 2: Decode image metadata
	width, height, err := decodeImageMetadata(file, mimeType)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidImage, err)
	}

	// Step 3: Generate file name and extension
	name := utils.GenerateImageName(fileHeader.Filename)
	extension := EXT_MAP[mimeType]
	if extension == "" {
		return nil, ErrUnsupportedFormat
	}
	filename := name + extension

	// Step 4: Prepare the image metadata
	image := &data.Image{
		Name:     name,
		Alt:      name,
		FileName: filename,
		Size:     int32(fileHeader.Size),
		Width:    int32(width),
		Height:   int32(height),
		MIMEType: mimeType,
		IsTemp:   isTemp,
	}

	// Step 5: Validate image metadata
	if err := validateImage(v, image); err != nil {
		return nil, err
	}

	// Step 6: Save the file to disk
	if err := saveFileToDisk(file, filepath.Join(path, filename)); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFileCreate, err)
	}

	return image, nil
}

func (s *ImageStorage) Move(source, destination string) error {
	// Open the source file
	src, err := os.Open(source)
	if err != nil {
		return ErrSystem
	}
	defer src.Close()

	// Get file information for permissions
	fi, err := src.Stat()
	if err != nil {
		return ErrSystem
	}

	// Prepare the destination file
	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	perm := fi.Mode() & os.ModePerm
	dst, err := os.OpenFile(destination, flag, perm)
	if err != nil {
		return ErrSystem
	}

	// Copy data from source to destination
	_, err = io.Copy(dst, src)
	if err != nil {
		dst.Close()
		os.Remove(destination)
		return ErrFileMove
	}

	// Close the destination file explicitly
	if err := dst.Close(); err != nil {
		return ErrSystem
	}

	// Ensure the source file is closed
	if err := src.Close(); err != nil {
		return ErrSystem
	}

	// Remove the source file
	if err := os.Remove(source); err != nil {
		return ErrSystem
	}

	return nil
}
