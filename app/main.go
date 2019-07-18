package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/sotah-inc/steamwheedle-cartel/pkg/logging"
)

var port int

func init() {
	parsedPort, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Failed to get port: %s", err.Error())

		return
	}

	port = parsedPort
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, "Hello, world!!!"); err != nil {
			logging.WithField("error", err.Error()).Error("")

			return
		}
	}).Methods("POST")
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
