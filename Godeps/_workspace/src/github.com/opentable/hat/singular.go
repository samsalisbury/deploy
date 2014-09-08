package hat

type ResolvedSingularNode struct {
	ResolvedNodeBase
	entity interface{} // Manifested singular entity.
}

// Creates both singular and collection resolved nodes that belong to a named member.
func newResolvedSingular(parent ResolvedNode, node *Node, id string, tag *Tag, entity interface{}) ResolvedNode {
	base := newResolvedNodeBase(parent, node, id, tag)
	return &ResolvedSingularNode{base, entity}
}

func (rn *ResolvedSingularNode) Locate(path ...string) (ResolvedNode, error) {
	if len(path) == 0 || (len(path) == 1 && len(path[0]) == 0) {
		return rn, nil
	}
	id := path[0]
	path = path[1:]
	child, err := rn.Resolve(id)
	if err != nil {
		return nil, err
	} else if child == nil {
		return nil, HttpError(404, id, "does not exist")
	}
	return child.Locate(path...)
}

func (n *ResolvedSingularNode) Resolve(id string) (ResolvedNode, error) {
	member, ok := n.Node.Members[id]
	if !ok {
		return nil, HttpError(404, n.Node.EntityType.Name(), "(", n.ID(), ") does not have a member called", quot(id))
	}
	if member.Node.IsCollection {
		return n.ResolveCollection(member.Tag, member.Node, id)
	} else {
		return n.ResolveSingular(member.Tag, member.Node, id)
	}
}

func (n *ResolvedSingularNode) Entity() interface{} {
	return n.entity
}

func (n *ResolvedSingularNode) SetEntity(e interface{}) {
	n.entity = e
}

func (n *ResolvedSingularNode) ResolveCollection(tag *Tag, collectionNode *Node, id string) (ResolvedNode, error) {
	collection, ids, err := collectionNode.ManifestCollection(n.Entity(), id)
	if err != nil {
		return nil, err
	}
	return newResolvedCollection(n, collectionNode, id, tag, collection, ids), nil
}

func (n *ResolvedSingularNode) ResolveSingular(tag *Tag, singularNode *Node, id string) (ResolvedNode, error) {
	entity, err := singularNode.ManifestSingular(n.Entity(), id)
	if err != nil {
		return nil, err
	}
	return newResolvedSingular(n, singularNode, id, tag, entity), nil
}

func (n *ResolvedSingularNode) Links() ([]Link, error) {
	links := []Link{Link{"self", n.Path()}}
	for name, member := range n.Node.Members {
		if member.Tag.Link {
			rel := member.Tag.Rel
			if rel == "" {
				rel = name
			}
			links = append(links, Link{rel, n.Path() + "/" + name})
		}
	}
	return links, nil
}

func (n *ResolvedSingularNode) Resource() (*Resource, error) {
	if entity, err := toSmap(n.entity); err != nil {
		return nil, err
	} else if embeddedMembers, err := n.EmbeddedMembers(); err != nil {
		return nil, err
	} else if links, err := n.Links(); err != nil {
		return nil, err
	} else {
		for k, _ := range embeddedMembers {
			entity.deleteIgnoringCase(k)
		}
		for _, l := range links {
			entity.deleteIgnoringCase(l.Rel)
		}
		return &Resource{n.Tag.Rel, entity, embeddedMembers, nil, links}, nil
	}
}

func (n *ResolvedSingularNode) EmbeddedResource(tag *Tag) (*Resource, error) {
	if len(tag.EmbedFields) != 0 {
		return n.FilteredEmbeddedResource(tag.EmbedFields)
	}
	return n.DefaultEmbeddedResource()
}

func (n *ResolvedSingularNode) DefaultEmbeddedResource() (*Resource, error) {
	return n.Resource()
}

func (n *ResolvedSingularNode) FilteredEmbeddedResource(fields []string) (*Resource, error) {
	r, err := n.Resource()
	if err != nil {
		return nil, err
	}
	m, err := toSmap(r.Entity)
	if err != nil {
		return nil, err
	}
	filtered := make(smap, len(fields))
	for _, f := range fields {
		filtered[f] = m[f]
	}
	r.Entity = filtered
	return r, nil
}

func (n *ResolvedSingularNode) EmbeddedMembers() (map[string]*Resource, error) {
	if n.Node.IsCollection {
		return nil, nil
	}
	embedded := map[string]*Resource{}
	for urlName, member := range n.Node.Members {
		if !member.Tag.Embed {
			continue
		}
		if memberNode, err := n.Locate(urlName); err != nil {
			return nil, err
		} else if resource, err := memberNode.EmbeddedResource(member.Tag); err != nil {
			return nil, err
		} else {
			embedded[member.Name] = resource
		}
	}
	return embedded, nil
}

func (n *ResolvedSingularNode) EmbeddedCollectionItems() ([]*Resource, error) {
	return nil, nil
}
