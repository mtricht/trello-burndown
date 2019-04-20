package server

import (
	"net/http"

	"github.com/spf13/viper"
	"github.com/mtricht/trello-burndown/pkg/trello"
)

type indexPage struct {
	Boards  []trello.Board
	BaseURL string
}

func index(w http.ResponseWriter, r *http.Request) {
	db := trello.GetDatabase()
	defer db.Close()
	boards := []trello.Board{}
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
