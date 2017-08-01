package HTTPSplunkEvent

import (
	"errors"
	"gopkg.in/ory-am/dockertest.v3"
	"log"
	"net/http"
	"testing"
	"encoding/json"
)

const (
	index         string = "main"
	source        string = "go-splunk-event-collector"
	sourcetype    string = "go-splunk-event-collector"
	host          string = "unit-tester"
	validToken    string = "122D68E5-EE08-4416-8FE6-A2CFDCF0F0A2"
	invalidToken  string = "invalid-token"
	invalidServer string = "http://invalidServer:8088"
)

func Test_IntegrationWrite(t *testing.T) {

	pool, resource := startHECContainer()
	defer pool.Purge(resource)

	validHECHost := hecHostFromResource(resource)

	tokens := []string{validToken, invalidToken, validToken}
	errValues := []error{nil, InvalidTokenError, ServerNotFoundError}
	servers := []string{validHECHost, validHECHost, invalidServer}

	for key, token := range tokens {
		hw, _ := NewHECWriter(servers[key], token, index, host, source, sourcetype, true)
		l := log.New(hw, "", log.Ldate|log.Ltime)
		err := l.Output(0, "test")
		if err != errValues[key] {
			t.Fail()
			t.Log(err)
		}
	}
}

func TestHECWriter_UseRawEndpoint(t *testing.T) {
	w, err := NewHECWriter("", "", "", "", "", "", true)
	if err != nil {
		t.Fail()
		t.Log(err)
	}

	w.UseRawEndpoint(false)
	if w.endpoint != endpointStandard {
		t.Fail()
		t.Logf("Wrong endpoint returned. Want=%s,\t Got=%s", endpointStandard, w.endpoint)
	}

	w.UseRawEndpoint(true)
	if w.endpoint != endpointRaw {
		t.Fail()
		t.Logf("Wrong endpoint returned. Want=%s,\t Got=%s", endpointRaw, w.endpoint)
	}
}

func TestHECWriter_RawTimeFunc(t *testing.T) {
	w := &HECWriter{}
	if w.rawTimeFunc() != 0 {
		t.Fail()
		t.Logf("Wrong value returned. Got=%s,\t Want=0", w.rawTimeFunc())
	}
}

func TestHECWriter_createEvent(t *testing.T) {
	message := "Hello World!"
	w, _ := NewHECWriter(invalidServer, validToken, index, host, source, sourcetype, true)
	buf, err := w.createEvent([]byte(message))

	if err != nil {
		t.Fail()
		t.Log(err)
	}

	dec := json.NewDecoder(buf)
	outEvent := &Event{}
	dec.Decode(outEvent)

	if message != outEvent.Text {
		t.Fail()
		t.Log("Failed to encode message properly.\t Want=%s,\tGot=%s", message, outEvent.Text)
	}
}

func startHECContainer() (*dockertest.Pool, *dockertest.Resource) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run(
		"andynortrup/go-splunk-event-collector",
		"latest",
		[]string{})

	if err != nil {
		log.Fatalf("Could not start Splunk: %s", err)
	}

	dockerSplunkHostWeb := "http://" + resource.GetBoundIP("8000/tcp") + ":" + resource.GetPort("8000/tcp")

	if err := pool.Retry(func() error {
		var err error

		resp, err := http.Get(dockerSplunkHostWeb)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return nil
		} else {
			return errors.New("wrong status")
		}

	}); err != nil {
		pool.Purge(resource)
		log.Fatalf("Failed to connect to service: %s", err)
	}

	return pool, resource
}

func hecHostFromResource(resource *dockertest.Resource) string {
	return "https://" + resource.GetBoundIP("8088/tcp") + ":" + resource.GetPort("8088/tcp")
}
