package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"tiktok-arena/internal/core/dtos"
	"tiktok-arena/internal/core/models"
	"tiktok-arena/internal/data/repository/scopes"
)

type TournamentRepository struct {
	db *gorm.DB
}

func NewTournamentRepository(db *gorm.DB) *TournamentRepository {
	return &TournamentRepository{db: db}
}

func (r *TournamentRepository) GetTournamentWithUserById(tournamentId uuid.UUID) (models.Tournament, error) {
	var tournament *models.Tournament
	record := r.db.
		Preload("User").
		First(&tournament, "id = ?", tournamentId)
	return *tournament, record.Error
}

func (r *TournamentRepository) GetAllTournamentsWithUsers(totalTournaments int64, queries dtos.PaginationQueries) (dtos.TournamentsResponse, error) {
	var tournaments []models.Tournament
	record := r.db.
		Preload("User").
		Scopes(scopes.Private(false)).
		Scopes(scopes.Search(queries.SearchText)).
		Scopes(scopes.Paginate(queries.Page, queries.Count)).
		Find(&tournaments)
	return dtos.TournamentsResponse{TournamentCount: totalTournaments, Tournaments: tournaments}, record.Error
}

func (r *TournamentRepository) CheckIfTournamentExistsByName(name string) (bool, error) {
	var tournament models.Tournament
	record := r.db.
		Select("id").
		First(&tournament, "name = ?", name)
	if record.Error == gorm.ErrRecordNotFound {
		return false, nil
	}
	return tournament.ID != uuid.Nil, record.Error
}

func (r *TournamentRepository) CheckIfNameIsTakenByOtherTournament(name string, id uuid.UUID) (bool, error) {
	var tournament models.Tournament
	record := r.db.
		Select("id").
		First(&tournament, "name = ? AND id != ?", name, id)
	return tournament.ID != uuid.Nil, record.Error
}

func (r *TournamentRepository) CheckIfTournamentExistsById(id uuid.UUID) (bool, error) {
	var tournament models.Tournament
	record := r.db.
		Select("id").
		First(&tournament, "id = ?", id)
	return tournament.ID != uuid.Nil, record.Error
}

func (r *TournamentRepository) CheckIfTournamentsExistsByIds(ids []string, userId uuid.UUID) (bool, error) {
	var tournaments []models.Tournament
	record := r.db.
		Where("user_id = ? AND id IN ?", userId, ids).
		Find(&tournaments)
	if len(tournaments) != len(ids) {
		return false, record.Error
	}
	return true, record.Error
}

func (r *TournamentRepository) CreateNewTournament(newTournament models.Tournament) error {
	record := r.db.
		Create(&newTournament)
	return record.Error
}

func (r *TournamentRepository) EditTournament(t models.Tournament) error {
	record := r.db.
		Model(&models.Tournament{}).
		Where("id = ?", &t.ID).
		Updates(t)
	return record.Error
}

func (r *TournamentRepository) DeleteTournamentById(id uuid.UUID, userId uuid.UUID) error {
	record := r.db.
		Where("id = ? AND user_id = ?", id, userId).
		Delete(&models.Tournament{})
	return record.Error
}

func (r *TournamentRepository) DeleteTournamentsByIds(ids []string, userId uuid.UUID) error {
	record := r.db.
		Where("user_id = ? AND id IN (?)", userId, ids).
		Delete(&models.Tournament{})
	return record.Error
}

func (r *TournamentRepository) GetTournamentsByUserID(id uuid.UUID, totalTournaments int64, queries dtos.PaginationQueries, isPrivate bool) (dtos.TournamentsResponseWithUser, error) {
	var tournaments []dtos.TournamentWithoutUser
	record := r.db.
		Model(models.Tournament{}).
		Where("user_id = ?", id).
		Scopes(scopes.Private(isPrivate)).
		Scopes(scopes.Search(queries.SearchText)).
		Scopes(scopes.Paginate(queries.Page, queries.Count)).
		Find(&tournaments)
	return dtos.TournamentsResponseWithUser{TournamentCount: totalTournaments, Tournaments: tournaments}, record.Error
}

func (r *TournamentRepository) TotalTournaments(isPrivate bool) (int64, error) {
	var totalTournaments int64
	record := r.db.
		Model(&models.Tournament{}).
		Scopes(scopes.Private(isPrivate)).
		Count(&totalTournaments)
	return totalTournaments, record.Error
}

func (r *TournamentRepository) TotalTournamentsByUserId(id uuid.UUID, isPrivate bool) (int64, error) {
	var totalTournaments int64
	record := r.db.
		Model(&models.Tournament{}).
		Scopes(scopes.Private(isPrivate)).
		Where("user_id = ?", id).
		Count(&totalTournaments)
	return totalTournaments, record.Error
}

func (r *TournamentRepository) UpdateTournamentTimesPlayed(tournamentId uuid.UUID) error {
	record := r.db.
		Model(&models.Tournament{}).
		Where("id = ?", tournamentId).
		UpdateColumn("times_played", gorm.Expr("times_played + ?", 1))
	return record.Error
}
