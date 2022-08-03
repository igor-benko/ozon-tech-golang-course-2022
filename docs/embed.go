package docs

import (
	"embed"
)

const SwaggerFileName = "person.swagger.json"

//go:embed person.swagger.json
var SwaggerFile embed.FS
