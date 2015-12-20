package HTTPSplunkEvent

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strconv"
)

//Event is a struct that holds all data to be sent to a Splunk HTTP logging
//endpoint
type Event struct {
	Time       int64       `json:"time"`
	Host       string      `json:"host"`
	Source     string      `json:"source"`
	Sourcetype string      `json:"sourcetype"`
	Index      string      `json:"index"`
	Event      interface{} `json:"event"`
}

//Send the event to the specified Splunk Server
func (e Event) Send(destination string, token string, disableCertValidation bool) error {

	//Ensure we have all of the values for the event
	if e.Time == 0 || e.Host == "" || e.Source == "" || e.Sourcetype == "" || e.Index == "" || e.Event == nil {
		return errors.New("All fields in Event must have a value")
	}

	//Create a byte array with the data
	b, err := json.Marshal(e)

	//Create client and request
	client := &http.Client{}
	request, err := http.NewRequest("POST", destination, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	if disableCertValidation {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		client.Transport = tr
	}

	header := http.Header{}
	header.Add("Authorization", "Splunk "+token)
	header.Set("Content-Type", "application/json")
	request.Header = header

	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		dump, err := httputil.DumpResponse(resp, true)
		fmt.Printf("RESPONSE: %s\n", dump)
		return err
	}

	//Any code other than 200 is an
	var splunkResponse Response
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&splunkResponse)

	if err != nil {
		fmt.Printf("Error decoding resposne\n")
		return err
	}

	if splunkResponse.Code != SplunkResponseOK {

		fmt.Printf("Splunk Response code not OK\n")
		return errors.New(strconv.Itoa(splunkResponse.Code) + ": " + splunkResponse.Text)
	}

	//Everything is fine return nil error
	return nil

}
