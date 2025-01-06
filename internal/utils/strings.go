package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/mnabil1718/blog.mnabil.dev/internal/validator"
)

func Slugify(text string) string {
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, " ", "-")
	text = strings.ReplaceAll(text, "_", "-")
	reg := regexp.MustCompile("[^a-z0-9-]+")
	text = reg.ReplaceAllString(text, "")
	return text
}

func GenerateImageName(fileName string) string {
	id := uuid.New().String()

	fileName = strings.Split(fileName, ".")[0]
	slugifiedName := Slugify(fileName)
	timestamp := time.Now().Format("20060102_150405")
	return fmt.Sprintf("%s-%s-%s", slugifiedName, id, timestamp)
}

func ValidateImageName(name string) error {

	if !validator.ImageNameRX.MatchString(name) {
		return errors.New("invalid image name")
	}

	return nil
}
