package server

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"github.com/mtricht/trello-burndown/pkg/trello"
)

type viewPage struct {
	Board   trello.Board
	Dates   []time.Time
	BaseURL string
}

func view(w http.ResponseWriter, r *http.Request) {
	ViewPage := getViewPage(r)
	err := templates.ExecuteTemplate(w, "view", ViewPage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getViewPage(r *http.Request) *viewPage {
	vars := mux.Vars(r)
	db := trello.GetDatabase()
	defer db.Close()
	board := trello.Board{}
	db.Preload("CardProgress", func(db *gorm.DB) *gorm.DB {
		return db.Order("date ASC")
	}).Where("id = ?", vars["board"]).First(&board)
	return &viewPage{
		Board:   board,
		Dates:   getDatesBetween(board.DateStart, board.DateEnd),
		BaseURL: viper.GetString("http.baseURL"),
	}
}

func getDatesBetween(start time.Time, end time.Time) []time.Time {
	delta := int(end.Sub(start).Hours())
	var dates []time.Time
	for index := 0; index <= delta; index++ {
		date, _ := time.Parse("2006-01-02", start.Format("2006-01-02"))
		date = date.Add(time.Hour * 24 * time.Duration(index))
		delta -= 24
		if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
			continue
		}
		dates = append(dates, date)
	}
	dates = append(dates, end)
	return dates
}
