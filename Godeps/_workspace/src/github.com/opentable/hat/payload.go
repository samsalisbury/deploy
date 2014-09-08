package hat

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
)

type Payload struct {
	Body        io.ReadCloser
	ContentType string
}

func newPayload(r *http.Request) *Payload {
	return &Payload{r.Body, r.Header.Get("Content-Type")}
}

func (p *Payload) Manifest(t reflect.Type) (interface{}, error) {
	v := reflect.New(t).Interface()
	if data, err := ioutil.ReadAll(p.Body); err != nil {
		return nil, err
	} else if err := json.Unmarshal(data, &v); err != nil {
		return nil, err
	} else {
		return v, nil
	}
}
