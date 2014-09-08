package hat

func (n *Node) initOps() error {
	n.Ops = map[string]*CompiledOp{}
	for name, op := range op_specs {
		if m, ok := n.EntityPtrType.MethodByName(name); !ok {
			if op.Required(n) {
				return n.Error("requires", name, "method")
			}
			continue
		} else if co, err := op.Compile(n, m); err != nil {
			return err
		} else {
			n.Ops[name] = co
		}
	}
	return nil
}
