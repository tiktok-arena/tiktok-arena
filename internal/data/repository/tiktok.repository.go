package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"tiktok-arena/internal/core/models"
)

type TiktokRepository struct {
	db *gorm.DB
}

func NewTiktokRepository(db *gorm.DB) *TiktokRepository {
	return &TiktokRepository{db: db}
}

func (r *TiktokRepository) CreateNewTiktok(newTiktok models.Tiktok) error {
	record := r.db.
		Create(&newTiktok)
	return record.Error
}

func (r *TiktokRepository) CreateNewTiktoks(t []models.Tiktok) error {
	record := r.db.
		Create(t)
	return record.Error
}

func (r *TiktokRepository) EditTiktok(t models.Tiktok) error {
	record := r.db.
		Where(&t.URL, &t.TournamentID).
		Updates(&t)
	return record.Error
}

func (r *TiktokRepository) DeleteTiktoks(t []models.Tiktok) error {
	record := r.db.Delete(t)
	return record.Error
}

func (r *TiktokRepository) DeleteTiktoksByIds(ids []string) error {
	record := r.db.
		Where("tournament_id IN (?)", ids).
		Delete(&models.Tiktok{})
	return record.Error
}

func (r *TiktokRepository) GetTournamentTiktoksById(tournamentId uuid.UUID) ([]models.Tiktok, error) {
	var tiktoks []models.Tiktok
	record := r.db.
		Select("*").
		Find(&tiktoks, "tournament_id = ?", tournamentId)
	return tiktoks, record.Error
}

func (r *TiktokRepository) UpdateTiktokWins(tournamentId uuid.UUID, tiktokURL string) error {
	record := r.db.
		Model(&models.Tiktok{}).
		Where("tournament_id = ? AND url = ?", tournamentId, tiktokURL).
		UpdateColumn("wins", gorm.Expr("wins + ?", 1))
	return record.Error
}
