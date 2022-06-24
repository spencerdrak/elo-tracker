package util

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type EloTrackerError struct {
	Inner             error
	UserReturnMessage string
	StatusCode        int
	Status            int
}

func (e *EloTrackerError) Error() string {
	return fmt.Sprintf("%v; Returning to user: %v;", e.Inner, e.UserReturnMessage)
}

func (e *EloTrackerError) Unwrap() error {
	return e.Inner
}

func HandleError(w http.ResponseWriter, r *http.Request, err *EloTrackerError) {
	log.Error("ERROR: " + err.Error())
	http.Error(w, fmt.Sprintf("%s: %s", http.StatusText(err.StatusCode), err.UserReturnMessage), err.Status)
}

// Liveness checks is webapp is ready to serve traffic
func Liveness(w http.ResponseWriter, r *http.Request) {
	log.Info("Liveness: 200")
	w.WriteHeader(200)
}

// Readiness checks if DB connection is good
func Readiness(w http.ResponseWriter, r *http.Request) {
	log.Info("Readiness: 200")
	w.WriteHeader(200)
}
