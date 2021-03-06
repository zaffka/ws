package ws

import (
	"context"
)

//Configurer is an interface returning params needed to configure ws realization
type Configurer interface {
	URL() string
	Header() map[string][]string
}

//Handler interface is a simplified wrapper for the ws realization
type Handler interface {
	//Conn makes new connection to the remote system and returns the connection tear down func or an error
	Conn(context.Context) (func() error, error)

	//Write writes binary messages to the socket
	Write([]byte) (int, error)

	//Handle returns a channel from where the binary results are being read
	Handle(ctx context.Context) <-chan []byte

	//Done returns a channel signaling about Handler's halt
	//The first error received is a halting reason
	//After that channel will be closed and receive nil values
	Done() <-chan error
}

//New is a constructor function masking web socket realization
func New(cnf Configurer) Handler {
	return &gorilla{
		Configurer: cnf,
		doneCh:     make(chan error, 1),
	}
}
