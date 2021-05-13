package keymap

import "gorm.io/gorm"

func getKeyMapForCharacter(db *gorm.DB, characterId uint32) ([]*Model, error) {
	return listGet(db, &entity{CharacterId: characterId})
}

func listGet(db *gorm.DB, query interface{}) ([]*Model, error) {
	var results []entity
	err := db.Where(query).Find(&results).Error
	if err != nil {
		return nil, err
	}

	var character = make([]*Model, 0)
	for _, e := range results {
		character = append(character, makeKeyMap(&e))
	}
	return character, nil
}