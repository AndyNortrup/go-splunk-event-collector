package HTTPSplunkEvent

import (
	"log"
	"testing"
	"gopkg.in/ory-am/dockertest.v3"
	"os"
	"net/http"
	"errors"
)

const (
	index string = "main"
	source string = "go-splunk-event-collector"
	sourcetype = "go-splunk-event-collector"
	host = "unit-tester"

)

func TestLogLocalSplunk(t *testing.T) {
	token := "122D68E5-EE08-4416-8FE6-A2CFDCF0F0A2"
	hw, _ := NewHECWriter(dockerSplunkHostHEC, token, index, host, source, sourcetype,true)
	l := log.New(hw, "", log.Ldate|log.Ltime)

	err := l.Output(0, "test")
	if err != nil {
		t.Fail()
		t.Log(err)
	}
}

var dockerSplunkHostHEC string

func TestMain(m *testing.M) {
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

	dockerSplunkHostHEC = "https://" + resource.GetBoundIP("8088/tcp") + ":" + resource.GetPort("8088/tcp")
	dockerSplunkHostWeb := "http://" + resource.GetBoundIP("8000/tcp") + ":" + resource.GetPort("8000/tcp")

	if err := pool.Retry(func() error {
		var err error

		resp, err := http.Get(dockerSplunkHostWeb)
		if err != nil {
			return err
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK  {
			return nil
		} else {
			return errors.New("wrong status")
		}


	}); err != nil {
		pool.Purge(resource)
		log.Fatalf("Failed to connect to service: %s", err)
	}

	code := m.Run()

	err = pool.Purge(resource)
	if err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
