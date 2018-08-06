package replier

import (
	"go/ast"
)

// MakeServerStr returns a string with the server instanciator
func makeServerStr(t *ast.TypeSpec) string {
	srvName := t.Name.Name + "Server"

	itp := t.Type.(*ast.InterfaceType)
	methodsStr := makeMethods(t.Name.Name, itp.Methods)
	server := `
// New` + srvName + ` returns a new replier server
func New` + srvName + ` (
	regAddr string,
	serverPort int,
	hdl  ` + t.Name.Name + `,
) {

	replier, _ := replier.NewServer(regAddr, "` + srvName + `-replier", "1")
	handler := handler.NewHandler()
	` + methodsStr + `
	replier.AddHandler(handler)
	replier.AddTransport(tcp.NewTransport())
	replier.AddTransport(inproc.NewTransport())
	go replier.Run(serverPort, "inproc", "` + srvName + `-replier")
	go replier.Run(serverPort, "tcp", "")
}
`

	return server
}

func makeHeader(pkg string) string {
	header := `package ` + pkg + `

import (
	"encoding/json"

	"bitbucket.org/ddanna79/mango-micro/handler"
	"bitbucket.org/ddanna79/mango-micro/replier"
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
