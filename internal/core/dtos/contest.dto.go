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
	CountMatches int
	Rounds       []Round
}

type Round struct {
	Round   int
	Matches []Match
}

type Match struct {
	MatchID      string
	FirstOption  Option
	SecondOption Option
}

type Option interface {
	isOption() bool
}

type MatchOption struct {
	MatchID string
}

func (m MatchOption) isOption() bool {
	return true
}

type TiktokOption struct {
	TiktokURL string
}

func (m TiktokOption) isOption() bool {
	return true
}

type ContestPayload struct {
	Type string `validate:"required"`
}
