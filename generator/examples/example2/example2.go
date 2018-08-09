package example2

import (
	"bitbucket.org/ConsentSystems/mango-micro/messages"
	"gopkg.in/mgo.v2/bson"
)

// APICallMessage represents a call to the API log message
type APICallMessage struct {
	OrganisationID bson.ObjectId
	UserID         bson.ObjectId
	Method         string
	Params         map[string]interface{}
}

// APICallsHandler handles teh API calls cunt
type APICallsHandler interface {
	RegisterCall(*APICallMessage) (*APICallMessage, error)
	External(*messages.Trigger) (*messages.Trigger, error)
}
