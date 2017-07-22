package HTTPSplunkEvent

import (
	"net/http"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
)

var authHeaderKey string = "Authorization"
var authHeaderUser string = "Splunk "
var headerSplunkRequestChannel string = "X-Splunk-Request-Channel"
var endpoint string = "/services/collector/raw"

type HECWriter struct {
	server string
	token string
	index string
	requestChannel string
}

//NewHECWriter returns a new HECWriter.  The server value should not include the endpoint, which is added by the
// constructor
func NewHECWriter(server, token, index string) (*HECWriter, error) {
	channel, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &HECWriter{
		server: server + endpoint,
		token: token,
		requestChannel: channel.String(),
		index: index,
	}, nil
}

func (w HECWriter) Write(p []byte) (n int, err error) {
	c := http.Client{}

	event := Event{ Event: string(p), Index: "main" }
	outBuf := bytes.NewBuffer([]byte{})
	en := json.NewEncoder(outBuf)
	en.Encode(event)

	request, err := http.NewRequest(http.MethodPost, w.server, outBuf)
	if err != nil {
		return 0, err
	}
	request.Header.Add(authHeaderKey, authHeaderUser + w.token)
	request.Header.Add(headerSplunkRequestChannel, w.requestChannel)
	resp, err := c.Do(request)

	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New(resp.Status)
	}

	var hecResp Response
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(hecResp)
	if err != nil {
		return 0, nil
	}

	decoder.Decode(&hecResp)

	if hecResp.Code != 0 {
		return 0, errors.New(hecResp.Text)
	}

	return len(p), nil
}
