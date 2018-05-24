package consul

import "fmt"

// Service represents a service in Consul
type Service struct {
	Name    string
	Version int
	Nodes   []*Node
}

// Node represents a node for a service
type Node struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	Port    int    `json:"port"`
}

// GetURL returns a string which can be used as URL for a socket connection
func (nd Node) GetURL() string {
	return fmt.Sprintf("%s:%d", nd.Address, nd.Port)
}
