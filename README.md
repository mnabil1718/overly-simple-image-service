# Overly Simple Image Service

Written in go, only supports image storage and processing endpoint

## Dependencies

- [https://github.com/chai2010/webp]
- [https://github.com/disintegration/imaging]
- [https://github.com/air-verse/air]
- [https://github.com/spf13/viper]
- [https://github.com/julienschmidt/httprouter]

## Available Endpoints

| Method     | Endpoint |
| ----------- | ----------- |
| GET     | /v1/images/:name?optional-params      |
| POST   | /v1/images        |

## Upload

Requires multi-part form data with key `file`

## Suported image formats

- image/jpeg
- image/png
- image/gif
- image/tiff
- image/webp
- image/bmp

## Supported Processing Operations

| Param       | Type        |  Description                                                                               |
| ----------- | ----------- | ------------------------------------------------------------------------------------------ |
| `w`         | int         |  Specify resized width of the image. If height not specified retain original aspect ratio. |
| `h`         | int         |  Specify resized height of the image. If width not specified retain original aspect ratio. |
| `crop`      | bool        |  If true height and width have to be specified. Resize and Crop image based on `w` and `h` |
| `blur`      | float64     |  Specify gaussian blur filter on image. Applied last after crop and resize.                |
| `quality`   | int         |  Specify quality of image upon encoding. Only works with Lossy (jpeg, webp)                |

