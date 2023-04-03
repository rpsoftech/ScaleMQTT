package global

import (
	"os"
	"path/filepath"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
)

var JWTKEY []byte
var Logger *zerolog.Logger
var Validator = validator.New()

func GetCuurentPath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}
