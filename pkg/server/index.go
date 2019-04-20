package server

import (
	"net/http"

	"github.com/mtricht/trello-burndown/pkg/trello"
	"github.com/spf13/viper"
)

type indexPage struct {
	Boards   []trello.Board
	BaseURL  string
	ReadOnly bool
}

func index(w http.ResponseWriter, r *http.Request) {
	db := trello.GetDatabase()
	defer db.Close()
	boards := []trello.Board{}
	db.Order("date_start desc").Find(&boards)
	indexPage := indexPage{
		Boards:   boards,
		BaseURL:  viper.GetString("http.baseURL"),
		ReadOnly: viper.GetBool("http.readOnly"),
	}
	err := templates.ExecuteTemplate(w, "index", indexPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
