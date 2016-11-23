package controller

import (
	"net/http"
	"time"

	"github.com/swordbeta/trello-burndown/src/backend"
)

// AddGet renders the form to add a trello board.
func AddGet(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "add", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// AddPost adds the new trello board to the SQLite database!
func AddPost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	db := backend.GetDatabase()
	defer db.Close()
	startDate, _ := time.Parse("2006-01-02", r.FormValue("start_date"))
	endDate, _ := time.Parse("2006-01-02", r.FormValue("end_date"))
	db.Save(&backend.Board{
		ID:        r.FormValue("id"),
		DateStart: startDate,
		DateEnd:   endDate,
	})
	backend.Run(r.FormValue("id"))
	http.Redirect(w, r, "index", 301)
}
