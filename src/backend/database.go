package backend

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/swordbeta/trello-burndown/src/util"
)

type Board struct {
	ID              string `gorm:"primary_key"`
	Name            string
	DateStart       time.Time
	DateEnd         time.Time
	Cards           uint
	Points          float64
	CardsCompleted  uint
	PointsCompleted float64
	CardProgress    []CardProgress
}

type CardProgress struct {
	gorm.Model
	BoardID string
	Date    time.Time
	Points  float64
}

func GetDatabase() *gorm.DB {
	db, err := gorm.Open("sqlite3", "./trello.db")
	if err != nil {
		log.Fatalln(err)
	}
	db.AutoMigrate(&Board{}, &CardProgress{})
	return db
}

func saveToDatabase(board *board, m map[string]float64) {
	db := GetDatabase()
	defer db.Close()
	oldBoard := Board{}
	db.Where("id = ?", board.ID).First(&oldBoard)
	db.Model(oldBoard).Updates(&Board{
		Name:            board.Name,
		Cards:           board.CardsTotal,
		Points:          board.PointsTotal,
		CardsCompleted:  board.CardsCompleted,
		PointsCompleted: board.PointsCompleted,
	})
	db.Delete(&CardProgress{
		BoardID: board.ID,
	})
	pointsInWeekend := 0.0
	for date, points := range m {
		date, _ := time.Parse("2006-01-02", date)
		if util.IsWeekend(date) {
			pointsInWeekend += points
			continue
		}
		db.Save(&CardProgress{
			Date:    date,
			Points:  points + pointsInWeekend,
			BoardID: board.ID,
		})
		pointsInWeekend = 0
	}
}
