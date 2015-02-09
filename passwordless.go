package main

import (
	"flag"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	passwordless "github.com/iamjem/go-passwordless-demo/lib"
	"github.com/stretchr/graceful"
	"net/http"
	"os"
	"path"
	"time"
)

var port = flag.String("port", "8080", "listen address")

func main() {
	flag.Parse()

	listen := fmt.Sprintf(":%s", *port)

	router := mux.NewRouter()

	// static files
	root, _ := os.Getwd()
	router.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir(path.Join(root, "public")))))

	// web routes
	router.PathPrefix("/").Handler(passwordless.BuildRoutes())

	// setup server
	var handler http.Handler

	// if Debug is true, enable logging
	if os.Getenv("DEBUG") == "true" {
		log.SetLevel(log.DebugLevel)
		handler = handlers.CombinedLoggingHandler(os.Stdout, router)
	} else {
		handler = router
	}

	log.WithFields(log.Fields{
		"listen": listen,
	}).Info("Server running")

	graceful.Run(listen, 10*time.Second, handler)
}
