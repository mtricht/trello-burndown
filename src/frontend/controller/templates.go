package controller

import (
	"html/template"
	"os"
	"path/filepath"
)

var basePath, _ = os.Getwd()
var viewPath = filepath.Join(basePath, "src", "frontend", "view")
var templates = template.Must(template.ParseGlob(
	filepath.Join(viewPath, "*"),
))
