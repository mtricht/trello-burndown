package controller

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/swordbeta/trello-burndown/src/backend"
	"github.com/swordbeta/trello-burndown/src/util"
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
	db.Preload("CardProgress", func(db *gorm.DB) *gorm.DB {
		return db.Order("date ASC")
	}).Where("id = ?", vars["board"]).First(&board)
	return &viewPage{
		Board: board,
		Dates: getDatesBetween(board.DateStart, board.DateEnd),
	}
}

func getDatesBetween(start time.Time, end time.Time) []time.Time {
	delta := int(end.Sub(start).Hours())
	var dates []time.Time
	dates = append(dates, start)
	for index := 0; index <= delta; index++ {
		date, _ := time.Parse("2006-01-02", start.Format("2006-01-02"))
		date = date.Add(time.Hour * 24 * time.Duration(index))
		delta -= 24
		if util.IsWeekend(date) {
			continue
		}
		dates = append(dates, date)
	}
	dates = append(dates, end)
	return dates
}
