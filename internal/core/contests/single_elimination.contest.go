package contests

import (
	"github.com/google/uuid"
	"math"
	"tiktok-arena/internal/core/dtos"
	"tiktok-arena/internal/core/models"
)

// SingleElimination
// https://en.wikipedia.org/wiki/Single-elimination_tournament
func SingleElimination(t []models.Tiktok) *dtos.Contest {
	countTiktok := len(t)
	countRound := int(math.Ceil(math.Log2(float64(countTiktok))))
	countSecondRoundParticipators := 1 << (countRound - 1) // Equivalent to int(math.Pow(2, float64(countRound)) / 2)
	countFirstRoundMatches := countTiktok - int(math.Pow(2, float64(countRound)-1))
	countFirstRoundTiktoks := countFirstRoundMatches * 2

	rounds := make([]dtos.Round, 0, countRound)

	firstRoundMatches := make([]dtos.Match, 0, countFirstRoundMatches)
	secondRoundMatches := make([]dtos.Match, 0, countSecondRoundParticipators/2)

	secondRoundParticipators := make([]dtos.Option, 0, countSecondRoundParticipators) // This slice should store MatchOption or TiktokOption

	// Filling first round with firstRoundMatches and appending MatchOptions to second round participators
	for j := 0; j < countFirstRoundTiktoks; j += 2 {
		matchID := uuid.NewString()
		firstRoundMatches = append(firstRoundMatches, dtos.Match{
			MatchID: matchID,
			FirstOption: dtos.TiktokOption{
				TiktokURL: t[j].URL,
			},
			SecondOption: dtos.TiktokOption{
				TiktokURL: t[j+1].URL,
			},
		})
		secondRoundParticipators = append(secondRoundParticipators,
			dtos.MatchOption{MatchID: matchID})
	}
	// Appending first round firstRoundMatches to rounds
	rounds = append(rounds, dtos.Round{
		Round:   1,
		Matches: firstRoundMatches,
	})
	// Appending TiktokOptions to second round participators
	for _, tiktok := range t[countFirstRoundTiktoks:] {
		secondRoundParticipators = append(secondRoundParticipators,
			dtos.TiktokOption{TiktokURL: tiktok.URL})
	}
	// Generating second round firstRoundMatches
	for i := 0; i < int(countSecondRoundParticipators); i += 2 {
		match := dtos.Match{
			MatchID:      uuid.NewString(),
			FirstOption:  secondRoundParticipators[i],
			SecondOption: secondRoundParticipators[i+1],
		}
		secondRoundMatches = append(secondRoundMatches, match)
	}
	// Generating second round
	secondRound := dtos.Round{
		Round:   2,
		Matches: secondRoundMatches}
	rounds = append(rounds, secondRound)

	previousRoundMatches := secondRoundMatches
	for roundID := 3; roundID <= countRound; roundID++ {
		// Generating Nth round matches (where N > 2)
		var currentRoundMatches []dtos.Match
		for matchID := 0; matchID < len(previousRoundMatches); matchID += 2 {
			match := dtos.Match{
				MatchID: uuid.NewString(),
				FirstOption: dtos.MatchOption{
					MatchID: previousRoundMatches[matchID].MatchID,
				},
				SecondOption: dtos.MatchOption{
					MatchID: previousRoundMatches[matchID+1].MatchID,
				},
			}
			currentRoundMatches = append(currentRoundMatches, match)
		}
		// Generating Nth round (where N > 2)
		round := dtos.Round{
			Round:   roundID,
			Matches: currentRoundMatches,
		}
		rounds = append(rounds, round)

		previousRoundMatches = currentRoundMatches
	}
	return &dtos.Contest{
		CountMatches: countTiktok - 1,
		Rounds:       rounds,
	}

}
