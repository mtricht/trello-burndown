package webserver

import (
	"net/http"

	"github.com/swordbeta/trello-burndown/src/webserver/controller"
)

type route struct {
	Name        string
	Method      []string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routes []route

var serverRoutes = routes{
	route{
		"Index page",
		[]string{"GET"},
		"/",
		controller.Index,
	},
	route{
		"Index page",
		[]string{"GET"},
		"/index",
		controller.Index,
	},
	route{
		"Add a trello board",
		[]string{"GET"},
		"/add",
		controller.AddGet,
	},
	route{
		"Add a trello board",
		[]string{"POST"},
		"/add",
		controller.AddPost,
	},
	route{
		"View a burndown chart",
		[]string{"GET"},
		"/view/{board}",
		controller.View,
	},
	route{
		"Delete a trello board",
		[]string{"GET"},
		"/delete/{board}",
		controller.Delete,
	},
	route{
		"Refresh a trello board",
		[]string{"GET"},
		"/refresh/{board}",
		controller.Refresh,
	},
}
