package src

type GameHistoryEntry struct {
	Winner  string `json:"winner"`
	Loser   string `json:"loser"`
	IsDraw  bool   `json:"isDraw"`
	DateStr string `json:"dateStr"`
}

