package services

import (
	"testing"
	"tiktok-arena/internal/core/dtos"
	"tiktok-arena/internal/data/repository"
)

func TestTournamentService_CreateTournament(t *testing.T) {
	test := SetupIntegration(t)
	test.Cleanup()

	userRepository := repository.NewUserRepository(test.db)
	tiktokRepository := repository.NewTiktokRepository(test.db)
	tournamentRepository := repository.NewTournamentRepository(test.db)

	authService := NewAuthService(userRepository)
	tournamentService := NewTournamentService(tournamentRepository, tiktokRepository)
	t.Run("Create new user for testing. "+
		"Create tournament with tiktoks in it."+
		"Get tournaments and tiktoks."+
		"List tournaments and tiktoks.", func(t *testing.T) {
		// Create new user for testing
		userDetails := dtos.AuthInput{Name: "test", Password: "test"}
		details, err := authService.NewUser(&userDetails)
		if err != nil {
			t.Fatalf("Error in creating user. Expected %v, got %v.", nil, err.Error())
		}

		// Create createTournament with tiktoks in it
		createTournament := dtos.CreateTournament{
			Name:     "test_tournament_name_1",
			PhotoURL: "http://test_url.ua",
			Size:     5,
			Tiktoks: []dtos.CreateTiktok{
				{Name: "test_tiktok_name_1", URL: "http://tiktok.com/@username/1"},
				{Name: "test_tiktok_name_2", URL: "http://tiktok.com/@username/2"},
				{Name: "test_tiktok_name_3", URL: "http://tiktok.com/@username/3"},
				{Name: "test_tiktok_name_4", URL: "http://tiktok.com/@username/4"},
				{Name: "test_tiktok_name_5", URL: "http://tiktok.com/@username/5"},
			},
		}
		err = tournamentService.CreateTournament(&createTournament, *details.ID)
		if err != nil {
			t.Fatalf("Error in creating createTournament. Expected %v, got %v.", nil, err.Error())
		}
	})

}
