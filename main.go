package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/gplus"
	"github.com/sirupsen/logrus"
)

// Global Variables
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

	// Routing
	router := mux.NewRouter()
	router.HandleFunc("/auth/{provider}/callback", loginCallbackHandler)
	router.HandleFunc("/auth/{provider}", loginHandler)
	router.HandleFunc("/logout", logoutHandler)
	router.PathPrefix("/chat").Handler(MustAuth(&templateHandler{filename: "chat.html"})).Methods("GET")
	router.PathPrefix("/login").Handler(&templateHandler{filename: "login.html"}).Methods("GET")
	room := newRoom()
	router.PathPrefix("/room").Handler(room).Methods("GET")
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	log.Debug("main: End routing")

	// Run the chatroom.
	go room.run()

	// Start the web server.
	log.Info("Start the web server - Listen port address ", *port)
	log.Fatal(http.ListenAndServe(*port, router))
}
