package utils

import (
	"strings"
	"tienlen-server/rules"
	"unicode"
)

func IsValidPlay(play []string, lastPlay []string) bool {
	playType, playValue := rules.DetectCombination(play)
	lastType, lastValue := rules.DetectCombination(lastPlay)

	// Không hợp lệ nếu bài hiện tại không phải tổ hợp đúng
	if playType == "invalid" {
		return false
	}

	// Nếu không có bài trước (vòng mới)
	if lastType == "invalid" || len(lastPlay) == 0 {
		return true
	}

	// Nếu khác loại bài → không được đánh (trừ một vài trường hợp đặc biệt)
	if playType != lastType {
		// Cho phép chặt heo bằng tứ quý
		if lastType == "single" && lastValue == rules.ValueOrder["2"] && playType == "four_of_a_kind" {
			return true
		}
		// Có thể mở rộng thêm logic chặt đôi heo bằng đôi thông ở đây
		return false
	}

	// Cùng loại → so sánh giá trị
	return playValue > lastValue
}

var cardOrder = map[string]int{
	"3": 0, "4": 1, "5": 2, "6": 3, "7": 4, "8": 5, "9": 6,
	"10": 7, "J": 8, "Q": 9, "K": 10, "A": 11, "2": 12,
}

func GetCardValue(card string) string {
	// card = "3♠" → return "3"
	return strings.TrimRightFunc(card, func(r rune) bool {
		return !unicode.IsNumber(r) && !unicode.IsLetter(r)
	})
}

func GetCardRank(card string) int {
	val := GetCardValue(card)
	return cardOrder[val]
}
