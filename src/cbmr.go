package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
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

type CurrentMatch struct {
	player1  string
	player2  string
	seed     string
	matchId  int32
	category string
}

var (
	connectedClients []ConnClientData
	matchInPlay      bool  // Is there a match currently going on?
	match_id         int64 = 0
	currentMatch     CurrentMatch
)

func main() {
	http.HandleFunc("/connect", connectClient)
	http.HandleFunc("/start_match", beginMatch)
	http.HandleFunc("/match_info", match_info)

	log.Println("Server running on port 3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func connectClient(w http.ResponseWriter, r *http.Request) {
	var currentClient ConnClientData

	// /connect?name=[name]
	query, _ := url.ParseQuery(r.URL.RawQuery)
	name := query.Get("name")

	currentClient.playerName = name
	currentClient.clientAddress = strings.Split(r.RemoteAddr, ":")[0]

	// TODO: check if already exists

	connectedClients = append(connectedClients, currentClient)
	log.Printf("Finished connecting client %d to server!\n", len(connectedClients))
}

func generateSeed() int32 {
	return rand.Int31n(((1 << 31) - 1))
}

var goodSeeds = []string{"badsigfile", "446456054", "33490196", "x9mc", "557110973",
	"33490196", "327675199", "990066099", "2s4n2z", "69589057"}

func beginMatch(w http.ResponseWriter, r *http.Request) {
	// start_match?player1=[player1]&category=[category]
	query, _ := url.ParseQuery(r.URL.RawQuery)

	player1 := query.Get("player1")
	category := query.Get("category")

	// Randomly select a player from the list
	if len(connectedClients) < 1 {
		io.WriteString(w, "FAILURE")
		return
	}

	player2 := connectedClients[rand.Intn(len(connectedClients))].playerName
	for player2 == player1 {
		player2 = connectedClients[rand.Intn(len(connectedClients))].playerName
	}

	log.Printf("Selected players for match: %s vs %s in %s%%.\n",
		player1, player2, category)

	/*
			{
			"id": [match_id],
			"players": [
				"player1": [player1],
				"player2": [player2]
			],
			"seed":	    [generate_seed]
			"category": [category]
		}
	*/
	if matchInPlay == true {
		return
	}
	matchInPlay = true

	var seedChosen string
	if category == "random" {
		seedChosen = strconv.Itoa(int(generateSeed()))
	} else if category == "any" {
		seedChosen = goodSeeds[rand.Intn(len(goodSeeds))]
	}

	currentMatch = CurrentMatch{player1, player2, seedChosen, int32(match_id), category}

	io.WriteString(w, fmt.Sprintf("{\n\t\"id\": %d,\n\t\"players\": [\n\t\t\"player1\": \"%s\",\n\t\t\"player2\": \"%s\"\n\t],\n\t\"seed\": \"%s\"\n\t\"category\": %s\n}\n",
		match_id, player1, player2, seedChosen, category))
	match_id++
}

func match_info(w http.ResponseWriter, r *http.Request) {
	if matchInPlay {
		io.WriteString(w, fmt.Sprintf("{\n\t\"id\": %d,\n\t\"players\": [\n\t\t\"player1\": \"%s\",\n\t\t\"player2\": \"%s\"\n\t],\n\t\"seed\": \"%s\"\n\t\"category\": %s\n}\n",
			currentMatch.matchId, currentMatch.player1, currentMatch.player2, currentMatch.seed, currentMatch.category))
	}
}
