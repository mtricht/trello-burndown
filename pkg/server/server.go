package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/itsjamie/go-bindata-templates"
	"github.com/muazamkamal/trello-burndown/assets"
	"github.com/spf13/viper"
)

var templates, err = binhtml.New(assets.Asset, assets.AssetDir).LoadDirectory("assets")

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
		index,
	},
	route{
		"Index page",
		[]string{"GET"},
		"/index",
		index,
	},
	route{
		"Add a trello board",
		[]string{"GET"},
		"/add",
		addGet,
	},
	route{
		"Add a trello board",
		[]string{"POST"},
		"/add",
		addPost,
	},
	route{
		"View a burndown chart",
		[]string{"GET"},
		"/view/{board}",
		view,
	},
	route{
		"Delete a trello board",
		[]string{"GET"},
		"/delete/{board}",
		delete,
	},
	route{
		"Refresh a trello board",
		[]string{"GET"},
		"/refresh/{board}",
		refresh,
	},
}

// Start starts the web server serving the frontend.
func Start() {
	log.Printf("Listening on :%s, open now: http://127.0.0.1:%s", viper.GetString("http.port"), viper.GetString("http.port"))
	router := newRouter()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", viper.GetString("http.port")), router))
}

func newRouter() http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range serverRoutes {

		router.
			Methods(route.Method...).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}
	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	return handlers.CORS(
		handlers.AllowedHeaders([]string{
			"content-type",
		}),
		handlers.AllowCredentials(),
	)(loggedRouter)
}
