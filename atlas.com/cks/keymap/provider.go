package keymap

import (
	"atlas-cks/database"
	"atlas-cks/model"
	"gorm.io/gorm"
)

func entityByCharacterId(id uint32) database.EntityListProvider[entity] {
	return func(db *gorm.DB) model.SliceProvider[entity] {
		return database.ListGet[entity](db, &entity{CharacterId: id})
	}
}
