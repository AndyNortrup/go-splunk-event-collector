# go-splunk-event-collector

[![Build Status](https://travis-ci.org/AndyNortrup/go-splunk-event-collector.svg?branch=master)](https://travis-ci.org/AndyNortrup/go-splunk-event-collector)
[![Coverage Status](https://coveralls.io/repos/github/AndyNortrup/go-splunk-event-collector/badge.svg?branch=master)](https://coveralls.io/github/AndyNortrup/go-splunk-event-collector?branch=master)


A go implementation for sending events to [Splunk's HTTP Event Collector](http://dev.splunk.com/view/event-collector/SP-CAAAE6M)

The HEC writer is an io.Writer compliant struct that can be used directly or with the log.Logger object to send logs to
Splunk Event Collector. Currently uses the raw endpoint so that the contents of your events will be written directly as an event,
so you have to ensure that you include a time stamp for Splunk to index.

```
	server := "http://localhost:8088"
	token := "<<your HEC token here>>"
	index := "main"
	hw, _ := NewHECWriter(server, token, index)
	l := log.New(hw, "", log.Ldate|log.Ltime)
	l.Print("test")
```

An unfortunate property of the log.Logger is that it returns no errors produced by the writer, so you don't have a great
mechanism for ensuring that your writes were completed.