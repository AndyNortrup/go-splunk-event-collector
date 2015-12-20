package HTTPSplunkEvent

import (
	"testing"
	"time"
)

func TestSendEvent(t *testing.T) {
	event := &Event{
		Time:       time.Now().Unix(),
		Index:      "tweet_harvest",
		Source:     "HTTPSplunkEvent",
		Sourcetype: "Test",
		Event:      "This is a test",
		Host:       "Golang",
	}

	err := event.Send("http://input-prd-p-zd5ktsgk9g47.cloud.splunk.com:8088/services/collector",
		"8CEB3C52-47B9-451E-81A9-4E45A299D41C", true)

	if err != nil {
		t.Fatalf("Error sending event: %#v\n", err.Error())
	}
}
