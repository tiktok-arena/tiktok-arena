package contests

import (
	"fmt"
	"testing"
	"tiktok-arena/internal/core/models"
)

func BenchmarkKingOfTheHill(b *testing.B) {
	var tiktoks []models.Tiktok
	for i := 0; i < 64; i++ {
		tiktoks = append(tiktoks, models.Tiktok{
			TournamentID: nil,
			Tournament:   nil,
			Name:         fmt.Sprint("name", i),
			URL:          fmt.Sprint("testurl", i),
			Wins:         5432,
		})
	}

	for n := 0; n < b.N; n++ {
		result := KingOfTheHill(tiktoks)
		_ = result
	}
}
