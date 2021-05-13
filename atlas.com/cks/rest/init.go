package rest

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func CreateRestService(l *logrus.Logger, db *gorm.DB) {
	rs := NewServer(l, db)
	go rs.Run()
}
