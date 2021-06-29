package modules

import (
	"github.com/bSkracic/similaritipsum/db"
	"github.com/bSkracic/similaritipsum/model"
)

func GetWordEntries() *[]model.WordEntry {
	var wordEntries *[]model.WordEntry
	db.GetConnection().WordEntries().Find(wordEntries)
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
