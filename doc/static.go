package doc

import (
	"embed"
	"io/fs"
)

//go:embed swagger
var swaggerFs embed.FS

var SwaggerFs fs.FS

func init() {
	var err error
	SwaggerFs, err = fs.Sub(swaggerFs, "swagger")
	if err != nil {
		panic(err)
	}
}
