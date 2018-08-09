package example2

import (
	"encoding/json"
	"errors"

	"bitbucket.org/ConsentSystems/mango-micro/requester"
	"nanomsg.org/go-mangos/transport/inproc"
	"nanomsg.org/go-mangos/transport/tcp"
)
import "bitbucket.org/ConsentSystems/mango-micro/messages"

// APICallsHandlerClient ...
type APICallsHandlerClient interface {
	RegisterCall(msg APICallMessage) (*APICallMessage, error)
	External(msg messages.Trigger) (*messages.Trigger, error)
}

type defaultAPICallsHandlerClient struct {
	rqstr requester.Server
}

// RegisterCall ...
func (cl *defaultAPICallsHandlerClient) RegisterCall(msg APICallMessage) (*APICallMessage, error) {

	bts, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	resBts, err := cl.rqstr.Sock().Request("APICallsHandler.RegisterCall", bts)
	if err != nil {
		return nil, err
	}
	resArr := &[]*[]byte{}
	json.Unmarshal(resBts, resArr)
	rsp := &APICallMessage{}
	json.Unmarshal(*(*resArr)[0], rsp)
	if len(*(*resArr)[1]) > 0 {
		return nil, errors.New(string(*(*resArr)[1]))
	}
	return rsp, nil
}

// External ...
func (cl *defaultAPICallsHandlerClient) External(msg messages.Trigger) (*messages.Trigger, error) {

	bts, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}

	resBts, err := cl.rqstr.Sock().Request("APICallsHandler.External", bts)
	if err != nil {
		return nil, err
	}
	resArr := &[]*[]byte{}
	json.Unmarshal(resBts, resArr)
	rsp := &messages.Trigger{}
	json.Unmarshal(*(*resArr)[0], rsp)
	if len(*(*resArr)[1]) > 0 {
		return nil, errors.New(string(*(*resArr)[1]))
	}
	return rsp, nil
}

// NewAPICallsHandlerClient  returns a handy client for the API Calls RPC service
func NewAPICallsHandlerClient(
	registryAddr string,
	transport string,
	responder ...requester.Responder,
) (APICallsHandlerClient, error) {
	cl, err := requester.NewServer(registryAddr)
	if err != nil {
		return nil, err
	}

	cl.AddTransport(tcp.NewTransport())
	cl.AddTransport(inproc.NewTransport())

	go cl.Run(responder...)

	return &defaultAPICallsHandlerClient{
		rqstr: cl,
	}, nil
}
