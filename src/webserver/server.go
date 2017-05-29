package webserver

import (
	"fmt"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

// Start starts the webserver serving the frontend.
func Start() {
	log.Printf("Listening on http://127.0.0.1:%s", viper.GetString("http.port"))
	router := newRouter()
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", viper.GetString("http.port")), router))
}
