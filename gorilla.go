package ws

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

var once sync.Once

type gorilla struct {
	url    url.URL
	header http.Header

	resCh chan []byte
	finCh chan struct{}

	conn *websocket.Conn
	ctx  context.Context

	finErr error
}

func (g *gorilla) Conn(ctx context.Context) (func() error, error) {
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, g.url.String(), g.header)
	if err != nil {
		return nil, err
	}

	//some initial assignments
	g.conn = conn
	var cancelF context.CancelFunc
	g.ctx, cancelF = context.WithCancel(ctx)

	return func() error {
		//shutting down reading
		cancelF()

		//unsubscribing
		err := g.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			g.finErr = fmt.Errorf("write error: %w", err)
		}

		//closing connection
		return g.conn.Close()
	}, nil
}

func (g *gorilla) Write(subscribeMsgB []byte) (i int, err error) {
	var wcl io.WriteCloser
	wcl, err = g.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return
	}

	i, err = wcl.Write(subscribeMsgB)
	if err != nil {
		return
	}

	err = wcl.Close()
	if err != nil {
		return
	}

	return
}

func (g *gorilla) handling() {
	defer close(g.finCh)
handleLoop:
	for {
		select {
		case <-g.ctx.Done():
			break handleLoop
		default:
			_, message, err := g.conn.ReadMessage()
			if err != nil {
				g.finErr = fmt.Errorf("read error: %w", err)
				break handleLoop
			}
			g.resCh <- message
		}
	}
}

func (g *gorilla) Data() <-chan []byte {
	once.Do(func() {
		//handling func must be executed only once
		//because of closing of the signal channel
		go g.handling()
	})

	return g.resCh
}

func (g *gorilla) Done() <-chan struct{} {
	return g.finCh
}

func (g *gorilla) Err() error {
	return g.finErr
}
