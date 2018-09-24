package example2

import (
	"encoding/json"
	"errors"
	"time"

	"bitbucket.org/ConsentSystems/mango-micro/requester"
	"nanomsg.org/go-mangos/transport/inproc"
	"nanomsg.org/go-mangos/transport/tcp"
)
import "bitbucket.org/ConsentSystems/mango-micro/mango-service/registry"

// APICallsHandlerClient ...
type APICallsHandlerClient interface {
	TestBool(msg APICallMessage) (*APICallMessage, error)
}

type defaultAPICallsHandlerClient struct {
	rqstr       requester.Server
	deadline    time.Duration
	serviceName string
}

// SetDeadline Sets the deadline for the requests
func (cl *defaultAPICallsHandlerClient) SetDeadline(dr time.Duration) {
	cl.deadline = dr
}

// TestBool ...
func (cl *defaultAPICallsHandlerClient) TestBool(msg APICallMessage) (*APICallMessage, error) {

	bts, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	sck := cl.rqstr.Sock()
	sck.SetDeadline(cl.deadline)
	resBts, err := sck.Request(cl.serviceName+".TestBool", bts)
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

// NewAPICallsHandlerClient  returns a handy client for the API Calls RPC service
func NewAPICallsHandlerClient(
	registry registry.Registry,
	transport string,
	serviceName string,
	responder ...requester.Responder,
) (APICallsHandlerClient, error) {
	cl, err := requester.NewServer(registry)
	if err != nil {
		return nil, err
	}

	if serviceName == "" {
		serviceName = "APICallsHandler"
	}

	cl.AddTransport(tcp.NewTransport())
	cl.AddTransport(inproc.NewTransport())

	go cl.Run(responder...)

	return &defaultAPICallsHandlerClient{
		rqstr:       cl,
		deadline:    time.Second * 10,
		serviceName: serviceName,
	}, nil
}
