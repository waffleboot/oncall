package model

import "time"

type Node struct {
	ID        int
	Name      string
	DeletedAt time.Time
}

func (s *Node) IsDeleted() bool {
	return !s.DeletedAt.IsZero()
}

func (s *Item) ActiveNodes() []Node {
	nodes := make([]Node, 0, len(s.Nodes))
	for _, node := range s.Nodes {
		if !node.IsDeleted() {
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func (s *Item) CreateNode() Node {
	var maxID int
	for i := range s.Nodes {
		node := s.Nodes[i]
		if node.ID > maxID {
			maxID = node.ID
		}
	}
	node := Node{ID: maxID + 1}
	s.Nodes = append(s.Nodes, node)
	return node
}

func (s *Item) UpdateNode(node Node) {
	for i := range s.Nodes {
		if s.Nodes[i].ID == node.ID {
			s.Nodes[i] = node
			break
		}
	}
}

func (s *Item) DeleteNode(node Node, at time.Time) {
	for i := range s.Nodes {
		if s.Nodes[i].ID == node.ID {
			s.Nodes[i].DeletedAt = at
			break
		}
	}
}
