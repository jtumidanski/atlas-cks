package rest

import (
	"atlas-cks/keymap"
	"context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"sync"
)

func CreateRestService(l *logrus.Logger, db *gorm.DB, ctx context.Context, wg *sync.WaitGroup) {
	go NewServer(l, ctx, wg, ProduceRoutes(db))
}

func ProduceRoutes(db *gorm.DB) func(l logrus.FieldLogger) http.Handler {
	return func(l logrus.FieldLogger) http.Handler {
		router := mux.NewRouter().PathPrefix("/ms/cks").Subrouter().StrictSlash(true)
		router.Use(CommonHeader)

		csr := router.PathPrefix("/characters").Subrouter()
		csr.HandleFunc("/{characterId}/keymap", keymap.HandleGetKeyMap(l, db)).Methods(http.MethodGet)
		csr.HandleFunc("/{characterId}/keymap/reset", keymap.HandleResetKeyMap(l, db)).Methods(http.MethodPost)

		return router
	}
}
