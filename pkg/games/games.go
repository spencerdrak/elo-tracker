package games

import (
	"database/sql"
	"math"
	"net/http"

	"github.com/spencerdrak/elotracker/pkg/players"
	"github.com/spencerdrak/elotracker/pkg/util"
)

type Game struct {
	Id             int64  `json:"id"`
	WinnerUsername string `json:"winnerUsername"`
	LoserUsername  string `json:"loserUsername"`
	WinMethod      string `json:"winMethod"`
}

func (g Game) AddTo(db *sql.DB) ([]players.Player, *util.EloTrackerError) {
	winnerUser, err := players.GetPlayerByUsername(db, g.WinnerUsername)
	if err != nil {
		return nil, err
	}

	loserUser, err := players.GetPlayerByUsername(db, g.LoserUsername)
	if err != nil {
		return nil, err
	}

	rCurrentWinner := math.Pow(10, float64(winnerUser.Rating)/400)
	rCurrentLoser := math.Pow(10, float64(loserUser.Rating)/400)

	expectedScoreWinner := rCurrentWinner / (rCurrentWinner + rCurrentLoser)
	expectedScoreLoser := rCurrentLoser / (rCurrentWinner + rCurrentLoser)

	// S(1) = 1 if player 1 wins / 0.5 if draw / 0 if player 2 wins
	// S(2) = 0 if player 1 wins / 0.5 if draw / 1 if player 2 wins

	sWinner := 1.0
	sLoser := 0.0

	if g.WinMethod == "draw" {
		sWinner = 0.5
		sLoser = 0.5
	}

	// r'(1) = r(1) + K * (S(1) – E(1))
	winnerNewRating := winnerUser.Rating + int(math.Round(32*(sWinner-expectedScoreWinner)))

	//r'(2) = r(2) + K * (S(2) – E(2))
	loserNewRating := loserUser.Rating + int(math.Round((32 * (sLoser - expectedScoreLoser))))

	insert, prepareErr := db.Prepare("INSERT INTO games (winner_username, loser_username, win_method) " +
		"VALUES" +
		"(?, ?, ?)")
	if prepareErr != nil {
		return nil, &util.EloTrackerError{
			Inner:             err,
			UserReturnMessage: "failed to insert game",
			StatusCode:        500,
			Status:            http.StatusInternalServerError,
		}
	}
	defer insert.Close()

	_, insertErr := insert.Exec(g.WinnerUsername, g.LoserUsername, g.WinMethod)

	if insertErr != nil {
		return nil, &util.EloTrackerError{
			Inner:             err,
			UserReturnMessage: "failed to insert game",
			StatusCode:        500,
			Status:            http.StatusInternalServerError,
		}
	}

	err = winnerUser.UpdateRating(winnerNewRating, db)

	if err != nil {
		return nil, err
	}

	err = loserUser.UpdateRating(loserNewRating, db)

	if err != nil {
		return nil, err
	}

	winnerUser.Rating = winnerNewRating
	loserUser.Rating = loserNewRating

	players := []players.Player{winnerUser, loserUser}

	return players, nil
}
