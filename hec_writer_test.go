package HTTPSplunkEvent

import (
	"log"
	"testing"
)

const (
	index string = "main"
	source string = "go-splunk-event-collector"
	sourcetype = "go-splunk-event-collector"
	host = "unit-tester"

)

func TestLogLocalSplunk(t *testing.T) {
	server := "http://splunks:8088"
	token := "122D68E5-EE08-4416-8FE6-A2CFDCF0F0A2"
	hw, _ := NewHECWriter(server, token, index, host, source, sourcetype,true)
	l := log.New(hw, "", log.Ldate|log.Ltime)

	err := l.Output(0, "test")
	if err != nil {
		t.Fail()
		t.Log(err)
	}
}
