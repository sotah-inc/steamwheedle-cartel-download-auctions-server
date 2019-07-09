package steamwheedle_cartel_download_auctions_server_app

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
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
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, world!")
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
