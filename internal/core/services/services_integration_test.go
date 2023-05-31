package services

import (
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
	"reflect"
	"testing"
	"tiktok-arena/configuration"
	"tiktok-arena/internal/core/dtos"
	"tiktok-arena/internal/data/repository"
)

// Test integrity of all services with database test instance
func TestTournamentUserAuthServices(t *testing.T) {
	test := SetupIntegration(t)
	defer test.Cleanup()

	userRepository := repository.NewUserRepository(test.db)
	tiktokRepository := repository.NewTiktokRepository(test.db)
	tournamentRepository := repository.NewTournamentRepository(test.db)

	authService := NewAuthService(userRepository)
	tournamentService := NewTournamentService(tournamentRepository, tiktokRepository)
	userService := NewUserService(userRepository, tournamentRepository)
	t.Run("Create new user for testing."+
		"Log in."+
		"Check whoami"+
		"Change user photo"+
		"Create tournament with tiktoks in it."+
		"Get all tournaments and for specific user and compare them"+
		"Edit tournament. Compare edited tournament to tournament from GetAllTournaments"+
		"Delete tournament. Check if tournament is deleted using GetAllTournaments"+
		"Create several tournaments and delete all of them with DeleteTournaments", func(t *testing.T) {
		// Create new user for testing
		userDetails := dtos.AuthInput{Name: "test", Password: "test"}
		newUser, err := authService.NewUser(&userDetails)
		if err != nil {
			t.Fatalf("Error in creating user. Expected %v, got %v.", nil, err.Error())
		}
		// Log in
		user, err := authService.GetUserByNameAndPassword(&userDetails)
		if newUser.Token != user.Token {
			t.Fatal("Token mismatch")
		}
		// Check whoami
		jwtTokenString := newUser.Token
		token, err := jwt.Parse(jwtTokenString, func(token *jwt.Token) (interface{}, error) {
			key := []byte(configuration.EnvConfig.JwtSecret)
			return key, nil
		})
		if err != nil {
			t.Fatal("jwt Parse error")
		}
		if !token.Valid {
			t.Fatal("JWT token is invalid")
		}

		whoami, err := authService.WhoAmI(token)
		if user.Token != whoami.Token {
			t.Fatal("NewUser and Whoami tokens mismatch")
		}

		// Change user photo
		userPhoto := dtos.ChangePhotoURL{PhotoURL: "http://test.url"}
		err = userService.ChangeUserPhoto(&userPhoto, *newUser.ID)
		if err != nil {
			t.Fatalf("Error in editing user photo. Expected %v, got %v.", nil, err.Error())
		}

		// Create tournament with tiktoks in it
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
		err = tournamentService.CreateTournament(&createTournament, *newUser.ID)
		if err != nil {
			t.Fatalf("Error in creating createTournament. Expected %v, got %v.", nil, err.Error())
		}

		// Get all tournaments and for specific user and compare them

		// Create simple pagination query
		paginationQuery := new(dtos.PaginationQueries)
		dtos.ValidatePaginationQueries(paginationQuery)

		tournamentsForUser, err := userService.TournamentsOfUser(*newUser.ID, paginationQuery)

		tournaments, err := tournamentService.GetAllTournaments(paginationQuery)
		if !reflect.DeepEqual(tournamentsForUser, tournaments) {
			t.Fatal("User tournaments and global tournaments mismatch")
		}
		// TournamentId for next test
		testTournament := tournaments.Tournaments[0]
		testTournamentID := testTournament.ID.String()

		// Check GetTournamentTiktoks
		tiktoks, err := tournamentService.GetTournamentTiktoks(testTournamentID)
		if len(*tiktoks) != len(createTournament.Tiktoks) {
			t.Fatal("Tiktoks mismatch")
		}

		// Edit tournament
		editTournament := dtos.EditTournament{
			Name:     "test_tournament_name_4",
			PhotoURL: "http://test_url2.ua",
			Size:     4,
			Tiktoks: []dtos.CreateTiktok{
				{Name: "test_tiktok_name_1", URL: "http://tiktok.com/@username/1"},
				{Name: "test_tiktok_name_6", URL: "http://tiktok.com/@username/6"},
				{Name: "test_tiktok_name_3", URL: "http://tiktok.com/@username/3"},
				{Name: "test_tiktok_name_2", URL: "http://tiktok.com/@username/2"},
			},
		}
		err = tournamentService.EditTournament(&editTournament, *newUser.ID, testTournamentID)
		if err != nil {
			t.Fatalf("Error in editing tournament. Expected %v, got %v.", nil, err.Error())
		}
		// Check tournament edit
		tournament, err := tournamentService.GetTournament(testTournamentID)
		if tournament.Name != editTournament.Name || tournament.Size != editTournament.Size {
			t.Fatal("Edited tournament mismatch")
		}

		// Delete tournament
		err = tournamentService.DeleteTournament(*newUser.ID, testTournamentID)
		if err != nil {
			t.Fatalf("Error in deleting tournament. Expected %v, got %v.", nil, err.Error())
		}
		// Check deletion
		tournament, err = tournamentService.GetTournament(testTournamentID)
		expectedError := RepositoryError{gorm.ErrRecordNotFound}
		if err != expectedError {
			t.Fatalf("Error in deleting tournament. Expected %v, got %v.", nil, err.Error())
		}
		// Create several tournaments and delete all of them with DeleteTournaments
		createTournaments := []dtos.CreateTournament{
			{
				Name:     "test_tournament_name_1",
				PhotoURL: "http://test_url.ua1",
				Size:     4,
				Tiktoks: []dtos.CreateTiktok{
					{Name: "test_tiktok_name_1", URL: "http://tiktok.com/@username/1"},
					{Name: "test_tiktok_name_2", URL: "http://tiktok.com/@username/2"},
					{Name: "test_tiktok_name_3", URL: "http://tiktok.com/@username/3"},
					{Name: "test_tiktok_name_4", URL: "http://tiktok.com/@username/4"},
				},
			},
			{
				Name:     "test_tournament_name_2",
				PhotoURL: "http://test_url.ua2",
				Size:     6,
				Tiktoks: []dtos.CreateTiktok{
					{Name: "test_tiktok_name_5", URL: "http://tiktok.com/@username/5"},
					{Name: "test_tiktok_name_6", URL: "http://tiktok.com/@username/6"},
					{Name: "test_tiktok_name_7", URL: "http://tiktok.com/@username/7"},
					{Name: "test_tiktok_name_8", URL: "http://tiktok.com/@username/8"},
					{Name: "test_tiktok_name_9", URL: "http://tiktok.com/@username/9"},
					{Name: "test_tiktok_name_10", URL: "http://tiktok.com/@username/10"},
				},
			},
			{
				Name:     "test_tournament_name_3",
				PhotoURL: "http://test_url.ua3",
				Size:     5,
				Tiktoks: []dtos.CreateTiktok{
					{Name: "test_tiktok_name_11", URL: "http://tiktok.com/@username/11"},
					{Name: "test_tiktok_name_12", URL: "http://tiktok.com/@username/12"},
					{Name: "test_tiktok_name_13", URL: "http://tiktok.com/@username/13"},
					{Name: "test_tiktok_name_14", URL: "http://tiktok.com/@username/14"},
					{Name: "test_tiktok_name_15", URL: "http://tiktok.com/@username/15"},
				},
			},
		}
		for _, createTournament := range createTournaments {
			err = tournamentService.CreateTournament(&createTournament, *newUser.ID)
			if err != nil {
				t.Fatalf("Error in creating createTournament. Expected %v, got %v.", nil, err.Error())
			}
		}
		tournamentResponse, err := tournamentService.GetAllTournaments(paginationQuery)
		if err != nil {
			t.Fatalf("Error in GetAllTournaments. Expected %v, got %v.", nil, err.Error())
		}
		var tournamentIds []string
		for _, createTournament := range tournamentResponse.Tournaments {
			tournamentIds = append(tournamentIds, createTournament.ID.String())
		}
		deleteTournaments := dtos.TournamentIds{TournamentIds: tournamentIds}
		err = tournamentService.DeleteTournaments(*newUser.ID, &deleteTournaments)
		if err != nil {
			t.Fatalf("Error in DeleteTournaments. Expected %v, got %v.", nil, err.Error())
		}
	})

}
