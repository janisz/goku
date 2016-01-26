package main

import (
	log "github.com/Sirupsen/logrus"
	"net/http"
	"fmt"
)

var VERSION string

func main() {

	log.WithField("Version", VERSION).Info("Starting marathon-consul")

	config, err := New()
	if err != nil {
		log.Fatal(err.Error())
	}

	// set up routes
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "OK")
	})

	Subscribe()

	log.WithField("Port", config.Web.Listen).Info("Listening")
	log.Fatal(http.ListenAndServe(config.Web.Listen, nil))
}
