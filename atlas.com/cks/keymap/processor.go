package keymap

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func GetKeyMapForCharacter(_ logrus.FieldLogger, db *gorm.DB) func(characterId uint32) ([]*Model, error) {
	return func(characterId uint32) ([]*Model, error) {
		return getKeyMapForCharacter(db, characterId)
	}
}

func ChangeKeyMap(l logrus.FieldLogger, db *gorm.DB) func(characterId uint32, changes []Change) error {
	return func(characterId uint32, changes []Change) error {
		return db.Transaction(func(tx *gorm.DB) error {
			err := deleteByCharacter(tx, characterId)
			if err != nil {
				l.WithError(err).Errorf("Unable to delete existing bindings for character %d.", characterId)
				return err
			}

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
