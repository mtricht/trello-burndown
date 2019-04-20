package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/mtricht/trello-burndown/pkg/trello"
)

func delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	db := trello.GetDatabase()
	defer db.Close()
	db.Delete(&trello.Board{
		ID: vars["board"],
	})
	http.Redirect(w, r, viper.GetString("http.baseURL")+"index", 302)
}
