package keymap

import (
	"atlas-cks/database"
	"atlas-cks/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var defaultKey = []int32{18, 65, 2, 23, 3, 4, 5, 6, 16, 17, 19, 25, 26, 27, 31, 34, 35, 37, 38, 40, 43, 44, 45, 46, 50, 56, 59, 60, 61, 62, 63, 64, 57, 48, 29, 7, 24, 33, 41, 39}
var defaultType = []int8{4, 6, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 4, 5, 5, 4, 4, 5, 6, 6, 6, 6, 6, 6, 5, 4, 5, 4, 4, 4, 4, 4}
var defaultAction = []int32{0, 106, 10, 1, 12, 13, 18, 24, 8, 5, 4, 19, 14, 15, 2, 17, 11, 3, 20, 16, 9, 50, 51, 6, 7, 53, 100, 101, 102, 103, 104, 105, 54, 22, 52, 21, 25, 26, 23, 27}

func ByCharacterModelProvider(db *gorm.DB) func(characterId uint32) model.SliceProvider[Model] {
	return func(characterId uint32) model.SliceProvider[Model] {
		return database.ModelSliceProvider[Model, entity](db)(entityByCharacterId(characterId), makeKeyMap)
	}
}

func GetKeyMapForCharacter(_ logrus.FieldLogger, db *gorm.DB) func(characterId uint32) ([]Model, error) {
	return func(characterId uint32) ([]Model, error) {
		return ByCharacterModelProvider(db)(characterId)()
	}
}

func Reset(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32) error {
	return func(characterId uint32) error {
		return db.Transaction(func(tx *gorm.DB) error {
			err := deleteByCharacter(tx, characterId)
			if err != nil {
				l.WithError(err).Errorf("Unable to delete for character %d.", characterId)
				return err
			}
			for i := 0; i < len(defaultKey); i++ {
				_, err := create(tx, characterId, defaultKey[i], defaultType[i], defaultAction[i])
				if err != nil {
					l.WithError(err).Errorf("Unable to create key binding for character %d. key = %d type = %d action = %d.", characterId, defaultKey[i], defaultType[i], defaultAction[i])
					return err
				}
			}
			return nil
		})
	}
}

func CreateDefault(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32) error {
	return func(characterId uint32) error {
		return db.Transaction(func(tx *gorm.DB) error {
			for i := 0; i < len(defaultKey); i++ {
				_, err := create(tx, characterId, defaultKey[i], defaultType[i], defaultAction[i])
				if err != nil {
					l.WithError(err).Errorf("Unable to create key binding for character %d. key = %d type = %d action = %d.", characterId, defaultKey[i], defaultType[i], defaultAction[i])
					return err
				}
			}
			return nil
		})
	}
}

func ChangeKeyMap(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, changes []Change) error {
	return func(characterId uint32, changes []Change) error {
		return db.Transaction(func(tx *gorm.DB) error {
			for _, change := range changes {
				_, err := create(tx, characterId, change.Key, change.ChangeType, change.Action)
				if err != nil {
					l.WithError(err).Errorf("Unable to create key binding for character %d. key = %d type = %d action = %d.", characterId, change.Key, change.ChangeType, change.Action)
					return err
				}
			}
			return nil
		})
	}
}
