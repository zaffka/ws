package ws

import (
	"context"
	"net/http"
	"net/url"
)

//Config is a struct holding params needed to configure ws realization
type Config struct {
	URL    url.URL
	Header http.Header
}

//Handler interface is a simplified wrapper for the ws realization
type Handler interface {
	//Conn makes new connection to the remote system and returns the connection tear down func or an error
	Conn(context.Context) (func() error, error)

	//Write writes binary messages to the socket
	Write([]byte) (int, error)

	//DataChan returns a channel from where the results are being read
	DataChan(context.Context) <-chan []byte

	//Err holds any error the Handler has when finished its work
	Err() error
}

//New is a constructor function masking web socket realization
func New(cnf Config) Handler {
	return &gorilla{
		url:    cnf.URL,
		header: cnf.Header,
	}
}
