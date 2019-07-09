package steamwheedle_cartel_download_auctions_server_app

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, world!")
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
