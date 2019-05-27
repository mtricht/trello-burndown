package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mtricht/trello-burndown/pkg/trello"
	"github.com/spf13/viper"
)

func refresh(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trello.Run(vars["board"], "")
	http.Redirect(w, r, viper.GetString("http.baseURL")+"index", 302)
}

func refreshLabel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	trello.Run(vars["board"], vars["label"])
	http.Redirect(w, r, viper.GetString("http.baseURL")+"index", 302)
}
