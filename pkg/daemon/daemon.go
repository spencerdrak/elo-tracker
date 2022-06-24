package daemon

import (
	"database/sql"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/spencerdrak/elotracker/pkg/handlers"
	"github.com/spencerdrak/elotracker/pkg/util"

	log "github.com/sirupsen/logrus"
)

func Run() {
	logger := log.WithFields(log.Fields{
		"app": "liveatc-archive-downloader",
	})

	// urlExample := "postgres://username:password@localhost:5432/database_name"
	db, err := sql.Open("mysql", os.Getenv("DATABASE_URL"))
	if err != nil {
		logger.Fatal("Unable to connect to database: %v\n", err)
	}

	err = db.Ping()
	if err != nil {
		logger.Fatal("Unable to ping DB: %v\n", err)
	}

	defer db.Close()

	contextPath := ""

	r := mux.NewRouter()
	health := r.PathPrefix(contextPath + "/health").Subrouter()
	health.HandleFunc("/liveness", util.Liveness).Methods("GET")
	health.HandleFunc("/readiness", util.Readiness).Methods("GET")

	app := r.PathPrefix(contextPath).Subrouter()

	app.HandleFunc(contextPath+"/players", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("received GET /players request")
		handlers.GetAllPlayersHandler(w, r, db)
	}).Methods("GET")

	app.HandleFunc(contextPath+"/player", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("received POST /player request")
		handlers.InsertPlayerHandler(w, r, db)
	}).Methods("POST")

	app.HandleFunc(contextPath+"/game", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("received POST /game request")
		handlers.InsertGameHandler(w, r, db)
	}).Methods("POST")

	// Bind to a port and pass our router in
	logger.Fatal(http.ListenAndServe(":8000", r))
}
