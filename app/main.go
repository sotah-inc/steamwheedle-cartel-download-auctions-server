package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/compute/metadata"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/sotah-inc/steamwheedle-cartel/pkg/logging"
	"github.com/sotah-inc/steamwheedle-cartel/pkg/logging/stackdriver"
	"github.com/sotah-inc/steamwheedle-cartel/pkg/sotah"
)

var port int
var serviceName string
var projectId string

func init() {
	parsedPort, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Failed to get port: %s", err.Error())

		return
	}

	port = parsedPort
	serviceName = os.Getenv("K_SERVICE")
	projectId, err = metadata.Get("project/project-id")
	if err != nil {
		log.Fatalf("Failed to get port: %s", err.Error())

		return
	}
}

func WriteErroneousResponse(w http.ResponseWriter, msg string, err error) {
	w.WriteHeader(http.StatusInternalServerError)

	if _, err := fmt.Fprint(w, msg); err != nil {
		logging.WithField("error", err.Error()).Error("Failed to write response")

		return
	}

	logging.WithField("error", err.Error()).Error(msg)
}

func main() {
	// establishing log verbosity
	logVerbosity, err := logrus.ParseLevel("info")
	if err != nil {
		logging.WithField("error", err.Error()).Fatal("Could not parse log level")

		return
	}
	logging.SetLevel(logVerbosity)

	// adding stackdriver hook
	logging.WithField("project-id", projectId).Info("Creating stackdriver hook")
	stackdriverHook, err := stackdriver.NewHook(projectId, serviceName)
	if err != nil {
		logging.WithFields(logrus.Fields{
			"error":     err.Error(),
			"projectID": projectId,
		}).Fatal("Could not create new stackdriver logrus hook")

		return
	}
	logging.AddHook(stackdriverHook)

	logging.WithField("service", serviceName).Info("Initializing service")

	r := mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logging.Info("Received request")

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			WriteErroneousResponse(w, "Could not read request body", err)

			return
		}

		regionRealm, err := sotah.NewRegionRealmTuple(string(body))
		if err != nil {
			WriteErroneousResponse(w, "Could not decode request body", err)

			return
		}

		logging.WithField("request-body", string(body)).Info("Received region-realm")

		<-time.After(time.Duration(rand.Intn(5000)) * time.Millisecond)

		logging.Info("Sending response")

		responseBody, err := regionRealm.EncodeForDelivery()
		if err != nil {
			WriteErroneousResponse(w, "Encode response body", err)

			return
		}

		if _, err := fmt.Fprint(w, responseBody); err != nil {
			logging.WithField("error", err.Error()).Error("Failed to return response")

			return
		}
	}).Methods("POST")
	http.Handle("/", r)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
