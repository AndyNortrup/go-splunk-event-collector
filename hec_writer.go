package HTTPSplunkEvent

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
	"io"
)

var authHeaderKey string = "Authorization"
var authHeaderUser string = "Splunk "
var headerSplunkRequestChannel string = "X-Splunk-Request-Channel"
var endpointRaw string = "/services/collector/raw"
var endpointStandard string = "/services/collector"

type HECWriter struct {
	server         string
	token          string
	index          string
	requestChannel string
	host           string
	source         string
	sourcetype     string
	endpoint       string

	//TimeFunc is the function used to set the event time when time extraction is not used.
	TimeFunc func() int64

	client *http.Client
}

var InvalidTokenError error = errors.New("Invalid Token")
var ServerNotFoundError error = errors.New("Server Not Found")

//NewHECWriter returns a new HECWriter.  The server value should not include the endpoint, which is added by the
// constructor
func NewHECWriter(server, token, index, host, source, sourcetype string, allowInsecure bool) (*HECWriter, error) {
	channel, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &HECWriter{
		server:         server,
		token:          token,
		requestChannel: channel.String(),
		index:          index,
		host:           host,
		source:         source,
		sourcetype:     sourcetype,
		endpoint:       endpointRaw,
		TimeFunc:       func() int64 { return time.Now().UnixNano() },
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: allowInsecure},
			},
		},
	}, nil
}

func (w *HECWriter) Write(p []byte) (n int, err error) {

	outBuf, err := w.createEvent(p)
	if err != nil {
		log.Print(err)
		return 0, err
	}

	request, err := http.NewRequest(http.MethodPost, w.getDest(), outBuf)
	if err != nil {
		return 0, err
	}

	request.Header.Add(authHeaderKey, authHeaderUser+w.token)
	request.Header.Add(headerSplunkRequestChannel, w.requestChannel)
	resp, err := w.client.Do(request)

	if err != nil {
		return 0, ServerNotFoundError
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("%s - %s", resp.StatusCode, resp.Status)
		if resp.StatusCode == http.StatusForbidden {
			return 0, InvalidTokenError
		}

		return 0, errors.New(resp.Status)
	}

	var hecResp Response
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(hecResp)
	if err != nil {
		return 0, err
	}

	decoder.Decode(&hecResp)

	if hecResp.Code != 0 {
		return 0, errors.New(hecResp.Text)
	}

	return len(p), nil
}

func (w *HECWriter) createEvent(p []byte) (io.Reader, error) {
	event := NewEvent(w.TimeFunc(), w.host, w.source, w.sourcetype, w.index, string(p))
	outBuf := bytes.NewBuffer([]byte{})
	en := json.NewEncoder(outBuf)
	err := en.Encode(event)
	return outBuf, err
}

// UseRawEndpoint indicates if data should be submitted to the HEC Raw endpoint or the standard endpoint.
// Using the Raw endpoint implies that dates are included in the event for extraction.  If you are using the writer as
// to support a log.Logger, the correct value is probably True because logger should include the date and time
// at the beginning of the event, which Splunk will then extract.
//
// If you happen to be using this as a standalone writer, you might want to set this as off and provide an appropriate
// HECWriter.TimeFunc() to add the time to the event.
func (w *HECWriter) UseRawEndpoint(extract bool) {
	if extract {
		w.endpoint = endpointRaw
		w.TimeFunc = w.rawTimeFunc
	} else {
		w.endpoint = endpointStandard
		w.TimeFunc = w.nowTimeFunc
	}
}

func (w *HECWriter) getDest() string {
	return w.server + w.endpoint
}

func (w *HECWriter) rawTimeFunc() int64 {
	return 0
}

func (w *HECWriter) nowTimeFunc() int64 {
	return time.Now().UnixNano()
}
