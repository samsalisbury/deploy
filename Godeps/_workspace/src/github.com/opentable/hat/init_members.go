package hat

import (
	"reflect"
)

func (n *Node) initMembers() error {
	switch n.EntityType.Kind() {
	case reflect.Struct:
		return n.initStructMembers()
	case reflect.Map:
		return n.initMapCollection()
	default:
		return Error("Nodes with kind " + n.EntityType.Kind().String() + " are not yet supported.")
	}
}

func (n *Node) initMapCollection() error {
	keyType := n.EntityType.Key()
	elementType := n.EntityType.Elem()
	if keyType.Kind() != reflect.String {
		return Error("Only collections with string keys are currently supported.")
	}
	if childNode, err := newNode(n, elementType, n.Tag); err != nil {
		return err
	} else {
		n.Collection = &Member{Node: childNode, Tag: n.Tag}
		n.CollectionName = elementType.Name()
	}
	return nil
}

func (n *Node) initStructMembers() error {
	n.Members = map[string]*Member{}
	t := n.EntityType
	numFields := t.NumField()
	for i := 0; i < numFields; i++ {
		f := t.Field(i)
		if tagData := f.Tag.Get("hat"); tagData != "" {
			// Embedded and linked items must be pointers. This makes the code in hat
			// easier to write.
			if f.Type.Kind() != reflect.Ptr {
				return n.MethodError(f.Name, "is", f.Type, "should be", reflect.PtrTo(f.Type), "because it is tagged hat:")
			}
			tag, err := newTag(f.Type.Elem().Name(), tagData)
			if err != nil {
				return err
			}
			if childNode, err := newNode(n, f.Type, tag); err != nil {
				return err
			} else if member, err := newMember(f.Name, childNode, tag); err != nil {
				return err
			} else {
				n.Members[member.URLName()] = member
			}
		}
	}
	return nil
}
