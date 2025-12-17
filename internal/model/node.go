package model

import (
	"fmt"
	"strings"
	"time"
)

type Node struct {
	ID        int
	Name      string
	DeletedAt time.Time
}

func (n *Node) IsDeleted() bool {
	return !n.DeletedAt.IsZero()
}

func (n *Node) MenuItem() string {
	if n.Name == "" {
		return "empty"
	}
	return n.Name
}

func (n *Node) ToPrint() string {
	return fmt.Sprintf("host: %s", n.Name)
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

func (n *Node) Printed() bool {
	return !n.IsDeleted() && strings.TrimSpace(n.Name) != ""
}

func (s *Item) PrintedNodes() []Node {
	m := make(map[string]Node)

	for _, node := range s.Nodes {
		if node.Printed() {
			m[node.Name] = node
		}
	}

	for _, vm := range s.PrintedVMs() {
		delete(m, vm.Node)
	}

	nodes := make([]Node, 0, len(m))
	for _, node := range m {
		nodes = append(nodes, node)
	}

	return nodes
}
