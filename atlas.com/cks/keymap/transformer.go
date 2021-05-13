package keymap

func makeKeyMap(e *entity) *Model {
	return &Model{
		id:          e.ID,
		characterId: e.CharacterId,
		key:         e.Key,
		theType:     e.Type,
		action:      e.Action,
	}
}
