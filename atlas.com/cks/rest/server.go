package rest

import (
	"atlas-cks/keymap"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"time"
)

type Server struct {
	l  *logrus.Logger
	hs *http.Server
}

func NewServer(l *logrus.Logger, db *gorm.DB) *Server {
	router := mux.NewRouter().PathPrefix("/ms/cks").Subrouter().StrictSlash(true)
	router.Use(commonHeader)

	csr := router.PathPrefix("/characters").Subrouter()
	csr.HandleFunc("/{characterId}/keymap", keymap.HandleGetKeyMap(l, db)).Methods(http.MethodGet)

	w := l.Writer()
	defer w.Close()

	hs := http.Server{
		Addr:         ":8080",
		Handler:      router,
		ErrorLog:     log.New(w, "", 0), // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}
	return &Server{l, &hs}
}

func (s *Server) Run() {
	s.l.Infoln("Starting server on port 8080")
	err := s.hs.ListenAndServe()
	if err != nil {
		s.l.Errorf("Starting server: %s\n", err)
		os.Exit(1)
	}
}

func commonHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
