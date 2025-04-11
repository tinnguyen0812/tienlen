package rules

import (
	"sort"
	"strings"
)

var ValueOrder = map[string]int{
	"3": 0, "4": 1, "5": 2, "6": 3, "7": 4, "8": 5,
	"9": 6, "10": 7, "J": 8, "Q": 9, "K": 10, "A": 11, "2": 12,
}

func parseCards(cards []string) []string {
	values := make([]string, 0, len(cards))
	for _, c := range cards {
		value := strings.TrimRight(c, "♠♥♦♣") // remove suit
		values = append(values, value)
	}
	return values
}

func DetectCombination(cards []string) (string, int) {
	values := parseCards(cards)
	countMap := map[string]int{}
	for _, v := range values {
		countMap[v]++
	}

	switch len(cards) {
	case 1:
		return "single", ValueOrder[values[0]]
	case 2:
		if values[0] == values[1] {
			return "pair", ValueOrder[values[0]]
		}
	case 3:
		if values[0] == values[1] && values[1] == values[2] {
			return "triple", ValueOrder[values[0]]
		}
	case 4:
		if values[0] == values[1] && values[1] == values[2] && values[2] == values[3] {
			return "four_of_a_kind", ValueOrder[values[0]]
		}
	default:
		if isStraight(values) {
			return "straight", ValueOrder[values[len(values)-1]]
		}
	}

	return "invalid", -1
}

func isStraight(values []string) bool {
	if len(values) < 3 || len(values) > 12 {
		return false
	}

	vals := []int{}
	seen := map[int]bool{}
	for _, v := range values {
		order := ValueOrder[v]
		if seen[order] {
			return false
		}
		seen[order] = true
		vals = append(vals, order)
	}
	sort.Ints(vals)

	for i := 1; i < len(vals); i++ {
		if vals[i]-vals[i-1] != 1 {
			return false
		}
	}
	return true
}
