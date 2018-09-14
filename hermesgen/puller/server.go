package puller

import (
	"go/ast"
)

// MakeServerStr returns a string with the server instanciator
func makeServerStr(t *ast.TypeSpec) string {
	srvName := t.Name.Name + "Server"

	itp := t.Type.(*ast.InterfaceType)
	methodsStr := makeMethods(t.Name.Name, itp.Methods)
	server := `
// New` + srvName + ` returns a new puller server
func New` + srvName + ` (
	registry registry.Registry,
	portStr string,
	hdl  ` + t.Name.Name + `,
) (puller.Server, error) {
	pullserver, _ := puller.NewServer(
		registry,
		"` + srvName + `-puller",
		"1",
	)

	pullserver.AddTransport(inproc.NewTransport())
	pullserver.AddTransport(tcp.NewTransport())

	handler := handler.NewHandler()

	` + methodsStr + `

	pullserver.AddHandler(handler)

	return pullserver, nil
}
`

	return server
}

func makeHeader(pkg string) string {
	header := `package ` + pkg + `

import (
	"encoding/json"
	"fmt"
	"strconv"

	"bitbucket.org/ConsentSystems/mango-micro/handler"
	"bitbucket.org/ConsentSystems/mango-micro/puller"
	"nanomsg.org/go-mangos/transport/inproc"
	"nanomsg.org/go-mangos/transport/tcp"
)
`
	return header
}

// Generator is a server generator for pull servers
type Generator struct {
}

// MakeHeader returns a heaer for the puller server
func (pg *Generator) MakeHeader(name string) string {
	return makeHeader(name)
}

// MakeBody returns the code for running the puller server
func (pg *Generator) MakeBody(t *ast.TypeSpec) string {
	return makeServerStr(t)
}
