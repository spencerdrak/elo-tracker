package players

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/spencerdrak/elotracker/pkg/util"
)

type Player struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Rating   int    `json:"rating"`
}

func (p Player) InsertInto(db *sql.DB) *util.EloTrackerError {
	insert, err := db.Prepare("INSERT INTO players (username, rating) " +
		"VALUES" +
		"($1, $2);")
	if err != nil {
		return &util.EloTrackerError{
			Inner:             err,
			UserReturnMessage: "inserting player failed",
			StatusCode:        500,
			Status:            http.StatusInternalServerError,
		}
	}
	defer insert.Close()

	_, err = insert.Exec(p.Username, p.Rating)

	if err != nil {
		return &util.EloTrackerError{
			Inner:             err,
			UserReturnMessage: "inserting player failed",
			StatusCode:        500,
			Status:            http.StatusInternalServerError,
		}
	}

	return nil
}

func GetAllPlayers(db *sql.DB) ([]Player, *util.EloTrackerError) {
	rows, err := db.Query("SELECT id, username, rating FROM players")

	if err != nil {
		return nil, &util.EloTrackerError{
			Inner:             err,
			UserReturnMessage: "querying all players failed",
			StatusCode:        500,
			Status:            http.StatusInternalServerError,
		}
	}

	var output = []Player{}
	for rows.Next() {
		var player = Player{}
		err := rows.Scan(&player.Id, &player.Username, &player.Rating)
		if err != nil {
			return nil, &util.EloTrackerError{
				Inner:             err,
				UserReturnMessage: "querying all players failed",
				StatusCode:        500,
				Status:            http.StatusInternalServerError,
			}
		}

		output = append(output, player)
	}

	return output, nil
}

func GetPlayerByUsername(db *sql.DB, username string) (Player, *util.EloTrackerError) {
	row := db.QueryRow("SELECT id, username, rating FROM players WHERE username = $1;", username)

	player := Player{}

	err := row.Scan(&player.Id, &player.Username, &player.Rating)
	if err != nil {
		if err == sql.ErrNoRows {
			return player, &util.EloTrackerError{
				Inner:             err,
				UserReturnMessage: fmt.Sprintf("player with username %s not found", username),
				StatusCode:        400,
				Status:            http.StatusInternalServerError,
			}
		} else {
			return player, &util.EloTrackerError{
				Inner:             err,
				UserReturnMessage: fmt.Sprintf("querying single player with username %s failed", username),
				StatusCode:        500,
				Status:            http.StatusInternalServerError,
			}
		}
	}
	return player, nil
}

func (p Player) UpdateRating(newRating int, db *sql.DB) *util.EloTrackerError {
	update, err := db.Prepare("UPDATE players SET rating = $1 WHERE id = $2")
	if err != nil {
		return &util.EloTrackerError{
			Inner:             err,
			UserReturnMessage: fmt.Sprintf("failed to update rating for username %s", p.Username),
			StatusCode:        500,
			Status:            http.StatusInternalServerError,
		}
	}
	defer update.Close()

	_, err = update.Exec(newRating, p.Id)

	if err != nil {
		return &util.EloTrackerError{
			Inner:             err,
			UserReturnMessage: fmt.Sprintf("failed to update rating for username %s", p.Username),
			StatusCode:        500,
			Status:            http.StatusInternalServerError,
		}
	}

	return nil
}
