package hat

import (
	"strings"
)

type Member struct {
	Name string
	Node *Node
	Tag  *Tag
}

func newMember(name string, n *Node, tag *Tag) (*Member, error) {
	return &Member{name, n, tag}, nil
}

func (m *Member) URLName() string {
	return strings.ToLower(m.Name)
}
