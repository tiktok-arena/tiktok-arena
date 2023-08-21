package models

import (
	"github.com/google/uuid"
	"math/rand"
	"time"
)

type Tiktok struct {
	TournamentID uuid.UUID  `gorm:"not null;primaryKey;default:null" json:"tournamentID"`
	Tournament   Tournament `gorm:"foreignKey:TournamentID" json:"tournament"`
	Name         string     `gorm:"not null;default:null" json:"name"`
	URL          string     `gorm:"not null;primaryKey;default:null" json:"url"`
	Wins         int        `json:"wins"`
}

func FindDifferenceOfTwoTiktokSlices(s1 []Tiktok, s2 []Tiktok) []Tiktok {
	var dif []Tiktok
	for _, t1 := range s1 {
		existsInS2 := false
		for _, t2 := range s2 {
			if t1.TournamentID.String() == t2.TournamentID.String() && t1.URL == t2.URL {
				existsInS2 = true
				break
			}
		}
		if !existsInS2 {
			dif = append(dif, t1)
		}
	}
	return dif
}

func ContainsTiktok(slice []Tiktok, t Tiktok) bool {
	for _, item := range slice {
		if item.TournamentID == t.TournamentID && item.URL == t.URL {
			return true
		}
	}
	return false
}

func ShuffleTiktok(t []Tiktok) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(t), func(i, j int) { t[i], t[j] = t[j], t[i] })
}
