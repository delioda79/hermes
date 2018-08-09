package example1

import (
	"encoding/json"

	"bitbucket.org/ConsentSystems/mango-micro/pusher"
	"nanomsg.org/go-mangos/transport/inproc"
	"nanomsg.org/go-mangos/transport/tcp"
)

type APICallsHandlerClient interface {
	RegisterCall(msg APICallMessage) error
	NoParamsCall() error
}

type defaultAPICallsHandlerClient struct {
	psh pusher.Pusher
}

func (cl *defaultAPICallsHandlerClient) RegisterCall(msg APICallMessage) error {
	bts, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return cl.psh.Push("APICallsHandler.RegisterCall", bts)
}

func (cl *defaultAPICallsHandlerClient) NoParamsCall() error {
	bts := []byte{}
	return cl.psh.Push("APICallsHandler.NoParamsCall", bts)
}

// NewAPICallsHandlerClient  returns a handy client for the API Calls Push/Pull service
func NewAPICallsHandlerClient(
	registryAddr string,
	transport string,
	puller ...pusher.Puller,
) (APICallsHandlerClient, error) {
	cl, err := pusher.NewServer(registryAddr)
	if err != nil {
		return nil, err
	}

	cl.AddTransport(tcp.NewTransport())
	cl.AddTransport(inproc.NewTransport())

	cl.Run(puller...)
	return &defaultAPICallsHandlerClient{
		psh: cl.Sock(),
	}, nil
}
