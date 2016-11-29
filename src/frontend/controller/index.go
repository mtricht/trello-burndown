package controller

import (
	"net/http"

	"github.com/swordbeta/trello-burndown/src/backend"
)

// Index renders the index page.
func Index(w http.ResponseWriter, r *http.Request) {
	db := backend.GetDatabase()
	defer db.Close()
	boards := []backend.Board{}
	db.Order("date_start desc").Find(&boards)
	err = templates.ExecuteTemplate(w, "index", boards)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
