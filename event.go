package HTTPSplunkEvent

//Event is a struct that holds all data to be sent to a Splunk HTTP logging
//endpoint
type Event struct {
	Time       int64   `json:"time"`
	Host       string  `json:"host"`
	Source     string  `json:"source"`
	Sourcetype string  `json:"sourcetype"`
	Index      string  `json:"index"`
	Event      []Value `json:"event"`
}
