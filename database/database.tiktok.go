package database

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"tiktok-arena/models"
)

func CreateNewTiktok(newTiktok *models.Tiktok) error {
	record := tiktokTable.Create(&newTiktok)
	return record.Error
}

func CreateNewTiktoks(t []models.Tiktok) error {
	record := tiktokTable.Create(t)
	return record.Error
}

func EditTiktok(t *models.Tiktok) error {
	record := tiktokTable.
		Where(&t.URL, &t.TournamentID).
		Updates(&t)
	return record.Error
}

func DeleteTiktoks(t []models.Tiktok) error {
	record := tiktokTable.Delete(t)
	return record.Error
}

func DeleteTiktoksByIds(ids []string) error {
	record := tiktokTable.
		Where("tournament_id IN (?)", ids).
		Delete(&models.Tiktok{})
	return record.Error
}

func GetTournamentTiktoksById(tournamentId uuid.UUID) ([]models.Tiktok, error) {
	var tiktoks []models.Tiktok
	record := tiktokTable.
		Select("*").
		Find(&tiktoks, "tournament_id = ?", tournamentId)
	return tiktoks, record.Error
}

func UpdateTiktokWins(tournamentId uuid.UUID, tiktokURL string) error {
	record := tiktokTable.
		Where("tournament_id = ? AND url = ?", tournamentId, tiktokURL).
		UpdateColumn("wins", gorm.Expr("wins + ?", 1))
	return record.Error
}
