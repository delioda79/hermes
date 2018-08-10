package requester

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
// New` + cltName + `  returns a handy client for the API Calls RPC service
func New` + cltName + `(
	registryAddr string,
	transport string,
	responder ...requester.Responder,
) (` + cltName + `, error) {
	cl, err := requester.NewServer(registryAddr)
	if err != nil {
		return nil, err
	}

	cl.AddTransport(tcp.NewTransport())
	cl.AddTransport(inproc.NewTransport())

	go cl.Run(responder...)

	return &default` + t.Name.Name + `Client{
		rqstr: cl,
		deadline: time.Second * 10,
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
	"errors"
	"time"

	"bitbucket.org/ConsentSystems/mango-micro/requester"
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
