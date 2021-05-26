package keymap

import (
	"atlas-cks/json"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func HandleGetKeyMap(fl logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := fl.WithFields(logrus.Fields{"originator": "HandleGetKeyMap", "type": "rest_handler"})

		characterId, err := strconv.Atoi(mux.Vars(r)["characterId"])
		if err != nil {
			l.WithError(err).Errorf("Unable to properly parse characterId from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		keys, err := GetKeyMapForCharacter(l, db)(uint32(characterId))
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

func HandleResetKeyMap(fl logrus.FieldLogger, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := fl.WithFields(logrus.Fields{"originator": "HandleResetKeyMap", "type": "rest_handler"})

		characterId, err := strconv.Atoi(mux.Vars(r)["characterId"])
		if err != nil {
			l.WithError(err).Errorf("Unable to properly parse characterId from path.")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = Reset(l, db)(uint32(characterId))
		if err != nil {
			l.WithError(err).Errorf("Unable to reset keybindings for character %d.", characterId)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
