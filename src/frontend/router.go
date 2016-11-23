package frontend

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

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
