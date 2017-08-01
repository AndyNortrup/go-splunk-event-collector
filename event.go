package HTTPSplunkEvent


//Event is a struct that holds all data to be sent to a Splunk HTTP logging
//endpoint
type Event struct {
	Time       int64       `json:"time"`
	Host       string      `json:"host"`
	Source     string      `json:"source"`
	Sourcetype string      `json:"sourcetype"`
	Index      string      `json:"index"`
	Text       string `json:"event"`
}

func NewEvent(time int64, host, source, sourcetype, index string, text string) *Event {
	return &Event{
		Time: time,
		Host: host,
		Source: source,
		Sourcetype: sourcetype,
		Index: index,
		Text: text,
	}
}
