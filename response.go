package HTTPSplunkEvent

//Response is a struct sent back by Splunk with the results of a sent event
type Response struct {
	Text               string `json:"text"`
	Code               int    `json:"code"`
	InvalidEventNumber int    `json:"invalid-event-number"`
	AckID              int    `json:"ackId"`
}

//SplunkResponseOK is a code sent back from Splunk indicating that the request
// was processed satisfactorally
var SplunkResponseOK = 200
