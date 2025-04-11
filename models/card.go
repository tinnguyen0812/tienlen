package models

import (
	"math/rand"
	"time"
)

type Card struct {
	Rank string `json:"rank"`
	Suit string `json:"suit"`
}

func GenerateDeck() []Card {
	ranks := []string{"3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A", "2"}
	suits := []string{"♠", "♣", "♦", "♥"}

	var deck []Card
	for _, r := range ranks {
		for _, s := range suits {
			deck = append(deck, Card{Rank: r, Suit: s})
		}
	}
	return deck
}

func Shuffle(deck []Card) []Card {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(deck), func(i, j int) {
		deck[i], deck[j] = deck[j], deck[i]
	})
	return deck
}
