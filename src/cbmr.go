package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type ConnClientData struct {
	// We need to send requests to the client like, twice.
	// We store it here. Additionally, we're going to assume
	// that the client's port is 8080.
	clientAddress string
	playerName    string
	// PlayerELO is not known to the server, only in the database,
	// so we don't need to store it at all.
	/* playerElo  int */
}

var (
	connectedClients []ConnClientData
	matchInPlay      bool // Is there a match currently going on?
)

func main() {
	http.HandleFunc("/connect", connectClient)
	http.HandleFunc("/start_match", beginMatch)

	log.Println("Server running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func connectClient(w http.ResponseWriter, r *http.Request) {
	var currentClient ConnClientData
	log.Printf("I have recieved a new client connecting to the server!")

	// /connect?name=[name]
	query, _ := url.ParseQuery(r.URL.RawQuery)
	name := query.Get("name")

	currentClient.playerName = name
	currentClient.clientAddress = strings.Split(r.RemoteAddr, ":")[0]

	connectedClients = append(connectedClients, currentClient)
	log.Printf("Finished connecting client to server! There are %d clients online\n", len(connectedClients))
	log.Println(currentClient)
}

func isConnectedPlayer(playerName string) bool {
	var r bool = false
	for i := 0; i < len(connectedClients); i++ {
		if connectedClients[i].playerName == playerName {
			r = true
		}
	}
	return r
}

func beginMatch(w http.ResponseWriter, r *http.Request) {
	query, _ := url.ParseQuery(r.URL.RawQuery)
	log.Println(query)

	player1 := query.Get("player1")
	player2 := query.Get("player2")

	// If either one or the other are not connected, do not continue,
	// make sure a match cannot start.
	if !isConnectedPlayer(player1) || !isConnectedPlayer(player2) {
		io.WriteString(w, "FAILURE")
		return
	}
}
