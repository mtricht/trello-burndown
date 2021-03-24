package trello

import (
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/adlio/trello"
	"github.com/jasonlvhit/gocron"
	"github.com/spf13/viper"
)

type cardResult struct {
	Error       error
	Date        string
	Created     string
	Complete    bool
	Points      float64
	TrelloError bool
}

var client *trello.Client

// Start starts watching boards that are active. Refreshes according
// to the refresh rate set in the configuration.
func Start() {
	client = trello.NewClient(
		viper.GetString("trello.apiKey"),
		viper.GetString("trello.userToken"),
	)
	go runBoards()
	ch := gocron.Start()
	refreshRate := uint64(viper.GetInt64("trello.refreshRate"))
	gocron.Every(refreshRate).Minutes().Do(runBoards)
	<-ch
}

func runBoards() {
	db := GetDatabase()
	defer db.Close()
	boards := []Board{}
	yesterday := time.Now().Add(-24 * time.Hour)
	db.Select("id").Where("date_start < ? AND date_end > ?", yesterday, yesterday).Find(&boards)
	for _, board := range boards {
		go Run(board.ID)
	}
}

// Run fetches and saves the points of a given board.
func Run(boardID string) {
	log.Printf("Checking board ID %s", boardID)
	board, err := client.GetBoard(boardID, trello.Defaults())
	if err != nil {
		log.Printf("Couldn't fetch board: %s", err)
		return
	}
	log.Printf("Board name: %s", board.Name)
	lastListID, err := getLastList(board)
	if err != nil {
		log.Printf("Couldn't fetch last list: %s", err)
	}
	resultChannel := make(chan *cardResult)
	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		log.Printf("Couldn't fetch cards: %s", err)
	}
	for _, card := range cards {
		go getCardDetails(card, lastListID, resultChannel)
	}
	boardEntity := Board{
		ID:   boardID,
		Name: board.Name,
	}
	var pointsPerDay = make(map[string]float64)
	var targetsPerDay = make(map[string]float64)
	for i := 0; i < len(cards); i++ {
		response := <-resultChannel
		if response.Error != nil {
			log.Fatalln(response.Error)
		}

		// Total points exist in the sprint on each day (regardless of completion)
		if _, ok := targetsPerDay[response.Created]; ok {
			targetsPerDay[response.Created] = response.Points + targetsPerDay[response.Created]
		} else {
			targetsPerDay[response.Created] = response.Points
		}

		if response.Complete {
			boardEntity.CardsCompleted++
			boardEntity.PointsCompleted += response.Points
			if _, ok := pointsPerDay[response.Date]; ok {
				pointsPerDay[response.Date] = response.Points + pointsPerDay[response.Date]
			} else {
				pointsPerDay[response.Date] = response.Points
			}
		}
		boardEntity.Cards++
		boardEntity.Points += response.Points
	}
	log.Printf("Cards progress: %d/%d", boardEntity.CardsCompleted, boardEntity.Cards)
	log.Printf("Total points: %f/%f", boardEntity.PointsCompleted, boardEntity.Points)
	saveToDatabase(boardEntity, pointsPerDay, targetsPerDay)
}

func getLastList(board *trello.Board) (string, error) {
	var highestPos float32
	var listID string
	lists, err := board.GetLists(trello.Defaults())
	if err != nil {
		return "", err
	}
	for _, list := range lists {
		if list.Pos > highestPos {
			listID = list.ID
		}
	}
	return listID, nil
}

func getCardDetails(card *trello.Card, listID string, res chan *cardResult) {
	points := getPoints(card)

	// Get action of a card with type "createCard"
	createdActions, err := card.GetActions(Arguments{"filter": "createCard"})
	if err != nil {
		res <- &cardResult{
			Error: err,
		}
		return
	}

	var dateCreated *time.Time

	// Get date of card creation
	for _, action := range createdActions {
		dateCreated = &action.Date
	}

	if card.IDList != listID {
		res <- &cardResult{
			Complete: false,
			Created:  dateCreated.Format("2006-01-02"),
			Points:   points,
		}
		return
	}
	actions, err := card.GetActions(trello.Defaults())
	if err != nil {
		res <- &cardResult{
			Error: err,
		}
		return
	}

	date := card.DateLastActivity
	for _, action := range actions {
		if action.Data.ListAfter != nil && action.Data.ListBefore != nil &&
			action.Data.ListAfter.ID != action.Data.ListBefore.ID && action.Data.ListAfter.ID == listID {
			date = &action.Date
			break
		}
	}

	res <- &cardResult{
		Complete: true,
		Date:     date.Format("2006-01-02"),
		Created:  dateCreated.Format("2006-01-02"),
		Points:   points,
	}
}

func getPoints(card *trello.Card) float64 {
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
