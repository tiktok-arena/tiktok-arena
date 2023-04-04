package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"tiktok-arena/models"
)

func GetTournamentById(tournamentId string) (models.Tournament, error) {
	var tournament models.Tournament
	record := tournamentsTable.First(&tournament, "id = ?", tournamentId)
	return tournament, record.Error
}

func CheckIfTournamentExistsByName(name string) (bool, error) {
	var tournament models.Tournament
	record := tournamentsTable.Select("id").First(&tournament, "name = ?", name)
	return tournament.ID != nil, record.Error
}

func CheckIfNameIsTakenByOtherTournament(name string, id uuid.UUID) (bool, error) {
	var tournament models.Tournament
	record := tournamentsTable.Select("id").First(&tournament, "name = ? AND id != ?", name, id)
	return tournament.ID != nil, record.Error
}

func CheckIfTournamentExistsById(id uuid.UUID) (bool, error) {
	var tournament models.Tournament
	record := tournamentsTable.Select("id").First(&tournament, "id = ?", id)
	return tournament.ID != nil, record.Error
}

func CheckIfTournamentsExistsByIds(ids []string, userId uuid.UUID) (bool, error) {
	var tournaments []models.Tournament
	record := tournamentsTable.Where("user_id = ? AND id IN ?", userId, ids).Find(&tournaments)
	if len(tournaments) != len(ids) {
		return false, record.Error
	}
	return true, record.Error
}

func CreateNewTournament(newTournament *models.Tournament) error {
	record := tournamentsTable.Create(&newTournament)
	return record.Error
}

func EditTournament(t *models.Tournament) error {
	record := tournamentsTable.Where("id = ?", &t.ID).Updates(t)
	return record.Error
}

func DeleteTournamentById(id uuid.UUID, userId uuid.UUID) error {
	record := tournamentsTable.Where("id = ? AND user_id = ?", id, userId).Delete(&models.Tournament{})
	return record.Error
}

func DeleteTournamentsByIds(ids []string, userId uuid.UUID) error {
	record := tournamentsTable.
		Where("user_id = ? AND id IN (?)", userId, ids).
		Delete(&models.Tournament{})
	return record.Error
}

func GetTournaments(totalTournaments int64, queries models.PaginationQueries) (models.TournamentsResponse, error) {
	var tournaments []models.Tournament
	record := tournamentsTable.
		Scopes(Search(queries.SearchText)).
		Scopes(Paginate(queries.Page, queries.Count)).
		Find(&tournaments)
	return models.TournamentsResponse{TournamentCount: totalTournaments, Tournaments: tournaments}, record.Error
}

func GetAllTournamentsForUserById(id uuid.UUID, totalTournaments int64, queries models.PaginationQueries) (models.TournamentsResponse, error) {
	var tournaments []models.Tournament
	record := tournamentsTable.
		Where("user_id = ?", id).
		Scopes(Search(queries.SearchText)).
		Scopes(Paginate(queries.Page, queries.Count)).
		Limit(100).Find(&tournaments)
	return models.TournamentsResponse{TournamentCount: totalTournaments, Tournaments: tournaments}, record.Error
}

func TotalTournaments() (int64, error) {
	var totalTournaments int64
	record := tournamentsTable.Count(&totalTournaments)
	return totalTournaments, record.Error
}

func TotalTournamentsByUserId(id uuid.UUID) (int64, error) {
	var totalTournaments int64
	record := tournamentsTable.Where("user_id = ?", id).Count(&totalTournaments)
	return totalTournaments, record.Error
}

func UpdateTournamentTimesPlayed(tournamentId uuid.UUID) error {
	record := tournamentsTable.Where("id = ?", tournamentId).
		UpdateColumn("times_played", gorm.Expr("times_played + ?", 1))
	return record.Error
}
