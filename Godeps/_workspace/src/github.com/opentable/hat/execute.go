package hat

import (
	"strings"
)

type inBinder func(ResolvedNode) map[IN]boundInput

func ExecuteRequest(root *Node, path string, method string, inputBinder inBinder) (int, *Resource, error) {
	target, err := LocateFromRoot(root, strings.Split(path[1:], "/")...)
	if err != nil {
		return 0, nil, err
	}
	methods := makeHTTPMethods(target, inputBinder(target))

	if m, ok := methods[method]; !ok {
		return 0, nil, HttpError(405, path+" does not support method "+method+"; it does support: "+supportedMethods(methods))
	} else {
		return m()
	}
}

func supportedMethods(methods map[string]StdHTTPMethod) string {
	m := []string{}
	for k, _ := range methods {
		m = append(m, k)
	}
	return strings.Join(m, ", ")
}
