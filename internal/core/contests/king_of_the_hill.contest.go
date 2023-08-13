package contests

import (
	"github.com/google/uuid"
	"tiktok-arena/internal/core/dtos"
	"tiktok-arena/internal/core/models"
)

// KingOfTheHill
// First match decided randomly between two participators.
// Loser of match leaves the game, winner will go to next match, next opponent decided randomly from standings.
// Procedure continues until last standing.
func KingOfTheHill(t []models.Tiktok) dtos.Contest {
	countTiktok := len(t)
	rounds := make([]dtos.Round, 0, countTiktok-1)
	match := dtos.Match{
		MatchID:      uuid.NewString(),
		FirstOption:  dtos.TiktokOption{TiktokURL: t[0].URL},
		SecondOption: dtos.TiktokOption{TiktokURL: t[1].URL},
	}
	rounds = append(rounds, dtos.Round{
		Round:   1,
		Matches: []dtos.Match{match},
	})
	previousMatch := match
	for i := 2; i < countTiktok-1; i++ {
		match = dtos.Match{
			MatchID:      uuid.NewString(),
			FirstOption:  dtos.MatchOption{MatchID: previousMatch.MatchID},
			SecondOption: dtos.TiktokOption{TiktokURL: t[i].URL},
		}
		rounds = append(rounds, dtos.Round{
			Round:   i,
			Matches: []dtos.Match{match},
		})
		previousMatch = match
	}
	return dtos.Contest{
		CountMatches: countTiktok - 1,
		Rounds:       rounds,
	}
}
