package keymap

import (
	"atlas-cks/json"
	"atlas-cks/rest"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

const (
	GetKeyMap   = "get_key_map"
	ResetKeyMap = "reset_key_map"
)

func InitResource(router *mux.Router, l logrus.FieldLogger, db *gorm.DB) {
	csr := router.PathPrefix("/characters").Subrouter()
	csr.HandleFunc("/{characterId}/keymap", registerGetKeyMap(l, db)).Methods(http.MethodGet)
	csr.HandleFunc("/{characterId}/keymap/reset", registerResetKeyMap(l, db)).Methods(http.MethodPost)
}

type IdHandler func(characterId uint32) http.HandlerFunc

func ParseId(l logrus.FieldLogger, next IdHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		value, err := strconv.Atoi(vars["characterId"])
		if err != nil {
			l.WithError(err).Errorln("Error parsing id as uint32")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		next(uint32(value))(w, r)
	}
}

func registerGetKeyMap(l logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return rest.RetrieveSpan(GetKeyMap, handleGetKeyMap(l, db))
}

func handleGetKeyMap(l logrus.FieldLogger, db *gorm.DB) rest.SpanHandler {
	return func(span opentracing.Span) http.HandlerFunc {
		return ParseId(l, func(characterId uint32) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				curriedGetKeyMap(l, db)(characterId)(w)
			}
		})
	}
}

func curriedGetKeyMap(fl logrus.FieldLogger, db *gorm.DB) func(characterId uint32) func(w http.ResponseWriter) {
	return func(characterId uint32) func(w http.ResponseWriter) {
		return func(w http.ResponseWriter) {
			l := fl.WithFields(logrus.Fields{"originator": "HandleGetKeyMap", "type": "rest_handler"})
			keys, err := GetKeyMapForCharacter(l, db)(characterId)
			if err != nil {
				l.WithError(err).Errorf("Unable to retrieve keybindings for character %d.", characterId)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			result := DataListContainer{Data: make([]DataBody, 0)}
			for _, key := range keys {
				result.Data = append(result.Data, DataBody{
					Id:   strconv.Itoa(int(key.Id())),
					Type: "KeyMap",
					Attributes: Attributes{
						Key:    key.Key(),
						Type:   key.Type(),
						Action: key.Action(),
					},
				})
			}

			w.WriteHeader(http.StatusOK)
			err = json.ToJSON(result, w)
			if err != nil {
				fl.WithError(err).Errorf("Writing response.")
			}
		}
	}
}

func registerResetKeyMap(l logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return rest.RetrieveSpan(ResetKeyMap, handleResetKeyMap(l, db))
}

func handleResetKeyMap(l logrus.FieldLogger, db *gorm.DB) rest.SpanHandler {
	return func(span opentracing.Span) http.HandlerFunc {
		return ParseId(l, func(characterId uint32) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				curriedResetKeyMap(l, db)(characterId)(w)
			}
		})
	}
}

func curriedResetKeyMap(fl logrus.FieldLogger, db *gorm.DB) func(characterId uint32) func(w http.ResponseWriter) {
	return func(characterId uint32) func(w http.ResponseWriter) {
		return func(w http.ResponseWriter) {
			l := fl.WithFields(logrus.Fields{"originator": "HandleResetKeyMap", "type": "rest_handler"})

			err := Reset(l, db)(characterId)
			if err != nil {
				l.WithError(err).Errorf("Unable to reset keybindings for character %d.", characterId)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusNoContent)
		}
	}
}
