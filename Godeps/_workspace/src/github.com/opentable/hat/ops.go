package hat

var op_specs = map[string]*Op{
	"Manifest": on(SELF_Nil).In().OptIn(IN_Parent, IN_ID).Out(OUT_Error).
		RequireIf(func(n *Node) bool { return !n.IsCollection }),
	"Page": on(SELF_Nil).In().OptIn(IN_PageNum, IN_Parent, IN_ID).Out(OUT_OtherEntity, OUT_Error).
		RequireIf(func(n *Node) bool { return n.IsCollection }),
	"Write": on(SELF_Payload).In().OptIn(IN_Parent, IN_ID).Out(OUT_Error).
		RequireIf(func(n *Node) bool { return false }),
}
