package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/pat"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/gplus"
	"github.com/sirupsen/logrus"
)

// Global variables
var (
	log  = logrus.New()
	port = flag.String("port", ":8080", "Application port address.")
)

func init() {
	// Parse flags.
	flag.Parse()

	// Setup goth.
	goth.UseProviders(
		gplus.New(os.Getenv("GPLUS_KEY"), os.Getenv("GPLUS_SECRET"), fmt.Sprintf("http://localhost%s/auth/gplus/callback", *port)),
	)

	// Set the logging level on a logger.
	log.SetLevel(logrus.DebugLevel)
}

func main() {
	log.Debug("main: Start routing")

	router := pat.New()
	router.Get("/auth/{provider}/callback", loginCallbackHandler)
	router.Get("/auth/{provider}", loginHandler)
	router.Get("/logout", logoutHandler)

	router.Add("GET", "/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	router.Add("GET", "/login", &templateHandler{filename: "login.html"})

	r := newRoom()
	router.Add("GET", "/room", r)

	log.Debug("main: End routing")

	// Run the chatroom.
	go r.run()

	// Start the web server.
	log.Info("Start the web server - Listen port address ", *port)
	log.Fatal(http.ListenAndServe(*port, router))
}
