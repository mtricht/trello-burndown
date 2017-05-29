package controller

import (
	"net/http"

	"github.com/spf13/viper"
	"github.com/swordbeta/trello-burndown/src/watcher"
)

type indexPage struct {
	Boards  []watcher.Board
	BaseURL string
}

// Index renders the index page.
func Index(w http.ResponseWriter, r *http.Request) {
	db := watcher.GetDatabase()
	defer db.Close()
	boards := []watcher.Board{}
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
