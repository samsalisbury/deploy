package hat

import (
	"encoding/json"
	"reflect"
	"strings"
)

// smap is a convenient handle on maps with string keys it is used throughout
// hat.
type smap map[string]interface{}

// This is only necessary because the toSmap method below
// uses JSON encoder to convert to smap, which means
func (m *smap) deleteIgnoringCase(k string) {
	for realKey, _ := range *m {
		if strings.EqualFold(k, realKey) {
			delete(*m, realKey)
		}
	}
}

// Converts the collection to a smap or panics.
func collectionToSmap(collection interface{}) smap {
	v := reflect.ValueOf(collection).Elem()
	m := smap{}
	for _, kv := range v.MapKeys() {
		m[kv.String()] = v.MapIndex(kv).Interface()
	}
	return m
}

// toSmap should convert anything (except a slice) to an smap. Presently it
// leans on JSON encoding to make this conversion. This is really bad, as
// json encoding is a separate concern, so this needs to be addressed soon.
func toSmap(v interface{}) (smap, error) {
	// THIS IS VERY BAD since JSON encoding has too many other rules.
	// TODO: Find an alternative... Probably hand-roll some reflection.
	return toSmapRespectingJSONTags(v)
}

// This is useful for rendering HAL
func toSmapRespectingJSONTags(v interface{}) (smap, error) {
	if j, err := json.Marshal(v); err != nil {
		return nil, err
	} else {
		var m map[string]interface{}
		if err := json.Unmarshal(j, &m); err != nil {
			return nil, err
		}
		return m, nil
	}
}
