package pusher

import (
	"go/ast"
)

// MakeClientStr returns a string with the server instanciator
func makeClientStr(t *ast.TypeSpec) string {
	cltName := t.Name.Name + "Client"

	itp := t.Type.(*ast.InterfaceType)
	intf := makeInterface(cltName, itp.Methods)
	mtds := makeMethods(t.Name.Name, itp.Methods)
	// methodsStr := MakeMethods(t.Name.Name, itp.Methods)
	server := `
// New` + cltName + `  returns a handy client for the API Calls Push/Pull service
func New` + cltName + `(
	registry registry.Registry,
	transport string,
	serviceName string,
	puller ...pusher.Puller,
) (` + cltName + `, error) {
	cl, err := pusher.NewServer(registry)
	if err != nil {
		return nil, err
	}

	cl.AddTransport(tcp.NewTransport())
	cl.AddTransport(inproc.NewTransport())

	cl.Run(puller...)
	return &default` + t.Name.Name + `Client{
		psh: cl.Sock(),
		serviceName: serviceName,
	}, nil
}
`

	return intf + mtds + server
}

// MakeHeader returns a heaer for the pusher client
func makeHeader(pkg string) string {
	header := `package ` + pkg + `

	import (
	"encoding/json"

	"bitbucket.org/ConsentSystems/mango-micro/pusher"
	"nanomsg.org/go-mangos/transport/inproc"
	"nanomsg.org/go-mangos/transport/tcp"
)
	`
	return header
}

// Generator is a server generator for push clients
type Generator struct {
}

// MakeHeader returns a heaer for the puller server
func (pg *Generator) MakeHeader(name string) string {
	return makeHeader(name)
}

// MakeBody returns the code for running the puller server
func (pg *Generator) MakeBody(t *ast.TypeSpec) string {
	return makeClientStr(t)
}
