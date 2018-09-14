package memory

import (
	"errors"
	"fmt"
	"strings"

	"bitbucket.org/ConsentSystems/mango-micro/mango-service/registry"
	"github.com/micro/util/go/lib/addr"
	"github.com/pborman/uuid"
)

type Node struct {
	registry.Node
	Tags []string
}

type Service struct {
	registry.Service
	Nodes []*Node
}
type Registry struct {
	services []Service
}

func (rg *Registry) Register(name, address string, port int, tags []string) (string, error) {

	privateAddr, err := addr.Extract(address)
	if err != nil {
		return "", err
	}

	nd := rg.findService(name)
	if nd == nil {
		nd = &Service{
			Service: registry.Service{
				Name: name,
			},
			Nodes: []*Node{},
		}
		rg.services = append(rg.services, *nd)
	}

	node := &Node{
		Node: registry.Node{
			ID:      uuid.NewUUID().String(),
			Address: privateAddr,
			Port:    port,
		},
		Tags: tags,
	}

	nd.Nodes = append(nd.Nodes, node)

	fmt.Println("Registered service: ", privateAddr)
	return node.ID, nil
}

func (rg *Registry) Deregister(id string) error {
	for _, srv := range rg.services {
		nd, pos := srv.FindNode(id)
		if nd != nil {
			srvs := rg.services
			srvs[len(srvs)-1], srvs[pos] = srvs[pos], srvs[len(srvs)-1]
			rg.services = srvs[:len(srvs)-1]
		}
	}
	fmt.Println("Deregistered service")
	return nil
}

func (rg *Registry) Get(name string, version, transport string) ([]string, error) {
	res := []string{}
	srv := rg.findService(name)
	nodes := []*registry.Node{}
	if srv != nil {
		for _, nd := range srv.Nodes {
			vrs := decodeVersion(nd.Tags)
			if version != "" && vrs != version {
				continue
			}

			servTransport := decodeTransp(nd.Tags)
			if transport != "" && servTransport != transport {
				continue
			}

			nodes = append(nodes, &registry.Node{
				ID:      nd.ID,
				Address: nd.Address,
				Port:    nd.Port,
			})
		}

		if len(nodes) < 1 {
			return []string{}, errors.New("Service not found")
		}

		urls := []string{}
		for i := 0; i < len(nodes); i++ {
			urls = append(urls, nodes[i].GetURL())
		}

		return urls, nil
	}
	return res, nil
}

func (rg *Registry) findService(name string) *Service {
	for _, v := range rg.services {
		if v.Name == name {
			return &v
		}
	}

	return nil
}

func decodeVersion(tags []string) string {
	for i := 0; i < len(tags); i++ {
		parts := strings.Split(tags[i], "=")
		if parts[0] == "v" && len(parts) == 2 {
			return parts[1]
		}
	}

	return ""
}

func decodeTransp(tags []string) string {
	for i := 0; i < len(tags); i++ {
		parts := strings.Split(tags[i], "=")
		if parts[0] == "transport" && len(parts) == 2 {
			return parts[1]
		}
	}

	return ""
}
