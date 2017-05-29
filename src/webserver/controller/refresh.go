package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/swordbeta/trello-burndown/src/watcher"
)

// Refresh refreshes a trello board!
func Refresh(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	watcher.Run(vars["board"])
	http.Redirect(w, r, viper.GetString("http.baseURL")+"index", 302)
}
