package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func handle_err(err error) {
	if err != nil {
		log.Println(err)
		return
	}
}

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

const (
	MATCH_IN_PLAY = iota
	MATCH_FOREFEIT
	MATCH_DRAW
	PLAYER1_WIN
	PLAYER2_WIN
)

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
	http.HandleFunc("/end_match", end_match)

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

func getIndexOfPlayer(s string) int {
	for i := 0; i < len(connectedClients); i++ {
		if connectedClients[i].playerName == s {
			return i
		}
	}
	return -1
}

func match_info(w http.ResponseWriter, r *http.Request) {
	// If the client GETs this route, then we just want to display a JSON object that displays what
	// the match looks like.
	if r.Method == http.MethodGet {
		if matchInPlay {
			io.WriteString(w, fmt.Sprintf("{\n\t\"id\": %d,\n\t\"players\": [\n\t\t\"player1\": \"%s\",\n\t\t\"player2\": \"%s\"\n\t],\n\t\"seed\": \"%s\"\n\t\"category\": %s\n}\n",
				currentMatch.matchId, currentMatch.player1, currentMatch.player2, currentMatch.seed, currentMatch.category))
		}
		return
	}

	// If POST, then we've got something about the timing/match ending!
	// We'll assume it's a JSON object of the form
	/*
		{
			"match_result": "win" 							// Based off the player who sent it.
			"match_time": 	"0:12.345",
			"sender":			  "[player1's name]"
		}
	*/
	if r.Method == http.MethodPost {
		// Decode JSON message
		var matchIf matchInfoPOST
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			return
		}
		defer r.Body.Close()
		err = json.Unmarshal(body, &matchIf)
		if err != nil {
			log.Fatalln(err)
		}

		// We need to send a message to player2, so let's get the player we need to send.
		messageToSendTo := currentMatch.player2
		if matchIf.sender == currentMatch.player2 {
			messageToSendTo = currentMatch.player1
		}

		// Next we'll send the packet to the client, and the client can handle it
		log.Printf("wtf idk %s\n", connectedClients[getIndexOfPlayer(messageToSendTo)])
		bruh := fmt.Sprintf("http://%s:8080/match_info?data=%s",
			connectedClients[getIndexOfPlayer(messageToSendTo)].clientAddress,
			url.QueryEscape(string(body)))
		fmt.Println(bruh)
		http.Get(bruh)
	}
}

type matchInfoPOST struct {
	match_result string
	match_time   string
	forefeit_opp string
	sender       string
}

// Once player 2 acknowledges (if win/loss) / accepts (draw/forfeit),
// *player2* will call this. Additionally, player1 should continously
// poll this route with a GET request every few seconds to check if the
// server is running

// player2's client should send a POST request with JSON,
// that should look like this:
/*
	{
		"winner": [player1/player2],
		"time": [time],
		"accepted": true/false
	}
*/
func end_match(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		io.WriteString(w, fmt.Sprintf("%t", matchInPlay))
		return
	}
	// method should be post, that's the only one that the client should expose.
	// now let's parse the JSON output object.

	var data map[string]interface{}
	body, err := io.ReadAll(r.Body)
	handle_err(err)

	defer r.Body.Close()
	err = json.Unmarshal(body, &data)
	handle_err(err)

	// We now have the JSON,  modify ELO and other stuff.
	log.Printf("match data posted!\n")
	log.Println(data)

	if data["accepted"] == true {
		matchInPlay = false
		log.Printf("gg! winner is %s in %s\n", data["winner"], data["time"])
	}
	if data["accepted"] == false {
		log.Printf("match is still on!")
	}
}
