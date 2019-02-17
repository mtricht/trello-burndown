package trello

import (
	"log"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
)

// Board contains data of a trello board.
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

// CardProgress represents the progress of a card.
type CardProgress struct {
	gorm.Model
	BoardID string
	Date    time.Time
	Points  float64
}

// GetDatabase returns a sqlite3 database connection.
func GetDatabase() *gorm.DB {
	db, err := gorm.Open(viper.GetString("database.dialect"), viper.GetString("database.url"))
	if err != nil {
		log.Fatalln(err)
	}
	db.AutoMigrate(&Board{}, &CardProgress{})
	return db
}

func saveToDatabase(board Board, m map[string]float64) {
	db := GetDatabase()
	defer db.Close()
	oldBoard := Board{}
	db.Where("id = ?", board.ID).First(&oldBoard)
	db.Model(oldBoard).Updates(&board)
	db.Unscoped().Where("board_id = ?", board.ID).Delete(CardProgress{})
	pointsInWeekend := 0.0
	for date, points := range m {
		date, _ := time.Parse("2006-01-02", date)
		if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
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
