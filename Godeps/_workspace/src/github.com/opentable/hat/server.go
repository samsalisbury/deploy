package hat

import (
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type Server struct {
	root *Node
}

func NewServer(root interface{}) (*Server, error) {
	if rootNode, err := newNode(nil, reflect.TypeOf(root), &Tag{}); err != nil {
		return nil, err
	} else {
		return &Server{rootNode}, nil
	}
}

func getMediaType(accept string) string {
	parts := strings.Split(accept, ";")
	parts[0] = strings.Trim(parts[0], " \t")
	return parts[0]
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer timeTrack(now(), "request")
	inputBinder := makeInputBinder(r)
	//acceptType := getMediaType(r.Header.Get("Accept"))
	fieldFilters := NewFieldFilter(r.URL.Query().Get("fields"))
	//embedFilters := NewEmbedFilter(r.URL.Query().Get("embed"))
	if statusCode, resource, err := ExecuteRequest(s.root, r.URL.Path, r.Method, inputBinder); err != nil {
		writeError(w, err)
	} else if hal, err := RenderAsHAL(resource, fieldFilters); err != nil {
		writeError(w, err)
	} else {
		halMimeType := "application/hal+json"
		w.Header().Add("Content-Type", halMimeType)
		writeResponse(w, statusCode, hal)
	}
}

type boundInput func(*BoundOp) (interface{}, error)

func makeInputBinder(r *http.Request) func(ResolvedNode) map[IN]boundInput {
	return func(n ResolvedNode) map[IN]boundInput {
		return map[IN]boundInput{
			IN_Parent: func(bo *BoundOp) (interface{}, error) {
				return n.Parent().Entity(), nil
			},
			IN_ID: func(*BoundOp) (interface{}, error) {
				return n.ID(), nil
			},
			IN_Payload: func(bo *BoundOp) (interface{}, error) {
				// TODO: These methods should all be compiled at the compile step.
				payload := newPayload(r)
				var t reflect.Type
				if bo.Compiled.Def.RequiresPayloadReceiver() {
					t = bo.Compiled.Node.EntityType
				} else if bo.Compiled.Def.Requires(IN_Payload) {
					t = bo.Compiled.OtherEntityType
				} else {
					return nil, bo.Compiled.Error("bindInputs: unable to determine required payload type")
				}
				return payload.Manifest(t)
			},
			IN_PageNum: func(*BoundOp) (interface{}, error) {
				if n := r.URL.Query().Get("page"); len(n) == 0 {
					return 0, nil
				} else if page, err := strconv.ParseInt(n, 10, 32); err != nil {
					return 0, Error("Page number", quot(n), "not recognised; expected integer:", err)
				} else {
					return page, nil
				}
			},
		}
	}
}

func writeError(w http.ResponseWriter, err error) {
	if httpErr, ok := err.(HTTPError); ok {
		writeResponse(w, httpErr.StatusCode(), httpErr.Err())
	} else {
		writeResponse(w, 500, errorResource(err))
	}
}

func errorResource(err error) smap {
	return smap{"error": err.Error()}
}

func writeResponse(w http.ResponseWriter, statusCode int, resource interface{}) {
	if data, err := json.Marshal(resource); err != nil {
		writeError(w, HttpError(500, "Unable to marshal response into json:", err.Error()))
	} else {
		w.WriteHeader(statusCode)
		w.Write(data)
	}
}
