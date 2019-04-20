package server

import (
	"net/http"
	"time"

	"github.com/mtricht/trello-burndown/pkg/trello"
	"github.com/spf13/viper"
)

type addPage struct {
	BaseURL string
}

func addGet(w http.ResponseWriter, r *http.Request) {
	if viper.GetBool("http.readOnly") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err := templates.ExecuteTemplate(w, "add", addPage{
		BaseURL: viper.GetString("http.baseURL"),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func addPost(w http.ResponseWriter, r *http.Request) {
	if viper.GetBool("http.readOnly") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	r.ParseForm()
	db := trello.GetDatabase()
	defer db.Close()
	startDate, _ := time.Parse("2006-01-02", r.FormValue("start_date"))
	endDate, _ := time.Parse("2006-01-02", r.FormValue("end_date"))
	db.Save(&trello.Board{
		ID:        r.FormValue("id"),
		DateStart: startDate,
		DateEnd:   endDate,
	})
	trello.Run(r.FormValue("id"))
	http.Redirect(w, r, viper.GetString("http.baseURL")+"index", 302)
}
