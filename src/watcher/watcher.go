package watcher

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/spf13/viper"
)

type board struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	URL             string `json:"url"`
	Lists           []list `json:"lists"`
	Cards           []card `json:"cards"`
	CardsCompleted  uint
	CardsTotal      uint
	PointsCompleted float64
	PointsTotal     float64
}

type list struct {
	ID       string `json:"id"`
	Position int    `json:"pos"`
}

type card struct {
	ID               string `json:"id"`
	ListID           string `json:"idList"`
	Name             string `json:"name"`
	DateLastActivity string `json:"dateLastActivity"`
	Actions          *[]action
}

type action struct {
	Data struct {
		ListAfter struct {
			ListID string `json:"id"`
		}
		ListBefore struct {
			ListID string `json:"id"`
		}
	}
	Date string `json:"date"`
}

type cardResult struct {
	Error       error
	Date        string
	Complete    bool
	Points      float64
	TrelloError bool
}

// Start starts watching boards that are active. Refreshes according
// to the refresh rate set in the configuration.
func Start() {
	go runBoards()
	ch := gocron.Start()
	refreshRate := uint64(viper.GetInt64("trello.refreshRate"))
	gocron.Every(refreshRate).Minutes().Do(runBoards)
	<-ch
}

func runBoards() {
	db := GetDatabase()
	defer db.Close()
	boards := []board{}
	yesterday := time.Now().Add(-24 * time.Hour)
	db.Select("id").Where("date_start < ? AND date_end > ?", yesterday, yesterday).Find(&boards)
	for _, board := range boards {
		go Run(board.ID)
	}
}

// Run fetches and saves the points of a given board. Called by
// the watcher and when refreshed on the frontend.
func Run(boardID string) {
	log.Printf("Checking board ID #%s", boardID)
	board, err := getBoard(boardID)
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}
	if board == nil {
		log.Println("Something went wrong requesting a board from Trello.")
		return
	}
	board.ID = boardID
	log.Printf("Board name: %s", board.Name)
	lastListID := getLastList(board)
	resultChannel := make(chan *cardResult)
	for _, card := range board.Cards {
		go determineCardComplete(card, lastListID, resultChannel)
	}
	var m = make(map[string]float64)
	for i := 0; i < len(board.Cards); i++ {
		response := <-resultChannel
		if response.Error != nil {
			log.Fatalln(response.Error)
		}
		if response.TrelloError {
			log.Println("Something went wrong requesting a card from Trello.")
			return
		}
		if response.Complete {
			board.CardsCompleted++
			board.PointsCompleted += response.Points
			if _, ok := m[response.Date]; ok {
				m[response.Date] = response.Points + m[response.Date]
			} else {
				m[response.Date] = response.Points
			}
		}
		board.CardsTotal++
		board.PointsTotal += response.Points
	}
	log.Printf("Cards progress: %d/%d", board.CardsCompleted, board.CardsTotal)
	log.Printf("Total points: %f/%f", board.PointsCompleted, board.PointsTotal)
	saveToDatabase(board, m)
}

func getBoard(id string) (*board, error) {
	url := fmt.Sprintf(
		"https://api.trello.com/1/boards/%s?key=%s&token=%s&cards=visible&lists=all",
		id,
		viper.GetString("trello.apiKey"),
		viper.GetString("trello.userToken"),
	)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.Header.Get("Content-Type") != "application/json; charset=utf-8" {
		return nil, nil
	}
	r := new(board)
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func getLastList(board *board) string {
	var highestPos int
	var listID string
	for _, list := range board.Lists {
		if list.Position > highestPos {
			listID = list.ID
		}
	}
	return listID
}

func determineCardComplete(card card, listID string, res chan *cardResult) {
	points := getPoints(&card)
	if card.ListID != listID {
		res <- &cardResult{
			Complete: false,
			Points:   points,
		}
		return
	}
	url := fmt.Sprintf(
		"https://api.trello.com/1/cards/%s/actions?key=%s&token=%s",
		card.ID,
		viper.GetString("trello.apiKey"),
		viper.GetString("trello.userToken"),
	)
	resp, err := http.Get(url)
	if err != nil {
		res <- &cardResult{
			Error: err,
		}
		return
	}
	defer resp.Body.Close()
	if resp.Header.Get("Content-Type") != "application/json; charset=utf-8" {
		res <- &cardResult{
			TrelloError: true,
		}
		return
	}
	var actions []action
	err = json.NewDecoder(resp.Body).Decode(&actions)
	if err != nil {
		res <- &cardResult{
			Error: err,
		}
		return
	}
	date, err := time.Parse(time.RFC3339Nano, card.DateLastActivity)
	if err != nil {
		res <- &cardResult{
			Error: err,
		}
		return
	}
	for _, action := range actions {
		if action.Data.ListAfter.ListID != action.Data.ListBefore.ListID &&
			action.Data.ListAfter.ListID == listID {
			date, err = time.Parse(time.RFC3339Nano, action.Date)
			if err != nil {
				res <- &cardResult{
					Error: err,
				}
				return
			}
			break
		}
	}
	res <- &cardResult{
		Complete: true,
		Date:     date.Format("2006-01-02"),
		Points:   points,
	}
}

func getPoints(card *card) float64 {
	r := regexp.MustCompile(`\(([0-9]*\.[0-9]+|[0-9]+)\)`)
	matches := r.FindStringSubmatch(card.Name)
	if len(matches) != 2 {
		return 0
	}
	points, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		log.Fatalln(err)
	}
	return points
}
