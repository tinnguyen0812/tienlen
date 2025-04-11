package models

type TurnManager struct {
	PlayerOrder []string
	CurrentTurn int
}

func NewTurnManager(playerIDs []string, room *Room) *TurnManager {
	firstTurn := findPlayerWithThreeSpades(room)
	return &TurnManager{
		PlayerOrder: playerIDs,
		CurrentTurn: firstTurn,
	}
}

func (tm *TurnManager) GetCurrentPlayer() string {
	return tm.PlayerOrder[tm.CurrentTurn]
}

func (tm *TurnManager) NextTurn() {
	tm.CurrentTurn = (tm.CurrentTurn + 1) % len(tm.PlayerOrder)
}

func findPlayerWithThreeSpades(room *Room) int {
	for i, playerID := range getSortedPlayerIDs(room) {
		cards := room.PlayerCards[playerID]
		for _, card := range cards {
			if card.Rank == "3" && card.Suit == Spades {
				return i
			}
		}
	}
	return 0
}

func getSortedPlayerIDs(room *Room) []string {
	var ids []string
	for id := range room.Players {
		ids = append(ids, id)
	}
	return ids
}
