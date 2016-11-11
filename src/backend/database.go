package backend

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
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
	db.Save(&Board{
		ID:              board.ID,
		Name:            board.Name,
		Cards:           board.CardsTotal,
		Points:          board.PointsTotal,
		CardsCompleted:  board.CardsCompleted,
		PointsCompleted: board.PointsCompleted,
	})
	db.Delete(&CardProgress{
		BoardID: board.ID,
	})
	for date, points := range m {
		date, _ := time.Parse("2006-01-02", date)
		db.Save(&CardProgress{
			Date:   date,
			Points: points,
		})
	}
}
