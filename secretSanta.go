package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sort"
)

type Participants struct {
	Participants []Person
}

type Pairing struct {
	Sender   Person
	Receiver Person
}

type Pairings struct {
	Pairings []Pairing
}

func getOffset(numParticipants int) int {
	return (rand.Intn(numParticipants-1) + 1)
}

func getReceiverIndex(senderIndex int, numParticipants int, offset int) int {
	return (senderIndex + offset) % numParticipants
}

func getPairing(senders []Person, senderIndex int, receiverIndex int) Pairing {
	return Pairing{senders[senderIndex], senders[receiverIndex]}
}

func getPairings(senders []Person) Pairings {
	numParticipants := len(senders)
	offset := getOffset(numParticipants)

	var pairings []Pairing
	for senderIndex := range senders {
		pairings = append(pairings, getPairing(senders, senderIndex,
			getReceiverIndex(senderIndex, numParticipants, offset)))
	}

	return Pairings{pairings}
}

func getSenders(participants Participants) []Person {
	senders := participants.Participants
	sort.Slice(senders, func(i, j int) bool {
		return rand.Int() < rand.Int()
	})

	return senders
}

func pair(w http.ResponseWriter, r *http.Request) {
	// Declare a new Person struct.
	var participants Participants

	err := json.NewDecoder(r.Body).Decode(&participants)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	pairings, _ := json.Marshal(getPairings(getSenders(participants)))
	fmt.Fprintf(w, string(pairings))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/pair", pair)

	error := http.ListenAndServe(":4000", mux)
	log.Fatal(error)
}
