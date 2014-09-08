package logging

import (
	j "encoding/json"
)

func json(a interface{}) ([]byte, error) {
	data, err := j.Marshal(a)
	if err != nil {
		return nil, Error("json.Marshal failed. " + err.Error())
	}
	debug("JSON Length:", len(data))
	return data, nil
}
