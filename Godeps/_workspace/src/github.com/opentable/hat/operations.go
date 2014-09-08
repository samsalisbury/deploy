package hat

import (
	"reflect"
)

type Op struct {
	On             SELF
	Inputs         []IN
	OptionalInputs []IN
	Outputs        []OUT
	Required       func(*Node) bool
}

func (o *Op) MinIn() int {
	return len(o.Inputs)
}

func (o *Op) MaxIn() int {
	return o.MinIn() + len(o.OptionalInputs)
}

func (o *Op) RequiresNilReceiver() bool {
	return o.On == SELF_Nil
}

func (o *Op) RequiresPayloadReceiver() bool {
	return o.On == SELF_Payload
}

func (o *Op) RequiresManifestedReceiver() bool {
	return o.On == SELF_Manifested
}

func (o *Op) RequiresPayload() bool {
	return o.Requires(IN_Payload)
}

func (o *Op) Requires(in IN) bool {
	for _, i := range o.Inputs {
		if i == in {
			return true
		}
	}
	return false
}

type CompiledOp struct {
	Def             *Op
	OtherEntityType reflect.Type
	Method          reflect.Method
	NumIn           int
	In              []IN
	InputTypes      []reflect.Type
	Node            *Node
}

func (o *Op) Compile(n *Node, m reflect.Method) (*CompiledOp, error) {
	var otherEntityType reflect.Type
	numIn := m.Type.NumIn() - 1
	exactInputs := make([]IN, numIn)
	inputTypes := make([]reflect.Type, numIn)
	i := 0
	if numIn < o.MinIn() || numIn > o.MaxIn() {
		return nil, n.wrongNumIn(m.Name, o, numIn)
	}
	// Validate the mandatory inputs of the user defined method (m).
	for userPos, in := range o.Inputs {
		realPos := 1 + userPos
		paramType := m.Type.In(realPos)
		if err := in.Accepts(n, m.Name, userPos, paramType); err != nil {
			return nil, err
		}
		// Special meanings of certain field types are resolved here.
		if in == IN_Payload {
			otherEntityType = paramType
		}
		exactInputs[i] = in
		inputTypes[i] = paramType
		i++
	}
	// Validate the optional inputs of the user defined method (m).
	for userPos, in := range o.OptionalInputs {
		userPos = userPos + len(o.Inputs)
		realPos := 1 + userPos
		if userPos == numIn {
			break
		}
		paramType := m.Type.In(realPos)
		if err := in.Accepts(n, m.Name, userPos, paramType); err != nil {
			return nil, err
		}
		exactInputs[i] = in
		inputTypes[i] = paramType
		i++
	}
	return &CompiledOp{o, otherEntityType, m, numIn, exactInputs, inputTypes, n}, nil
}

// A BoundOp has all data resolved, and a method bound to the
// specified receiver value (nil, payload, or manifested).
type BoundOp struct {
	Compiled *CompiledOp
	Receiver interface{}
	Inputs   map[IN]boundInput
	Method   reflect.Value
}

func (co *CompiledOp) BindNilReceiver() *BoundOp {
	return &BoundOp{Compiled: co, Receiver: reflect.New(co.Node.EntityType).Interface()}
}

func (co *CompiledOp) BindManifestedOrPayloadReciever(rcvr interface{}) (*BoundOp, error) {
	return &BoundOp{Compiled: co, Receiver: rcvr}, nil
}

func (co *CompiledOp) Invoke(inputs map[IN]boundInput) (entity interface{}, other interface{}, err error) {
	if bo, err := co.Bind(inputs); err != nil {
		return nil, nil, err
	} else {
		return bo.Invoke()
	}
}

func (bo *BoundOp) Invoke() (entity interface{}, other interface{}, err error) {
	var in []reflect.Value
	if in, err = bo.PrepareInputs(); err != nil {
		return nil, nil, err
	}
	out := bo.Method.Call(in)
	entity = bo.Receiver
	for i, o := range bo.Compiled.Def.Outputs {
		if o == OUT_Error {
			if !out[i].IsNil() {
				err = out[i].Interface().(error)
			}
		} else if o == OUT_OtherEntity {
			if !out[i].IsNil() {
				other = out[i].Interface()
			}
		}
	}
	return entity, other, err
}

func (bo *BoundOp) PrepareInputs() ([]reflect.Value, error) {
	prepared := make([]reflect.Value, bo.Compiled.NumIn)
	for i, in := range bo.Compiled.In {
		if it, err := bo.Inputs[in](bo); err != nil {
			return nil, err
		} else {
			if it == nil {
				it = reflect.New(bo.Compiled.InputTypes[i])
			}
			prepared[i] = reflect.ValueOf(it)
		}
	}
	return prepared, nil
}

func (co *CompiledOp) Bind(inputs map[IN]boundInput) (*BoundOp, error) {
	if bo, err := co.BindReceiver(inputs); err != nil {
		return nil, err
	} else if err := bo.BindMethod(); err != nil {
		return nil, err
	} else {
		bo.Inputs = inputs
		return bo, nil
	}
}

func (bo *BoundOp) BindMethod() error {
	bo.Method = reflect.ValueOf(bo.Receiver).MethodByName(bo.Compiled.Method.Name)
	return nil
}

func (co *CompiledOp) Error(args ...interface{}) hatError {
	return co.Node.MethodError(co.Method.Name, args...)
}

func (co *CompiledOp) BindReceiver(inputs map[IN]boundInput) (*BoundOp, error) {
	if co.Def.RequiresNilReceiver() {
		return co.BindNilReceiver(), nil
	} else {
		var rcvr interface{}
		var err error
		if co.Def.RequiresManifestedReceiver() {
			rcvr, err = co.BindManifestedReceiver(inputs)
		} else if co.Def.RequiresPayloadReceiver() {
			rcvr, err = co.BindPayloadReceiver(inputs)
		}
		if err != nil {
			return nil, co.Error("BindReceiver(...)", err.Error())
		}
		return co.BindManifestedOrPayloadReciever(rcvr)
	}
	return nil, co.Node.MethodError(co.Method.Name, "no receiver type specified")
}

func (co *CompiledOp) BindManifestedReceiver(inputs map[IN]boundInput) (interface{}, error) {
	rcvr, _, err := co.Node.Ops["Manifest"].Invoke(inputs)
	return rcvr, err
}

func (co *CompiledOp) BindPayloadReceiver(inputs map[IN]boundInput) (interface{}, error) {
	bo := &BoundOp{Compiled: co}
	return inputs[IN_Payload](bo)
}

type SELF int

const (
	SELF_Nil        = SELF(iota)
	SELF_Payload    = SELF(iota)
	SELF_Manifested = SELF(iota) // TODO: Remove; Deprecated: all selves are manifested.
)

type OUT int

const (
	OUT_Error       = OUT(iota)
	OUT_OtherEntity = OUT(iota)
)

func on(self SELF) *Op                         { return &Op{On: self} }
func (o *Op) In(inputs ...IN) *Op              { o.Inputs = inputs; return o }
func (o *Op) OptIn(inputs ...IN) *Op           { o.OptionalInputs = inputs; return o }
func (o *Op) Out(outputs ...OUT) *Op           { o.Outputs = outputs; return o }
func (o *Op) RequireIf(p func(*Node) bool) *Op { o.Required = p; return o }

func iType(nilPtr interface{}) reflect.Type {
	return reflect.TypeOf(nilPtr).Elem()
}

func typ(example interface{}) reflect.Type {
	return reflect.TypeOf(example)
}
