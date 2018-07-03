package consul

import (
	"errors"
	"fmt"
	"strings"

	"bitbucket.org/ConsentSystems/mango-micro/mango-service/registry"
	"github.com/hashicorp/consul/api"
	"github.com/pborman/uuid"
)

// NewRegistry returns a new registry instance
func NewRegistry(config *api.Config) registry.Registry {
	reg := &consulRegistry{}
	var conf *api.Config
	if config == nil {
		conf = api.DefaultConfig()
	} else {
		conf = config
	}
	reg.configs = conf
	return reg
}

type consulRegistry struct {
	configs *api.Config
}

func (cr consulRegistry) Register(name, address string, port int, tags []string) (string, error) {
	cl, err := api.NewClient(cr.configs)
	if err != nil {
		return "", err
	}
	serviceCOnf := &api.AgentServiceRegistration{
		ID:      uuid.NewUUID().String(),
		Name:    name,
		Port:    port,
		Address: address,
		Tags:    tags,
	}
	err = cl.Agent().ServiceRegister(serviceCOnf)
	if err != nil {
		return "", err
	}

	fmt.Println("Registered service")
	return serviceCOnf.ID, nil
}

func (cr consulRegistry) Deregister(id string) error {

	cl, err := api.NewClient(cr.configs)
	if err != nil {
		return err
	}

	err = cl.Agent().ServiceDeregister(id)
	if err != nil {
		return err
	}

	fmt.Println("Deregistered service")
	return nil
}

func (cr consulRegistry) Get(name string, version string) ([]string, error) {
	cl, err := api.NewClient(cr.configs)
	if err != nil {
		return []string{}, err
	}
	rsp, _, err := cl.Health().Service(name, "", false, nil)
	if err != nil {
		return []string{}, err
	}
	nodes := []*Node{}
	for _, s := range rsp {

		if s.Service.Service != name {
			continue
		}

		// version is now a tag
		servVers := decodeVersion(s.Service.Tags)
		if version != "" && servVers != version {
			continue
		}

		// service ID is now the node id
		id := s.Service.ID

		// address is service address
		address := s.Service.Address

		// use node address
		if len(address) == 0 {
			address = s.Node.Address
		}

		var del bool

		for _, check := range s.Checks {
			// delete the node if the status is critical
			if check.Status == "critical" {
				del = true
				break
			}
		}

		// if delete then skip the node
		if del {
			continue
		}

		nodes = append(nodes, &Node{
			ID:      id,
			Address: address,
			Port:    s.Service.Port,
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

func decodeVersion(tags []string) string {
	for i := 0; i < len(tags); i++ {
		parts := strings.Split(tags[i], "=")
		if parts[0] == "v" && len(parts) == 2 {
			return parts[1]
		}
	}

	return ""
}
