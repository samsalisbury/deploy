package hat

import "reflect"

type StdHTTPMethod func() (statusCode int, resource *Resource, err error)

func makeHTTPMethods(n ResolvedNode, inputs map[IN]boundInput) map[string]StdHTTPMethod {
	methods := map[string]StdHTTPMethod{
		"GET": makeGET(n, inputs),
	}
	if _, ok := n.UnderlyingNode().Ops["Write"]; ok {
		methods["PUT"] = makePUT(n, inputs)
	}
	return methods
}

func makeGET(n ResolvedNode, inputs map[IN]boundInput) StdHTTPMethod {
	return func() (statusCode int, resource *Resource, err error) {
		if !n.UnderlyingNode().IsCollection && reflect.DeepEqual(reflect.ValueOf(n.Entity()).Elem().Interface(), reflect.Zero(n.UnderlyingNode().EntityType).Interface()) {
			return 0, nil, HttpError(404, "Not found.")
		}
		if r, err := n.Resource(); err != nil {
			return 0, nil, err
		} else {
			return 200, r, nil
		}
	}
}

func makePUT(n ResolvedNode, inputs map[IN]boundInput) StdHTTPMethod {
	return func() (statusCode int, resource *Resource, err error) {
		successStatus := 200
		if n.Entity() == nil {
			successStatus = 201
		}
		op := n.UnderlyingNode().Ops["Write"]
		_, _, err = op.Invoke(inputs)
		//n.SetEntity(entity)
		entity, err := n.UnderlyingNode().Manifest(n.Parent().Entity(), n.ID())
		if err != nil {
			return 0, nil, err
		}
		n.SetEntity(entity)
		if err != nil {
			return 0, nil, err
		}
		if r, err := n.Resource(); err != nil {
			return 0, nil, err
		} else {
			return successStatus, r, nil
		}
	}
}

func (parentNode *Node) createChildManifestInputs(parentEntity interface{}, id string) map[IN]boundInput {
	if parentEntity == nil && parentNode != nil {
		parentEntity = reflect.New(parentNode.EntityType).Interface()
	}
	return map[IN]boundInput{
		IN_Parent: func(_ *BoundOp) (interface{}, error) {
			return parentEntity, nil
		},
		IN_ID: func(_ *BoundOp) (interface{}, error) {
			return id, nil
		},
		IN_PageNum: func(_ *BoundOp) (interface{}, error) {
			return 1, nil
		},
	}
}
