package dtos

const (
	SingleElimination = "single_elimination"
	KingOfTheHill     = "king_of_the_hill"
)

func GetAllowedContestType() map[string]bool {
	return map[string]bool{
		SingleElimination: true,
		KingOfTheHill:     true,
	}
}

func CheckIfAllowedContestType(contestType string) bool {
	return GetAllowedContestType()[contestType]
}

type Contest struct {
	CountMatches int     `json:"countMatches"`
	Rounds       []Round `json:"rounds"`
}

type Round struct {
	Round   int     `json:"round"`
	Matches []Match `json:"matches"`
}

type Match struct {
	MatchID      string `json:"matchID"`
	FirstOption  Option `json:"firstOption"`
	SecondOption Option `json:"secondOption"`
}

type Option interface {
	isOption() bool
}

type MatchOption struct {
	MatchID string `json:"matchID"`
}

func (m MatchOption) isOption() bool {
	return true
}

type TiktokOption struct {
	TiktokURL string `json:"tiktokURL"`
}

func (m TiktokOption) isOption() bool {
	return true
}

type ContestPayload struct {
	Type string `validate:"required" json:"type"`
}
