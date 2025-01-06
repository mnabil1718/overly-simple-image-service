package storage

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"strings"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/mnabil1718/blog.mnabil.dev/internal/data"
	"github.com/mnabil1718/blog.mnabil.dev/internal/validator"
	"golang.org/x/image/bmp"
	"golang.org/x/image/tiff"
)

const MAX_IMAGE_DIM = 6000

type ImageProcessingOption struct {
	Width     int
	Height    int
	Crop      bool
	BlurSigma float64
	Quality   int
}

func ValidateImageProcessingOption(v *validator.Validator, opts *ImageProcessingOption) {

	if opts.Width < 50 && opts.Width > 0 {
		v.AddError("width", "have to be atleast 50 pixels wide")
	}

	if opts.Height < 50 && opts.Height > 0 {
		v.AddError("height", "have to be atleast 50 pixels wide")
	}

	v.Check(opts.Quality >= 0, "quality", "cannot be less than 0")
	v.Check(opts.Quality <= 100, "quality", "cannot be more than 100")
	v.Check(opts.BlurSigma >= 0, "blur", "cannot be less than 0")
	v.Check(opts.BlurSigma <= 10, "blur", "cannot be more than 10")
	v.Check(opts.Width <= MAX_IMAGE_DIM, "width", fmt.Sprintf("cannot be more than %d pixels wide", MAX_IMAGE_DIM))
	v.Check(opts.Height <= MAX_IMAGE_DIM, "height", fmt.Sprintf("cannot be more than %d pixels tall", MAX_IMAGE_DIM))

	if opts.Crop {
		v.Check(opts.Width >= 50, "width", "have to be atleast 50 pixels wide")
		v.Check(opts.Height >= 50, "height", "have to be atleast 50 pixels tall")
	}

	if !opts.Crop {
		if opts.Width <= 0 && opts.Height <= 0 {
			v.AddError("width", "cannot be empty")
			v.AddError("height", "cannot be empty")
		}
	}
}

func ProcessImage(path string, opts *ImageProcessingOption) (image.Image, error) {
	img, err := imaging.Open(path)
	if err != nil {
		return nil, ErrOpenImage
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	WValid := opts.Width >= 50
	HValid := opts.Height >= 50
	WFitBound := opts.Width < width
	HFitBound := opts.Height < height

	WHValid := WValid && HValid
	WHFitBounds := WFitBound && HFitBound
	WORHValid := WValid || HValid

	if opts.Crop {
		if WHFitBounds && WHValid {
			img = imaging.Fill(img, opts.Width, opts.Height, imaging.Center, imaging.Lanczos)
		}
	}

	if !opts.Crop {
		if WHFitBounds && WORHValid {
			img = imaging.Resize(img, opts.Width, opts.Height, imaging.Lanczos)
		}
	}

	if opts.BlurSigma > 0 {
		return imaging.Blur(img, opts.BlurSigma), nil
	}

	return img, nil

}

func EncodeImage(w http.ResponseWriter, r *http.Request, img image.Image, opts *ImageProcessingOption, image *data.Image) error {

	accept := r.Header.Get("Accept")
	isWEBPSupported := strings.Contains(accept, "image/webp")

	SetImageHeaders(w, image.Name+EXT_MAP[image.MIMEType], image.MIMEType)

	switch image.MIMEType {
	case "image/jpeg":
		if isWEBPSupported {
			SetImageHeaders(w, image.Name+EXT_MAP["image/webp"], "image/webp")
			return webp.Encode(w, img, &webp.Options{Quality: float32(opts.Quality)})
		}
		return jpeg.Encode(w, img, &jpeg.Options{Quality: opts.Quality})

	case "image/png":
		if isWEBPSupported {
			SetImageHeaders(w, image.Name+EXT_MAP["image/webp"], "image/webp")
			return webp.Encode(w, img, &webp.Options{Quality: float32(opts.Quality)})
		}
		return png.Encode(w, img)

	case "image/gif":
		return gif.Encode(w, img, nil)

	case "image/tiff":
		if isWEBPSupported {
			SetImageHeaders(w, image.Name+EXT_MAP["image/webp"], "image/webp")
			return webp.Encode(w, img, &webp.Options{Quality: float32(opts.Quality)})
		}
		return tiff.Encode(w, img, nil)

	case "image/bmp":
		if isWEBPSupported {
			SetImageHeaders(w, image.Name+EXT_MAP["image/webp"], "image/webp")
			return webp.Encode(w, img, &webp.Options{Quality: float32(opts.Quality)})
		}
		return bmp.Encode(w, img)

	case "image/webp":
		return webp.Encode(w, img, &webp.Options{Quality: float32(opts.Quality)})

	default:
		return ErrUnsupportedFormat
	}
}
