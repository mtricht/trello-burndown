package controller

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/swordbeta/trello-burndown/src/backend"
)

type viewPage struct {
	Board backend.Board
	Dates []time.Time
}

// View renders the burndown chart!
func View(w http.ResponseWriter, r *http.Request) {
	ViewPage := getViewPage(r)
	err := templates.ExecuteTemplate(w, "view", ViewPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getViewPage(r *http.Request) *viewPage {
	vars := mux.Vars(r)
	db := backend.GetDatabase()
	defer db.Close()
	board := backend.Board{}
	db.Preload("CardProgress").Where("id = ?", vars["board"]).First(&board)
	return &viewPage{
		Board: board,
		Dates: getDatesBetween(board.DateStart, board.DateEnd),
	}
}

func getDatesBetween(start time.Time, end time.Time) []time.Time {
	delta := int(end.Sub(start).Hours())
	dates := make([]time.Time, (delta/24)+1)
	index := 1
	dates[0] = start
	for delta != 0 {
		date, _ := time.Parse("2006-01-02", start.Format("2006-01-02"))
		dates[index] = date.Add(time.Hour * 24 * time.Duration(index))
		index++
		delta -= 24
	}
	return dates
}
