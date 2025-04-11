package models

import (
	"math/rand"
	"time"
)

func GenerateDeck() []Card {
	var deck []Card
	for _, suit := range Suits {
		for _, rank := range Ranks {
			deck = append(deck, Card{Suit: suit, Rank: rank})
		}
	}
	return deck
}

func Shuffle(deck []Card) []Card {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	shuffled := make([]Card, len(deck))
	perm := r.Perm(len(deck))
	for i, v := range perm {
		shuffled[v] = deck[i]
	}
	return shuffled
}
