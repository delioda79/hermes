package registry

// Registry is an interface for service discovery
type Registry interface {
	Register(name, address string, port int, tags []string) (string, error)
	Deregister(id string) error
	Get(name string, version string) ([]string, error)
}
