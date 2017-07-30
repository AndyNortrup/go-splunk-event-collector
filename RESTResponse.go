package HTTPSplunkEvent

import (
	"bytes"
	"encoding/xml"
	"net/http"
)

/*
<?xml version="1.0" encoding="UTF-8"?>
<response>
  <messages>
    <msg type="ERROR">Unauthorized</msg>
  </messages>
</response>
*/

//RESTResponse converts an error message back from the server
type RESTResponse struct {
	Messages []Message `xml:"messages"`
}

//Message holds the list of messages
type Message struct {
	Message string `xml:"msg"`
	Type    string `xml:"type,attr"`
}

//NewRESTResponse creates a RESTResponse from an http.Response
func NewRESTResponse(input *http.Response) (RESTResponse, error) {
	decode := xml.NewDecoder(input.Body)
	var response RESTResponse
	err := decode.Decode(&response)
	return response, err
}

//String converts the RESTResponse into a string for use in errors or output
func (response RESTResponse) String() string {
	var buf bytes.Buffer
	for _, msg := range response.Messages {
		buf.WriteString(msg.Type)
		buf.WriteString(": ")
		buf.WriteString(msg.Message)
		buf.WriteByte('\n')
	}
	return buf.String()
}
