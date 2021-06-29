package modules

import (
	"github.com/bSkracic/similaritipsum/db"
	"github.com/bSkracic/similaritipsum/model"
	"gorm.io/gorm"
)

func GetWordEntries(page int, pageSize int) *[]model.WordEntry {

	var wordEntries *[]model.WordEntry
	db.GetConnection().Scopes(Paginate(page, pageSize)).Find(&wordEntries)
	return wordEntries
}

func GetWordEntry(id uint) *model.WordEntry {
	var wordEntry *model.WordEntry
	db.GetConnection().WordEntries().Find(wordEntry, id)
	return wordEntry
}

func CreateWordEntry(wordEntry *model.WordEntry) uint {
	db.GetConnection().WordEntries().Create(wordEntry)
	return wordEntry.Id
}

func Paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}
