package dtos

const (
	SingleElimination = "single_elimination"
	KingOfTheHill     = "king_of_the_hill"
)

func GetAllowedBracketType() map[string]bool {
	return map[string]bool{
		SingleElimination: true,
		KingOfTheHill:     true,
	}
}

func CheckIfAllowedBracketType(bracketType string) bool {
	return GetAllowedBracketType()[bracketType]
}

type Bracket struct {
	CountMatches int
	Rounds       []Round
}

type Round struct {
	Round   int
	Matches []Match
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

type Match struct {
	MatchID      string
	FirstOption  Option
	SecondOption Option
}

type BracketPayload struct {
	BracketType string `validate:"required"`
}
