package controller

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/swordbeta/trello-burndown/src/watcher"
)

// Delete deletes a trello board!
func Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	db := watcher.GetDatabase()
	defer db.Close()
	db.Delete(&watcher.Board{
		ID: vars["board"],
	})
	http.Redirect(w, r, viper.GetString("http.baseURL")+"index", 302)
}
