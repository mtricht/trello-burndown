package controller

import (
	"net/http"

	"github.com/spf13/viper"
	"github.com/swordbeta/trello-burndown/src/backend"
)

type indexPage struct {
	Boards  []backend.Board
	BaseURL string
}

// Index renders the index page.
func Index(w http.ResponseWriter, r *http.Request) {
	db := backend.GetDatabase()
	defer db.Close()
	boards := []backend.Board{}
	db.Order("date_start desc").Find(&boards)
	indexPage := indexPage{
		Boards:  boards,
		BaseURL: viper.GetString("http.baseURL"),
	}
	err := templates.ExecuteTemplate(w, "index", indexPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
