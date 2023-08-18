package contests

import (
	"fmt"
	"github.com/google/uuid"
	"testing"
	"tiktok-arena/internal/core/models"
)

func BenchmarkSingleElimination(b *testing.B) {
	var tiktoks []models.Tiktok
	for i := 0; i < 64; i++ {
		tiktoks = append(tiktoks, models.Tiktok{
			TournamentID: uuid.UUID{},
			Tournament:   models.Tournament{},
			Name:         fmt.Sprint("name", i),
			URL:          fmt.Sprint("testurl", i),
			Wins:         5432,
		})
	}

	for n := 0; n < b.N; n++ {
		result := SingleElimination(tiktoks)
		_ = result
	}
}
