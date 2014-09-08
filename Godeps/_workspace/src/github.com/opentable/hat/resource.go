package hat

type Resource struct {
	Rel                     string
	Entity                  interface{}
	EmbeddedMembers         map[string]*Resource
	EmbeddedCollectionItems []*Resource
	Links                   []Link
}
