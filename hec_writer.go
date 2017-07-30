package HTTPSplunkEvent

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"time"
)

var authHeaderKey string = "Authorization"
var authHeaderUser string = "Splunk "
var headerSplunkRequestChannel string = "X-Splunk-Request-Channel"
var endpoint string = "/services/collector"

type HECWriter struct {
	server         string
	token          string
	index          string
	requestChannel string
	host           string
	source         string
	sourcetype     string
	client         *http.Client
}

//NewHECWriter returns a new HECWriter.  The server value should not include the endpoint, which is added by the
// constructor
func NewHECWriter(server, token, index, host, source, sourcetype string, allowInsecure bool) (*HECWriter, error) {
	channel, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &HECWriter{
		server:         server + endpoint,
		token:          token,
		requestChannel: channel.String(),
		index:          index,
		host:           host,
		source:         source,
		sourcetype:     sourcetype,
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: allowInsecure},
			},
		},
	}, nil
}

func (w HECWriter) Write(p []byte) (n int, err error) {

	//TODO: Figure out how to handle setting time in the JSON vs just letting splunk pull it out of the event data (p)
	event := NewEvent(time.Now().Unix(), w.host, w.source, w.sourcetype, w.index, p)
	outBuf := bytes.NewBuffer([]byte{})
	en := json.NewEncoder(outBuf)
	en.Encode(event)

	request, err := http.NewRequest(http.MethodPost, w.server, outBuf)
	if err != nil {
		return 0, err
	}

	request.Header.Add(authHeaderKey, authHeaderUser+w.token)
	request.Header.Add(headerSplunkRequestChannel, w.requestChannel)
	resp, err := w.client.Do(request)

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
