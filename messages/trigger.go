package messages

// Trigger is a trigger for our hooks
type Trigger struct {
	Name   string
	Params []byte
	UID    string
}
