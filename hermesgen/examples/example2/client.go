package example2

import (
	"encoding/json"
	"errors"
	"time"

	"bitbucket.org/ConsentSystems/mango-micro/requester"
	"nanomsg.org/go-mangos/transport/inproc"
	"nanomsg.org/go-mangos/transport/tcp"
)
import "bitbucket.org/ConsentSystems/mango-micro/messages"

// APICallsHandlerClient ...
type APICallsHandlerClient interface {
	RegisterCall(msg APICallMessage) (*APICallMessage, error)
	External(msg messages.Trigger) (*messages.Trigger, error)
	NoParams() (*messages.Trigger, error)
}

type defaultAPICallsHandlerClient struct {
	rqstr    requester.Server
	deadline time.Duration
}

// SetDeadline Sets the deadline for the requests
func (cl *defaultAPICallsHandlerClient) SetDeadline(dr time.Duration) {
	cl.deadline = dr
}

// RegisterCall ...
func (cl *defaultAPICallsHandlerClient) RegisterCall(msg APICallMessage) (*APICallMessage, error) {

	bts, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	sck := cl.rqstr.Sock()
	sck.SetDeadline(cl.deadline)
	resBts, err := sck.Request("APICallsHandler.RegisterCall", bts)
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
	sck := cl.rqstr.Sock()
	sck.SetDeadline(cl.deadline)
	resBts, err := sck.Request("APICallsHandler.External", bts)
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

// NoParams ...
func (cl *defaultAPICallsHandlerClient) NoParams() (*messages.Trigger, error) {
	sck := cl.rqstr.Sock()
	sck.SetDeadline(cl.deadline)
	resBts, err := sck.Request("APICallsHandler.NoParams", []byte{})
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
		rqstr:    cl,
		deadline: time.Second * 10,
	}, nil
}
