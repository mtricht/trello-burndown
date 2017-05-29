package controller

import (
	"github.com/itsjamie/go-bindata-templates"
	"github.com/swordbeta/trello-burndown/assets"
)

var templates, err = binhtml.New(assets.Asset, assets.AssetDir).LoadDirectory("assets")
