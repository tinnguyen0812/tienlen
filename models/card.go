package models

type Suit string
type Rank string

const (
	Spades   Suit = "♠"
	Hearts   Suit = "♥"
	Diamonds Suit = "♦"
	Clubs    Suit = "♣"
)

var Suits = []Suit{Spades, Hearts, Diamonds, Clubs}
var Ranks = []Rank{
	"3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K", "A", "2",
}

type Card struct {
	Suit Suit `json:"suit"`
	Rank Rank `json:"rank"`
}
