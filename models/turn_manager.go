package models

type TurnManager struct {
	PlayerOrder []string
	CurrentIdx  int
}

func NewTurnManager(playerIDs []string) *TurnManager {
	return &TurnManager{
		PlayerOrder: playerIDs,
		CurrentIdx:  0,
	}
}

func (tm *TurnManager) GetCurrentPlayer() string {
	return tm.PlayerOrder[tm.CurrentIdx]
}

func (tm *TurnManager) NextTurn() {
	tm.CurrentIdx = (tm.CurrentIdx + 1) % len(tm.PlayerOrder)
}
