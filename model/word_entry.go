package model

type WordEntry struct {
	Id     uint   `json:"id" gorm:"primaryKey"`
	Word   string `json:"word"`
	Skewer string `json:"skewer"`
}
