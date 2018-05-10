package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"github.com/swordbeta/trello-burndown/pkg/trello"
)

func refresh(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trello.Run(vars["board"])
	http.Redirect(w, r, viper.GetString("http.baseURL")+"index", 302)
}
