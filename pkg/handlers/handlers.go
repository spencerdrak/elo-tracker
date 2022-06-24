package handlers

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/spencerdrak/elotracker/pkg/games"
	"github.com/spencerdrak/elotracker/pkg/players"
	"github.com/spencerdrak/elotracker/pkg/util"
)

func GetAllPlayersHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	players, err := players.GetAllPlayers(db)

	if err != nil {
		util.HandleError(w, r, err)
		return
	}

	jsonPlayers, marshalErr := json.Marshal(players)

	if marshalErr != nil {
		util.HandleError(w, r, &util.EloTrackerError{
			Inner:             err,
			UserReturnMessage: "could not marshal response body",
			StatusCode:        400,
			Status:            http.StatusBadRequest,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonPlayers)
}

func InsertPlayerHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.HandleError(w, r, &util.EloTrackerError{
			Inner:             err,
			UserReturnMessage: "could not read request body",
			StatusCode:        400,
			Status:            http.StatusBadRequest,
		})
		return
	}

	var player players.Player
	if err := json.Unmarshal(body, &player); err != nil {
		util.HandleError(w, r, &util.EloTrackerError{
			Inner:             err,
			UserReturnMessage: "could not unmarshal request body",
			StatusCode:        400,
			Status:            http.StatusBadRequest,
		})
		return
	}

	if err := player.InsertInto(db); err != nil {
		util.HandleError(w, r, err)
		return
	}

	jsonPlayer, err := json.Marshal(player)

	if err != nil {
		util.HandleError(w, r, &util.EloTrackerError{
			Inner:             err,
			UserReturnMessage: "could not marshal (request) response body",
			StatusCode:        400,
			Status:            http.StatusBadRequest,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonPlayer)
}

func InsertGameHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		util.HandleError(w, r, &util.EloTrackerError{
			Inner:             err,
			UserReturnMessage: "could not read request body",
			StatusCode:        400,
			Status:            http.StatusBadRequest,
		})
		return
	}

	var game games.Game
	if err := json.Unmarshal(body, &game); err != nil {
		util.HandleError(w, r, &util.EloTrackerError{
			Inner:             err,
			UserReturnMessage: "could not unmarshal request body",
			StatusCode:        400,
			Status:            http.StatusBadRequest,
		})
		return
	}

	gamePlayers, gameAddErr := game.AddTo(db)

	if gameAddErr != nil {
		util.HandleError(w, r, gameAddErr)
		return
	}

	jsonPlayer, err := json.Marshal(gamePlayers)

	if err != nil {
		util.HandleError(w, r, &util.EloTrackerError{
			Inner:             err,
			UserReturnMessage: "could not marshal (request) response body",
			StatusCode:        400,
			Status:            http.StatusBadRequest,
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonPlayer)
}
