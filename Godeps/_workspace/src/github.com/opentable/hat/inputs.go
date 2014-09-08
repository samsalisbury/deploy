package hat

import (
	"reflect"
)

type IN int

const (
	IN_Parent   = IN(iota)
	IN_Payload  = IN(iota)
	IN_ID       = IN(iota)
	IN_PageNum  = IN(iota)
	ON_PageSize = IN(iota)
)

func (in IN) Accepts(n *Node, name string, pos int, t reflect.Type) error {
	switch in {
	default:
		panic("The programmer has made a serious error.")

	case IN_ID:
		if t.Kind() != reflect.String {
			return n.MethodError(name, "cannot accept input type", t, "at position", pos)
		}

	case IN_Parent:
		if n.Parent != nil && t != n.Parent.EntityPtrType {
			return n.MethodError(name, "expects a pointer to its parent type", n.Parent.EntityPtrType, "at position", pos)
		}

	case IN_Payload:
		if t.Kind() == reflect.Ptr {
			elemKind := t.Elem().Kind()
			switch elemKind {
			case reflect.Struct, reflect.Map, reflect.Slice:
				return nil // This is the only ok case, otherwise we return the below error.
			}
		}
		return n.MethodError(name, "expects a pointer to a struct, map, or slice at position", pos)

	case IN_PageNum:
		if t.Kind() != reflect.Int {
			return n.MethodError(name, "expects an int at position", pos)
		}
	}
	return nil
}
