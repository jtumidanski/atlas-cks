package keymap

import "gorm.io/gorm"

func Migration(db *gorm.DB) error {
	return db.AutoMigrate(&entity{})
}

type entity struct {
	ID          uint32 `gorm:"primaryKey;autoIncrement;not null"`
	CharacterId uint32 `gorm:"not null"`
	Key         int32  `gorm:"not null"`
	Type        int8   `gorm:"not null"`
	Action      int32  `gorm:"not null"`
}

func (e entity) TableName() string {
	return "keymap"
}