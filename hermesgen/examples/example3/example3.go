package example2

import (
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
	TestBool(*APICallMessage) (*APICallMessage, error)
}
