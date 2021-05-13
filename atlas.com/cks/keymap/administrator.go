package keymap

import "gorm.io/gorm"

func create(db *gorm.DB, characterId uint32, key int32, theType int8, action int32) (*Model, error) {
	e := &entity{
		CharacterId: characterId,
		Key:         key,
		Type:        theType,
		Action:      action,
	}

	err := db.Create(e).Error
	if err != nil {
		return nil, err
	}
	return makeKeyMap(e), nil
}

func deleteByCharacter(db *gorm.DB, characterId uint32) error {
	return db.Where(&entity{CharacterId: characterId}).Delete(&entity{}).Error
}
