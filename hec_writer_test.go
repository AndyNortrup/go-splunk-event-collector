package HTTPSplunkEvent

import (
	"log"
	"testing"
)

const (
	index string = "main"
)

//This test should only compile if HECWriter is a valid implementation of io.Writer.
func TestLogWriterIsWriter(t *testing.T) {
	server := "address"
	token := "token"
	w, err := NewHECWriter(server, token, index)
	if err != nil {
		t.Fail()
		t.Log(err)
	}
	l := log.New(w, "", 0)
	l.Println("Test event")
}

func TestLogLocalSplunk(t *testing.T) {
	server := "http://localhost:8088"
	token := "F9D3427D-EAE4-42BA-8DA0-620A8EB2E2B5"
	hw, _ := NewHECWriter(server, token, index)
	l := log.New(hw, "", log.Ldate|log.Ltime)
	l.Print("test")
	//err := l.Output(0, "test")
	//if err != nil {
	//	t.Fail()
	//	t.Log(err)
	//}
}
