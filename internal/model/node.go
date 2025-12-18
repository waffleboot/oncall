package model

import (
	"fmt"
	"strings"
	"time"
)

type Node struct {
	ID          int
	Name        string
	Description string
	DeletedAt   time.Time
}

func (n *Node) NotDeleted() bool {
	return n.DeletedAt.IsZero()
}

func (n *Node) MenuItem() string {
	if n.Name == "" {
		return "empty"
	}
	return n.Name
}

func (n *Node) ToPrint() string {
	description := n.Description
	if description != "" {
		description = "\n\n" + description
	}
	return fmt.Sprintf("host: %s%s", n.Name, description)
}

func (s *Item) ActiveNodes() []Node {
	nodes := make([]Node, 0, len(s.Nodes))
	for _, node := range s.Nodes {
		if node.NotDeleted() {
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func (s *Item) CreateNode() Node {
	return Node{}
}

func (s *Item) UpdateNode(node Node) {
	var maxID int
	var found bool
	for i, n := range s.Nodes {
		if n.ID == node.ID {
			s.Nodes[i] = node
			return
		}
		if n.Name == node.Name {
			found = true
		}
		if n.ID > maxID {
			maxID = n.ID
		}
	}
	if found {
		return
	}
	node.ID = maxID + 1
	s.Nodes = append(s.Nodes, node)
}

func (s *Item) DeleteNode(node Node) {
	for i := range s.Nodes {
		if s.Nodes[i].ID == node.ID {
			s.Nodes[i].DeletedAt = time.Now()
			return
		}
	}
}

func (n *Node) Printed() bool {
	return n.NotDeleted() && strings.TrimSpace(n.Name) != ""
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
