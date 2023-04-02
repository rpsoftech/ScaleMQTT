package global

import (
	"os"
	"path/filepath"
)

var JWTKEY string

func GetCuurentPath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(ex)
}
