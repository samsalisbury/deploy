package hat

import "strings"

type ResolvedCollectionNode struct {
	ResolvedNodeBase
	Collection    interface{} // Manifested collection.
	CollectionIDs []string    // Manifested collection IDs.
}

func newResolvedCollection(parent ResolvedNode, node *Node, id string, tag *Tag, collection interface{}, ids []string) ResolvedNode {
	base := newResolvedNodeBase(parent, node, id, tag)
	return &ResolvedCollectionNode{base, collection, ids}
}

func (rn *ResolvedCollectionNode) Locate(path ...string) (ResolvedNode, error) {
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

func (n *ResolvedCollectionNode) Resolve(id string) (ResolvedNode, error) {
	collection := n.Node.Collection
	if collection.Node.IsCollection {
		return n.ResolveCollection(collection.Tag, collection.Node, id)
	} else {
		return n.ResolveSingular(collection.Tag, collection.Node, id)
	}
}

func (n *ResolvedCollectionNode) Entity() interface{} {
	return n.Collection
}

func (n *ResolvedCollectionNode) SetEntity(e interface{}) {
	n.Collection = e
}

func (n *ResolvedCollectionNode) ResolveCollection(tag *Tag, collectionNode *Node, id string) (ResolvedNode, error) {
	collection, ids, err := collectionNode.ManifestCollection(n.Collection, id)
	if err != nil {
		return nil, err
	}
	if collection == nil {
		return nil, HttpError(404, "collection", n.ID, "does not have an item with ID", quot(id))
	}
	return newResolvedCollection(n, collectionNode, id, tag, collection, ids), nil
}

func (n *ResolvedCollectionNode) ResolveSingular(tag *Tag, singularNode *Node, id string) (ResolvedNode, error) {
	entity, err := singularNode.ManifestSingular(n.Collection, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, HttpError(404, "collection", n.ID, "does not have an item with ID", quot(id))
	}
	return newResolvedSingular(n, singularNode, id, tag, entity), nil
}

func (n *ResolvedCollectionNode) Links() ([]Link, error) {
	return []Link{Link{"self", n.Path()}}, nil
}

func (n *ResolvedCollectionNode) Resource() (*Resource, error) {
	if embeddedCollectionItems, err := n.EmbeddedCollectionItems(); err != nil {
		return nil, err
	} else if links, err := n.Links(); err != nil {
		return nil, err
	} else {
		return &Resource{n.Tag.Rel, nil, nil, embeddedCollectionItems, links}, nil
	}
}

func (n *ResolvedCollectionNode) EmbeddedResource(tag *Tag) (*Resource, error) {
	if len(tag.EmbedFields) != 0 {
		return n.FilteredEmbeddedResource(tag.EmbedFields)
	}
	return n.DefaultEmbeddedResource()
}

func (n *ResolvedCollectionNode) DefaultEmbeddedResource() (*Resource, error) {
	return n.Resource()
}

func (n *ResolvedCollectionNode) FilteredEmbeddedResource(fields []string) (*Resource, error) {
	// There are no fields to filter on a collection
	return n.Resource()
}

func (n *ResolvedCollectionNode) EmbeddedMembers() (map[string]*Resource, error) {
	return nil, nil
}

func (n *ResolvedCollectionNode) EmbeddedCollectionItems() ([]*Resource, error) {
	if !n.Node.IsCollection {
		return nil, nil
	}
	items := make([]*Resource, len(n.CollectionIDs))
	for i, id := range n.CollectionIDs {
		if childNode, err := n.Resolve(id); err != nil {
			return nil, err
		} else {
			if n.Tag.EmbedFields == nil {
				println("Embedding ", childNode.ID(), "; with all fields:")
			} else {
				println("Embedding ", childNode.ID(), "; with fields:", strings.Join(n.Tag.EmbedFields, ", "))
			}
			resource, err := childNode.EmbeddedResource(n.Tag)
			if err != nil {
				return nil, err
			} else {
				items[i] = resource
			}
		}
	}
	return items, nil
}
